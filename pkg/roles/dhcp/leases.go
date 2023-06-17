package dhcp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/dhcp/types"
	"beryju.io/gravity/pkg/roles/dns/utils"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/rfc1035label"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type Lease struct {
	inst roles.Instance

	scope      *Scope
	log        *zap.Logger
	Identifier string `json:"-"`

	Address          string `json:"address"`
	Hostname         string `json:"hostname"`
	AddressLeaseTime string `json:"addressLeaseTime,omitempty"`
	ScopeKey         string `json:"scopeKey"`
	DNSZone          string `json:"dnsZone,omitempty"`

	etcdKey string
}

func (r *Role) NewLease(identifier string) *Lease {
	return &Lease{
		inst:       r.i,
		Identifier: identifier,
		log:        r.log.With(zap.String("identifier", identifier)),
	}
}

func (r *Role) leaseFromKV(raw *mvccpb.KeyValue) (*Lease, error) {
	prefix := r.i.KV().Key(
		types.KeyRole,
		types.KeyLeases,
	).Prefix(true).String()
	identifier := strings.TrimPrefix(string(raw.Key), prefix)
	l := r.NewLease(identifier)
	err := json.Unmarshal(raw.Value, &l)
	if err != nil {
		return l, err
	}
	l.etcdKey = string(raw.Key)

	r.scopesM.RLock()
	scope, ok := r.scopes[l.ScopeKey]
	r.scopesM.RUnlock()
	if !ok {
		return l, fmt.Errorf("DHCP lease with invalid scope key: %s", l.ScopeKey)
	}
	l.scope = scope
	return l, nil
}

func (l *Lease) Put(ctx context.Context, expiry int64, opts ...clientv3.OpOption) error {
	raw, err := json.Marshal(&l)
	if err != nil {
		return err
	}

	if expiry > 0 {
		exp, err := l.inst.KV().Lease.Grant(ctx, expiry)
		if err != nil {
			return err
		}
		opts = append(opts, clientv3.WithLease(exp.ID))
	}

	leaseKey := l.inst.KV().Key(
		types.KeyRole,
		types.KeyLeases,
		l.Identifier,
	)
	_, err = l.inst.KV().Put(
		ctx,
		leaseKey.String(),
		string(raw),
		opts...,
	)
	if err != nil {
		return err
	}

	var zone string
	if l.scope != nil && l.scope.DNS != nil {
		zone = l.scope.DNS.Zone
	}
	if l.DNSZone != "" {
		zone = l.DNSZone
	}
	ev := roles.NewEvent(
		ctx,
		map[string]interface{}{
			"hostname":   l.Hostname,
			"address":    l.Address,
			"identifier": l.Identifier,
			"fqdn":       utils.EnsureTrailingPeriod(strings.Join([]string{l.Hostname, zone}, ".")),
		},
	)
	ev.Payload.RelatedObjectKey = leaseKey
	ev.Payload.RelatedObjectOptions = opts
	l.inst.DispatchEvent(types.EventTopicDHCPLeasePut, ev)

	l.log.Debug("put lease", zap.Int64("expiry", expiry))
	go l.scope.calculateUsage()
	return nil
}

func (l *Lease) createReply(req *Request4) *dhcpv4.DHCPv4 {
	rep, err := dhcpv4.NewReplyFromRequest(req.DHCPv4)
	if err != nil {
		req.log.Warn("failed to create reply", zap.Error(err))
		return nil
	}
	rep.UpdateOption(dhcpv4.OptSubnetMask(l.scope.ipam.GetSubnetMask()))
	rep.UpdateOption(dhcpv4.OptIPAddressLeaseTime(time.Duration(l.scope.TTL * int64(time.Second))))

	if l.AddressLeaseTime != "" {
		pl, err := time.ParseDuration(l.AddressLeaseTime)
		if err != nil {
			req.log.Warn("failed to parse address lease duration, defaulting", zap.Error(err), zap.String("default", pl.String()))
		} else {
			rep.UpdateOption(dhcpv4.OptIPAddressLeaseTime(pl))
		}
	}

	// DNS Options
	rep.UpdateOption(dhcpv4.OptDNS(net.ParseIP(extconfig.Get().Instance.IP)))
	if l.scope.DNS != nil {
		rep.UpdateOption(dhcpv4.OptDomainName(l.scope.DNS.Zone))
		if len(l.scope.DNS.Search) > 0 {
			rep.UpdateOption(dhcpv4.OptDomainSearch(&rfc1035label.Labels{Labels: l.scope.DNS.Search}))
		}
	}
	if l.Hostname != "" {
		hostname := l.Hostname
		if l.scope.DNS != nil && l.scope.DNS.AddZoneInHostname {
			fqdn := strings.Join([]string{l.Hostname, l.scope.DNS.Zone}, ".")
			hostname = fqdn
		}
		rep.UpdateOption(dhcpv4.OptHostName(strings.TrimSuffix(hostname, ".")))
	}

	rep.ServerIPAddr = net.ParseIP(extconfig.Get().Instance.IP)
	rep.UpdateOption(dhcpv4.OptServerIdentifier(rep.ServerIPAddr))
	rep.YourIPAddr = net.ParseIP(l.Address)

	for _, opt := range l.scope.Options {
		finalVal := make([]byte, 0)
		if opt.Tag == nil && opt.TagName == "" {
			continue
		}
		if opt.TagName != "" {
			tag, ok := types.TagMap[types.OptionTagName(opt.TagName)]
			if !ok {
				req.log.Warn("invalid tag name", zap.String("tagName", opt.TagName))
				continue
			}
			opt.Tag = &tag
		}

		// Values which are directly converted from string to byte
		if opt.Value != nil {
			finalVal = []byte(*opt.Value)
			if _, ok := types.IPTags[*opt.Tag]; ok {
				i := net.ParseIP(*opt.Value)
				finalVal = dhcpv4.IPs([]net.IP{i}).ToBytes()
			}
		}

		// For non-stringable values, get b64 decoded values
		if len(opt.Value64) > 0 {
			values64 := make([]byte, 0)
			for _, v := range opt.Value64 {
				va, err := base64.StdEncoding.DecodeString(v)
				if err != nil {
					req.log.Warn("failed to convert base64 value to byte", zap.Error(err))
				} else {
					values64 = append(values64, va...)
				}
			}
			finalVal = values64
		}
		dopt := dhcpv4.OptGeneric(dhcpv4.GenericOptionCode(*opt.Tag), finalVal)
		rep.UpdateOption(dopt)
	}
	return rep
}
