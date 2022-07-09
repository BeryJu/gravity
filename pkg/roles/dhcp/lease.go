package dhcp

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"net"
	"strings"
	"time"

	"beryju.io/ddet/pkg/extconfig"
	"beryju.io/ddet/pkg/roles"
	"beryju.io/ddet/pkg/roles/dhcp/types"
	"beryju.io/ddet/pkg/roles/dns/utils"
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
	AddressLeaseTime string `json:"addressLeaseTime"`
	ScopeKey         string `json:"scopeKey"`

	scope   *Scope
	etcdKey string
	inst    roles.Instance
	log     *log.Entry
}

func (r *DHCPRole) leaseFromKV(raw *mvccpb.KeyValue) (*Lease, error) {
	s := &Lease{
		inst: r.i,
	}
	err := json.Unmarshal(raw.Value, &s)
	if err != nil {
		return nil, err
	}
	prefix := r.i.KV().Key(
		types.KeyRole,
		types.KeyScopes,
		// l.Scope.Name,
		types.KeyLeases,
		"",
	)
	s.Identifier = strings.TrimPrefix(string(raw.Key), prefix)
	// Get full etcd key without leading slash since this usually gets passed to Instance Key()
	s.etcdKey = string(raw.Key)[1:]

	s.log = r.log.WithField("lease", prefix)
	return s, nil
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
		types.KeyScopes,
		l.scope.Name,
		types.KeyLeases,
		l.Identifier,
	)
	_, err = l.inst.KV().Put(
		context.TODO(),
		leaseKey,
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
	l.inst.DispatchEvent(types.EventTopicDHCPLeaseGiven, ev)

	l.log.WithField("expiry", expiry).Debug("put lease")
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

	ipLeaseDuration, err := time.ParseDuration(l.AddressLeaseTime)
	if err != nil {
		l.log.WithField("default", "24h").WithError(err).Warning("failed to parse address lease duration, defaulting")
		ipLeaseDuration = time.Hour * 24
	}
	rep.UpdateOption(dhcpv4.OptIPAddressLeaseTime(ipLeaseDuration))
	rep.UpdateOption(dhcpv4.OptSubnetMask(l.scope.ipam.GetSubnetMask()))

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
		if opt.Tag == nil {
			continue
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

	l.log.Trace(rep.Summary(), "peer", peer.String())
	if _, err := conn.WriteTo(rep.ToBytes(), peer); err != nil {
		l.log.WithError(err).Warning("failed to write reply")
	}
}
