package dhcp

import (
	"context"
	"encoding/json"
	"fmt"
	"math/big"
	"net"
	"net/netip"
	"strings"
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/dhcp/types"
	"github.com/netdata/go.d.plugin/pkg/iprange"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type ScopeDNS struct {
	Zone              string   `json:"zone"`
	Search            []string `json:"search"`
	AddZoneInHostname bool     `json:"addZoneInHostname"`
}

type Scope struct {
	ipam IPAM
	inst roles.Instance
	DNS  *ScopeDNS `json:"dns"`

	IPAM map[string]string `json:"ipam"`
	role *Role
	log  *zap.Logger

	cidr netip.Prefix
	Name string `json:"-"`

	etcdKey string

	SubnetCIDR string              `json:"subnetCidr"`
	Options    []*types.DHCPOption `json:"options"`
	TTL        int64               `json:"ttl"`
	Default    bool                `json:"default"`
}

func (r *Role) NewScope(name string) *Scope {
	return &Scope{
		Name: name,
		inst: r.i,
		role: r,
		TTL:  int64((7 * 24 * time.Hour).Seconds()),
		log:  r.log.With(zap.String("scope", name)),
		DNS:  &ScopeDNS{},
		IPAM: make(map[string]string),
	}
}

func (r *Role) scopeFromKV(raw *mvccpb.KeyValue) (*Scope, error) {
	prefix := r.i.KV().Key(types.KeyRole, types.KeyScopes).Prefix(true).String()
	name := strings.TrimPrefix(string(raw.Key), prefix)

	s := r.NewScope(name)
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
	case InternalIPAMType:
		fallthrough
	default:
		ipamInst, err = NewInternalIPAM(r, s)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to create ipam: %w", err)
	}
	s.ipam = ipamInst
	return s, nil
}

func (r *Role) findScopeForRequest(req *Request4) *Scope {
	var match *Scope
	longestBits := 0
	r.scopesM.RLock()
	defer r.scopesM.RUnlock()
	for _, scope := range r.scopes {
		// Handle local broadcast, check with the instance's listening IP
		if match == nil {
			subBits := scope.match(net.ParseIP(extconfig.Get().Instance.IP), req)
			if subBits > -1 && subBits > longestBits {
				req.log.Debug("selected scope based on cidr match (instance IP)", zap.String("scope", scope.Name))
				match = scope
				longestBits = subBits
			}
		}
		// Check based on gateway IP
		subBits := scope.match(req.GatewayIPAddr, req)
		if subBits > -1 && subBits > longestBits {
			req.log.Debug("selected scope based on cidr match (peer IP)", zap.String("scope", scope.Name))
			match = scope
			longestBits = subBits
		}
		// Handle default scope
		if match == nil && scope.Default {
			req.log.Debug("selected scope based on default flag", zap.String("scope", scope.Name))
			match = scope
		}
	}
	return match
}

func (s *Scope) match(peer net.IP, req *Request4) int {
	ip, err := netip.ParseAddr(peer.String())
	if err != nil {
		s.log.Warn("failed to parse client ip", zap.Error(err))
		return -1
	}
	if s.cidr.Contains(ip) {
		return s.cidr.Bits()
	}
	return -1
}

func (s *Scope) createLeaseFor(req *Request4) *Lease {
	ident := s.role.DeviceIdentifier(req.DHCPv4)
	lease := &Lease{
		Identifier: ident,

		Hostname: req.HostName(),
		ScopeKey: s.Name,

		inst:  s.inst,
		log:   req.log.With(zap.String("lease", ident)),
		scope: s,
	}
	requestedIP := req.RequestedIPAddress()
	if requestedIP != nil {
		req.log.Debug("checking requested IP", zap.String("ip", requestedIP.String()))
		ip, err := netip.ParseAddr(requestedIP.String())
		if err != nil {
			req.log.Warn("failed to parse requested ip", zap.Error(err))
		} else if s.ipam.IsIPFree(ip) {
			req.log.Debug("requested IP is free", zap.String("ip", requestedIP.String()))
			lease.Address = requestedIP.String()
		}
	}
	if lease.Address == "" {
		ip := s.ipam.NextFreeAddress()
		if ip == nil {
			return nil
		}
		lease.Address = ip.String()
	}
	req.log.Info("creating new DHCP lease", zap.String("ip", requestedIP.String()), zap.String("lease", ident))
	return lease
}

func (s *Scope) Put(ctx context.Context, expiry int64, opts ...clientv3.OpOption) error {
	raw, err := json.Marshal(&s)
	if err != nil {
		return err
	}

	if expiry > 0 {
		exp, err := s.inst.KV().Lease.Grant(ctx, expiry)
		if err != nil {
			return err
		}
		opts = append(opts, clientv3.WithLease(exp.ID))
	}

	leaseKey := s.inst.KV().Key(
		types.KeyRole,
		types.KeyScopes,
		s.Name,
	)
	_, err = s.inst.KV().Put(
		ctx,
		leaseKey.String(),
		string(raw),
		opts...,
	)
	if err != nil {
		return err
	}
	return nil
}

// calculateUsage Calculate scope usage for prometheus metrics
func (s *Scope) calculateUsage() {
	if s.IPAM["type"] != InternalIPAMType {
		return
	}
	ii := s.ipam.(*InternalIPAM)
	ips := iprange.New(ii.Start.AsSlice(), ii.End.AsSlice())
	usable := ips.Size()
	dhcpScopeSize.WithLabelValues(s.Name).Set(float64(usable.Uint64()))
	used := big.NewInt(0)
	for _, lease := range s.role.leases {
		if lease.ScopeKey != s.Name {
			continue
		}
		used = used.Add(used, big.NewInt(1))
	}
	dhcpScopeUsage.WithLabelValues(s.Name).Set(float64(used.Uint64()))
}
