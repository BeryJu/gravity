package dhcp

import (
	"context"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"net"
	"net/netip"
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
	inst  roles.Instance
	scope *Scope
	log   *zap.Logger

	Identifier string `json:"-"`

	Address          string `json:"address"`
	Hostname         string `json:"hostname"`
	AddressLeaseTime string `json:"addressLeaseTime,omitempty"`
	ScopeKey         string `json:"scopeKey"`
	DNSZone          string `json:"dnsZone,omitempty"`
	// Set to -1 for a reservation
	Expiry      int64  `json:"expiry"`
	Description string `json:"description"`

	etcdKey string
}

func (r *Role) FindLease(req *Request4) *Lease {
	lease, ok := r.leases.GetPrefix(r.DeviceIdentifier(req.DHCPv4))
	if !ok {
		return nil
	}
	// Check if the leases's scope matches the expected scope to handle this request
	expectedScope := r.findScopeForRequest(req)
	if expectedScope != nil && lease.scope != expectedScope {
		// We have a specific scope to handle this request but it doesn't match the lease
		lease.scope = expectedScope
		lease.ScopeKey = expectedScope.Name
		lease.setLeaseIP(req)
		lease.log.Info("Re-assigning address for lease due to changed request scope", zap.String("newIP", lease.Address))
		err := lease.Put(req.Context, lease.scope.TTL)
		if err != nil {
			r.log.Warn("failed to update lease for re-assigned IP", zap.Error(err))
		}
	}
	return lease
}

func (r *Role) NewLease(identifier string) *Lease {
	return &Lease{
		inst:       r.i,
		Identifier: identifier,
		log:        r.log.With(zap.String("identifier", identifier)),
		Expiry:     0,
	}
}

func (l *Lease) setLeaseIP(req *Request4) {
	requestedIP := req.RequestedIPAddress()
	if requestedIP != nil {
		req.log.Debug("checking requested IP", zap.String("ip", requestedIP.String()))
		ip, _ := netip.AddrFromSlice(requestedIP)
		if l.scope.ipam.IsIPFree(ip, &l.Identifier) {
			req.log.Debug("requested IP is free", zap.String("ip", requestedIP.String()))
			l.Address = requestedIP.String()
			l.scope.ipam.UseIP(ip, l.Identifier)
			return
		}
	}
	ip := l.scope.ipam.NextFreeAddress(l.Identifier)
	if ip == nil {
		return
	}
	req.log.Debug("using next free IP from IPAM", zap.String("ip", ip.String()))
	l.Address = ip.String()
	l.scope.ipam.UseIP(*ip, l.Identifier)
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

	scope, ok := r.scopes.GetPrefix(l.ScopeKey)
	if !ok {
		return l, fmt.Errorf("DHCP lease with invalid scope key: %s", l.ScopeKey)
	}
	l.scope = scope
	return l, nil
}

func (l *Lease) IsReservation() bool {
	return l.Expiry == -1
}

func (l *Lease) Delete(ctx context.Context) error {
	leaseKey := l.inst.KV().Key(
		types.KeyRole,
		types.KeyLeases,
		l.Identifier,
	)
	_, err := l.inst.KV().Delete(
		ctx,
		leaseKey.String(),
	)
	return err
}

func (l *Lease) Put(ctx context.Context, expiry int64, opts ...clientv3.OpOption) error {
	if expiry > 0 && !l.IsReservation() {
		l.Expiry = time.Now().Add(time.Duration(expiry) * time.Second).Unix()

		exp, err := l.inst.KV().Grant(ctx, expiry)
		if err != nil {
			return err
		}
		opts = append(opts, clientv3.WithLease(exp.ID))
	}

	raw, err := json.Marshal(&l)
	if err != nil {
		return err
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
			"expiry":     expiry,
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
		} else if pl.Seconds() < 1 {
			req.log.Warn("invalid duration: less than 1", zap.String("duration", l.AddressLeaseTime))
		} else if pl.Seconds() > math.MaxInt32 {
			req.log.Warn("invalid duration: duration too long", zap.String("duration", l.AddressLeaseTime))
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

	// Check if the request has a different hostname, and update the lease
	if req.HostName() != l.Hostname {
		l.Hostname = req.HostName()
		// Update lease with new hostname
		err := l.Put(req.Context, l.Expiry)
		if err != nil {
			l.log.Warn("failed to update lease for updated hostname", zap.Error(err))
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
		if opt == nil || opt.Tag == nil && opt.TagName == "" {
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
					continue
				}
				values64 = append(values64, va...)
			}
			finalVal = values64
		}
		if len(opt.ValueHex) > 0 {
			valuesHex := make([]byte, 0)
			for _, v := range opt.ValueHex {
				va, err := hex.DecodeString(v)
				if err != nil {
					req.log.Warn("failed to convert hex value to byte", zap.Error(err))
					continue
				}
				valuesHex = append(valuesHex, va...)
			}
			finalVal = valuesHex
		}
		dopt := dhcpv4.OptGeneric(dhcpv4.GenericOptionCode(*opt.Tag), finalVal)
		rep.UpdateOption(dopt)
		if dopt.Code.Code() == uint8(dhcpv4.OptionBootfileName) {
			rep.BootFileName = dhcpv4.GetString(dopt.Code, rep.Options)
		}
	}
	return rep
}
