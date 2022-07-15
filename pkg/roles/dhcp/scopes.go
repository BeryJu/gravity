package dhcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/netip"
	"strings"
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/dhcp/types"
	"github.com/insomniacslk/dhcp/dhcpv4"
	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Option struct {
	Tag     *uint8   `json:"tag"`
	TagName string   `json:"tagName"`
	Value   *string  `json:"value"`
	Value64 []string `json:"value64"`
}

var TagMap map[string]uint8 = map[string]uint8{
	"subnet_mask": dhcpv4.OptionSubnetMask.Code(),
	"router":      dhcpv4.OptionRouter.Code(),
	"time_server": dhcpv4.OptionTimeServer.Code(),
	"name_server": dhcpv4.OptionNameServer.Code(),
	"domain_name": dhcpv4.OptionDomainName.Code(),
	"bootfile":    dhcpv4.OptionBootfileName.Code(),
	"tftp_server": dhcpv4.OptionTFTPServerName.Code(),
}

type Scope struct {
	Name string `json:"-"`

	SubnetCIDR string    `json:"subnetCidr"`
	Default    bool      `json:"default"`
	Options    []*Option `json:"options"`
	TTL        int64     `json:"ttl"`
	Range      struct {
		Start string `json:"start"`
		End   string `json:"end"`
	} `json:"range"`
	DNS struct {
		Zone              string   `json:"zone"`
		Search            []string `json:"search"`
		AddZoneInHostname bool     `json:"addZoneInHostname"`
	} `json:"dns"`

	cidr    netip.Prefix
	etcdKey string
	ipam    IPAM
	inst    roles.Instance
	role    *Role
	log     *log.Entry
}

func (r *Role) newScope(name string) *Scope {
	return &Scope{
		Name: name,
		inst: r.i,
		role: r,
		TTL:  int64((7 * 24 * time.Hour).Seconds()),
		log:  r.log.WithField("scope", name),
	}
}

func (r *Role) scopeFromKV(raw *mvccpb.KeyValue) (*Scope, error) {
	prefix := r.i.KV().Key(types.KeyRole, types.KeyScopes).Prefix(true).String()
	name := strings.TrimPrefix(string(raw.Key), prefix)

	s := r.newScope(name)
	err := json.Unmarshal(raw.Value, &s)
	if err != nil {
		return nil, err
	}
	cidr, err := netip.ParsePrefix(s.SubnetCIDR)
	if err != nil {
		return nil, err
	}
	s.cidr = cidr

	s.etcdKey = string(raw.Key)

	var ipamInst IPAM
	ipamInst, err = NewInternalIPAM(r, s.SubnetCIDR, s.Range.Start, s.Range.End)
	if err != nil {
		return nil, fmt.Errorf("failed to create ipam: %w", err)
	}
	s.ipam = ipamInst
	return s, nil
}

func (r *Role) findScopeForRequest(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) *Scope {
	var match *Scope
	longestBits := 0
	for _, scope := range r.scopes {
		ip := peer.(*net.UDPAddr).IP
		// Handle cases where peer is an actual IP (most likely relay)
		subBits := scope.match(conn, ip, m)
		if subBits > -1 && subBits > longestBits {
			r.log.WithField("name", scope.Name).Trace("selected scope based on cidr match (peer IP)")
			match = scope
			longestBits = subBits
		}
		// Handle local broadcast, check with the instance's listening IP
		if match == nil && peer.(*net.UDPAddr).IP.Equal(net.ParseIP("255.255.255.255")) {
			subBits := scope.match(conn, net.ParseIP(extconfig.Get().Instance.IP), m)
			if subBits > -1 && subBits > longestBits {
				r.log.WithField("name", scope.Name).Trace("selected scope based on cidr match (instance IP)")
				match = scope
				longestBits = subBits
			}
		}
		// Handle default scope
		if match == nil && scope.Default {
			r.log.WithField("name", scope.Name).Debug("selected scope based on default state")
			match = scope
		}
	}
	if match != nil {
		r.log.Trace("found scope for request")
	}
	return match
}

func (s *Scope) match(conn net.PacketConn, peer net.IP, m *dhcpv4.DHCPv4) int {
	ip, err := netip.ParseAddr(peer.String())
	if err != nil {
		s.log.WithError(err).Warning("failed to parse client ip")
		return -1
	}
	if s.cidr.Contains(ip) {
		return s.cidr.Bits()
	}
	return -1
}

func (s *Scope) createLeaseFor(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) *Lease {
	ident := s.role.DeviceIdentifier(m)
	lease := &Lease{
		Identifier: ident,

		Hostname: m.HostName(),
		Address:  s.ipam.NextFreeAddress().String(),
		ScopeKey: s.Name,

		inst:  s.inst,
		log:   s.log.WithField("lease", ident),
		scope: s,
	}
	if requestIp, ok := netip.AddrFromSlice(m.Options.Get(dhcpv4.OptionRequestedIPAddress)); ok {
		s.log.WithField("ip", requestIp).Debug("checking requested IP")
		if s.ipam.IsIPFree(requestIp) {
			s.log.WithField("ip", requestIp).Debug("requested IP is free")
			lease.Address = requestIp.String()
		}
	}
	return lease
}

func (s *Scope) put(expiry int64, opts ...clientv3.OpOption) error {
	raw, err := json.Marshal(&s)
	if err != nil {
		return err
	}

	if expiry > 0 {
		exp, err := s.inst.KV().Lease.Grant(context.TODO(), expiry)
		if err != nil {
			return err
		}
		opts = append(opts, clientv3.WithLease(exp.ID))
	}

	leaseKey := s.inst.KV().Key(
		types.KeyRole,
		types.KeyLeases,
		s.Name,
	)
	_, err = s.inst.KV().Put(
		context.TODO(),
		leaseKey.String(),
		string(raw),
		opts...,
	)
	if err != nil {
		return err
	}
	return nil
}
