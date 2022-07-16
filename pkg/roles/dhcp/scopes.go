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

	SubnetCIDR string            `json:"subnetCidr"`
	Default    bool              `json:"default"`
	Options    []*Option         `json:"options"`
	TTL        int64             `json:"ttl"`
	IPAM       map[string]string `json:"ipam"`
	DNS        struct {
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
	switch s.IPAM["type"] {
	case "internal":
	default:
		ipamInst, err = NewInternalIPAM(r, s)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create ipam: %w", err)
	}
	s.ipam = ipamInst
	return s, nil
}

func (r *Role) findScopeForRequest(req *Request) *Scope {
	var match *Scope
	longestBits := 0
	for _, scope := range r.scopes {
		ip := req.peer.(*net.UDPAddr).IP
		// Handle cases where peer is an actual IP (most likely relay)
		subBits := scope.match(ip, req)
		if subBits > -1 && subBits > longestBits {
			req.log.WithField("name", scope.Name).Trace("selected scope based on cidr match (peer IP)")
			match = scope
			longestBits = subBits
		}
		// Handle local broadcast, check with the instance's listening IP
		if match == nil {
			subBits := scope.match(net.ParseIP(extconfig.Get().Instance.IP), req)
			if subBits > -1 && subBits > longestBits {
				req.log.WithField("name", scope.Name).Trace("selected scope based on cidr match (instance IP)")
				match = scope
				longestBits = subBits
			}
		}
		// Handle default scope
		if match == nil && scope.Default {
			req.log.WithField("name", scope.Name).Debug("selected scope based on default state")
			match = scope
		}
	}
	return match
}

func (s *Scope) match(peer net.IP, req *Request) int {
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

func (s *Scope) createLeaseFor(req *Request) *Lease {
	ident := s.role.DeviceIdentifier(req.DHCPv4)
	lease := &Lease{
		Identifier: ident,

		Hostname: req.HostName(),
		ScopeKey: s.Name,

		inst:  s.inst,
		log:   req.log.WithField("lease", ident),
		scope: s,
	}
	requestedIP := req.RequestedIPAddress()
	if requestedIP != nil {
		s.log.WithField("ip", requestedIP).Debug("checking requested IP")
		ip, err := netip.ParseAddr(requestedIP.String())
		if err != nil {
			s.log.WithError(err).Warning("failed to parse requested ip")
		} else if s.ipam.IsIPFree(ip) {
			s.log.WithField("ip", requestedIP).Debug("requested IP is free")
			lease.Address = requestedIP.String()
		}
	}
	if lease.Address == "" {
		ip := s.ipam.NextFreeAddress()
		if ip == nil {
			return nil
		}
		lease.Address = s.ipam.NextFreeAddress().String()
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
