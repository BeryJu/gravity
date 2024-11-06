package dhcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/netip"
	"strings"
	"sync"
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
	role *Role
	log  *zap.Logger

	cidr netip.Prefix
	Name string `json:"-"`

	etcdKey        string
	leases         map[string]*Lease
	leasesWatchCtx context.CancelFunc
	leasesSync     sync.RWMutex

	DNS        *ScopeDNS           `json:"dns"`
	IPAM       map[string]string   `json:"ipam"`
	SubnetCIDR string              `json:"subnetCidr"`
	Options    []*types.DHCPOption `json:"options"`
	TTL        int64               `json:"ttl"`
	Default    bool                `json:"default"`
	Hook       string              `json:"hook"`
}

func (r *Role) NewScope(name string) *Scope {
	return &Scope{
		Name:       name,
		inst:       r.i,
		role:       r,
		TTL:        int64((7 * 24 * time.Hour).Seconds()),
		log:        r.log.With(zap.String("scope", name)),
		DNS:        &ScopeDNS{},
		IPAM:       make(map[string]string),
		leases:     make(map[string]*Lease),
		leasesSync: sync.RWMutex{},
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
	lease := s.NewLease(ident)
	lease.Hostname = req.HostName()
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

func (s *Scope) watchScopeLeases(ctx context.Context) {
	evtHandler := func(ev *clientv3.Event) {
		lease, err := s.leaseFromKV(ev.Kv)
		defer s.calculateUsage()
		if ev.Type == clientv3.EventTypeDelete {
			delete(s.leases, lease.Identifier)
		} else {
			// Check if the record parsed above actually was parsed correctly,
			// we don't care for that when removing, but prevent adding
			// empty leases
			if err != nil {
				return
			}
			s.leasesSync.Lock()
			defer s.leasesSync.Unlock()
			s.leases[lease.Identifier] = lease
		}
	}
	ctx, canc := context.WithCancel(ctx)
	s.leasesWatchCtx = canc

	prefix := s.inst.KV().Key(s.etcdKey).Prefix(true).String()

	leases, err := s.inst.KV().Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		s.log.Warn("failed to list initial leases", zap.Error(err))
		time.Sleep(5 * time.Second)
		s.watchScopeLeases(ctx)
		return
	}
	for _, lease := range leases.Kvs {
		evtHandler(&clientv3.Event{
			Type: mvccpb.PUT,
			Kv:   lease,
		})
	}

	watchChan := s.inst.KV().Watch(
		ctx,
		prefix,
		clientv3.WithPrefix(),
	)
	go func() {
		for watchResp := range watchChan {
			for _, event := range watchResp.Events {
				go evtHandler(event)
			}
		}
	}()
}

func (s *Scope) StopWatchingLeases() {
	if s != nil && s.leasesWatchCtx != nil {
		s.leasesWatchCtx()
	}
}
