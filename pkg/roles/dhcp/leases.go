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
	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Lease struct {
	Identifier string `json:"-"`

	Address          string `json:"address"`
	Hostname         string `json:"hostname"`
	AddressLeaseTime string `json:"addressLeaseTime,omitempty"`
	ScopeKey         string `json:"scopeKey"`

	scope   *Scope
	etcdKey string
	inst    roles.Instance
	log     *log.Entry
}

func (r *Role) newLease(identifier string) *Lease {
	return &Lease{
		inst:       r.i,
		Identifier: identifier,
		log:        r.log.WithField("identifier", identifier),
	}
}

func (r *Role) leaseFromKV(raw *mvccpb.KeyValue) (*Lease, error) {
	prefix := r.i.KV().Key(
		types.KeyRole,
		types.KeyLeases,
	).Prefix(true).String()
	identifier := strings.TrimPrefix(string(raw.Key), prefix)
	l := r.newLease(identifier)
	err := json.Unmarshal(raw.Value, &l)
	if err != nil {
		return nil, err
	}
	l.etcdKey = string(raw.Key)

	scope, ok := r.scopes[l.ScopeKey]
	if !ok {
		return nil, fmt.Errorf("DHCP lease with invalid scope key: %s", l.ScopeKey)
	}
	l.scope = scope
	return l, nil
}

func (l *Lease) put(expiry int64, opts ...clientv3.OpOption) error {
	raw, err := json.Marshal(&l)
	if err != nil {
		return err
	}

	if expiry > 0 {
		exp, err := l.inst.KV().Lease.Grant(context.TODO(), expiry)
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
		context.TODO(),
		leaseKey.String(),
		string(raw),
		opts...,
	)
	if err != nil {
		return err
	}
	ev := roles.NewEvent(
		map[string]interface{}{
			"hostname": l.Hostname,
			"address":  l.Address,
			"fqdn":     utils.EnsureTrailingPeriod(strings.Join([]string{l.Hostname, l.scope.DNS.Zone}, ".")),
		},
	)
	ev.Payload.RelatedObjectKey = leaseKey
	ev.Payload.RelatedObjectOptions = opts
	l.inst.DispatchEvent(types.EventTopicDHCPLeasePut, ev)

	l.log.WithField("expiry", expiry).WithField("identifier", l.Identifier).Debug("put lease")
	return nil
}

func (l *Lease) reply(
	conn net.PacketConn,
	peer net.Addr,
	m *dhcpv4.DHCPv4,
	modifyResponse func(*dhcpv4.DHCPv4) *dhcpv4.DHCPv4,
) {
	rep, err := dhcpv4.NewReplyFromRequest(m)
	if err != nil {
		l.log.WithError(err).Warning("failed to create reply")
		return
	}

	rep = modifyResponse(rep)
	rep.UpdateOption(dhcpv4.OptSubnetMask(l.scope.ipam.GetSubnetMask()))

	if l.AddressLeaseTime != "" {
		pl, err := time.ParseDuration(l.AddressLeaseTime)
		if err != nil {
			l.log.WithField("default", pl.String()).WithError(err).Warning("failed to parse address lease duration, defaulting")
		} else {
			rep.UpdateOption(dhcpv4.OptIPAddressLeaseTime(pl))
		}
	} else {
		rep.UpdateOption(dhcpv4.OptIPAddressLeaseTime(time.Duration(l.scope.TTL * int64(time.Second))))
	}

	// DNS Options
	rep.UpdateOption(dhcpv4.OptDNS(net.ParseIP(extconfig.Get().Instance.IP)))
	rep.UpdateOption(dhcpv4.OptDomainName(l.scope.DNS.Zone))
	if len(l.scope.DNS.Search) > 0 {
		rep.UpdateOption(dhcpv4.OptDomainSearch(&rfc1035label.Labels{Labels: l.scope.DNS.Search}))
	}
	if l.scope.DNS.AddZoneInHostname {
		fqdn := strings.Join([]string{l.Hostname, l.scope.DNS.Zone}, ".")
		rep.UpdateOption(dhcpv4.OptHostName(fqdn))
	}

	rep.YourIPAddr = net.ParseIP(l.Address)
	rep.UpdateOption(dhcpv4.OptHostName(l.Hostname))

	for _, opt := range l.scope.Options {
		finalVal := make([]byte, 0)
		if opt.Tag == nil && opt.TagName == "" {
			continue
		}
		if opt.TagName != "" {
			tag, ok := TagMap[opt.TagName]
			if !ok {
				l.log.WithError(err).Warningf("invalid tag name %s", opt.TagName)
				continue
			}
			opt.Tag = &tag
		}

		// Values which are directly converted from string to byte
		if opt.Value != nil {
			i := net.ParseIP(*opt.Value)
			if i == nil {
				finalVal = []byte(*opt.Value)
			} else {
				finalVal = dhcpv4.IPs([]net.IP{i}).ToBytes()
			}
		}

		// For non-stringable values, get b64 decoded values
		if len(opt.Value64) > 0 {
			values64 := make([]byte, 0)
			for _, v := range opt.Value64 {
				va, err := base64.StdEncoding.DecodeString(v)
				if err != nil {
					l.log.WithError(err).Warning("failed to convert base64 value to byte")
				} else {
					values64 = append(values64, va...)
				}
			}
			finalVal = values64
		}
		dopt := dhcpv4.OptGeneric(dhcpv4.GenericOptionCode(*opt.Tag), finalVal)
		rep.UpdateOption(dopt)
	}

	l.log.Trace(rep.Summary())
	if _, err := conn.WriteTo(rep.ToBytes(), peer); err != nil {
		l.log.WithError(err).Warning("failed to write reply")
	}
}
