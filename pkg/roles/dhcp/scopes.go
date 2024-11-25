package dhcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/netip"
	"strings"
	"time"

	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/dhcp/types"
	"github.com/insomniacslk/dhcp/dhcpv4"
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
	Hook       string              `json:"hook"`
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

	previous := r.scopes[s.Name]

	ipamInst, err := s.ipamType(previous)
	if err != nil {
		return nil, fmt.Errorf("failed to create ipam: %w", err)
	}
	s.ipam = ipamInst
	return s, nil
}

func (s *Scope) ipamType(previous *Scope) (IPAM, error) {
	if previous != nil && s.IPAM["type"] == previous.IPAM["type"] {
		err := previous.ipam.UpdateConfig(s)
		return previous.ipam, err
	}
	switch s.IPAM["type"] {
	case InternalIPAMType:
		fallthrough
	default:
		return NewInternalIPAM(s.role, s)
	}
}

func (r *Role) findScopeForRequest(req *Request4) *Scope {
	var match *Scope
	longestBits := 0
	r.scopesM.RLock()
	defer r.scopesM.RUnlock()
	// To prioritise requests from a DHCP relay being matched correctly, give their subnet
	// match a 1 bit more priority
	const dhcpRelayBias = 1
	const clientIPBias = 2
	for _, scope := range r.scopes {
		// Check based on Client IP Address (highest priority)
		clientIPMatchBits := scope.match(req.ClientIPAddr)
		if clientIPMatchBits > -1 && clientIPMatchBits+clientIPBias > longestBits {
			req.log.Debug("selected scope based on client IP", zap.String("scope", scope.Name))
			match = scope
			longestBits = clientIPMatchBits + clientIPBias
		}
		// Check based on gateway IP (next highest priority)
		gatewayMatchBits := scope.match(req.GatewayIPAddr)
		if gatewayMatchBits > -1 && gatewayMatchBits+dhcpRelayBias > longestBits {
			req.log.Debug("selected scope based on cidr match (gateway IP)", zap.String("scope", scope.Name))
			match = scope
			longestBits = gatewayMatchBits + dhcpRelayBias
		}
		// Handle local broadcast, check with the instance's listening IP
		// Only consider local scopes if we don't have a match already
		localMatchBits := scope.match(net.ParseIP(req.LocalIP()))
		if localMatchBits > -1 && localMatchBits > longestBits {
			req.log.Debug("selected scope based on cidr match (instance/interface IP)", zap.String("scope", scope.Name))
			match = scope
			longestBits = localMatchBits
		}
		// Fallback to default scope if we don't already have a match
		if match == nil && scope.Default {
			req.log.Debug("selected scope based on default flag", zap.String("scope", scope.Name))
			match = scope
		}
	}
	if match != nil {
		req.log.Debug("final scope selection", zap.String("scope", match.Name))
	}
	return match
}

func (s *Scope) match(peer net.IP) int {
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
	lease := s.role.NewLease(ident)
	lease.Hostname = req.HostName()

	lease.scope = s
	lease.ScopeKey = s.Name
	lease.setLeaseIP(req)
	req.log.Info("creating new DHCP lease", zap.String("ip", lease.Address), zap.String("identifier", ident))
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

func (s *Scope) executeHook(method string, args ...interface{}) {
	s.role.i.ExecuteHook(roles.HookOptions{
		Source: s.Hook,
		Method: method,
		Env: map[string]interface{}{
			"dhcp": map[string]interface{}{
				"Opt": func(code uint8, data []byte) dhcpv4.Option {
					return dhcpv4.OptGeneric(dhcpv4.GenericOptionCode(code), data)
				},
			},
		},
	}, args...)
}
