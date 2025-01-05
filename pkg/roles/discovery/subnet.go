package discovery

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/discovery/types"
	"github.com/Ullaakut/nmap/v2"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/concurrency"
	"go.uber.org/zap"
)

type Subnet struct {
	inst       roles.Instance
	log        *zap.Logger
	role       *Role
	Identifier string `json:"-"`

	CIDR        string `json:"cidr"`
	DNSResolver string `json:"dnsResolver"`

	etcdKey      string
	DiscoveryTTL int `json:"discoveryTTL"`
}

func (r *Role) NewSubnet(name string) *Subnet {
	return &Subnet{
		DiscoveryTTL: int((24 * time.Hour).Seconds()),
		inst:         r.i,
		Identifier:   name,
		log:          r.log.With(zap.String("subnet", name)),
		role:         r,
	}
}

func (r *Role) subnetFromKV(raw *mvccpb.KeyValue) (*Subnet, error) {
	prefix := r.i.KV().Key(types.KeyRole, types.KeySubnets).Prefix(true).String()
	name := strings.TrimPrefix(string(raw.Key), prefix)

	sub := r.NewSubnet(name)
	err := json.Unmarshal(raw.Value, &sub)
	if err != nil {
		return nil, err
	}
	sub.etcdKey = string(raw.Key)
	return sub, nil
}

func (s *Subnet) RunDiscovery(ctx context.Context) []Device {
	dev := []Device{}
	se, err := concurrency.NewSession(s.inst.KV().Client)
	if err != nil {
		s.log.Warn("Failed to create concurrency session", zap.Error(err))
		return dev
	}

	m := concurrency.NewMutex(
		se,
		s.role.i.KV().Key(types.KeyRole, types.KeySubnets, s.Identifier).String(),
	)
	m.Lock(ctx)
	defer m.Unlock(ctx)

	s.log.Debug("starting scan for subnet")
	s.inst.DispatchEvent(types.EventTopicDiscoveryStarted, roles.NewEvent(s.role.ctx, map[string]interface{}{
		"subnet": s,
	}))
	defer s.inst.DispatchEvent(types.EventTopicDiscoveryEnded, roles.NewEvent(s.role.ctx, map[string]interface{}{
		"subnet": s,
	}))

	dns := s.DNSResolver
	if dns == "" {
		dns = extconfig.Get().FallbackDNS
	}

	scanner, err := nmap.NewScanner(
		nmap.WithContext(ctx),
		nmap.WithTargets(s.CIDR),
		nmap.WithPingScan(),
		nmap.WithForcedDNSResolution(),
		nmap.WithCustomDNSServers(dns),
	)
	s.log.Debug("nmap args", zap.Strings("args", scanner.Args()))
	if err != nil {
		s.log.Warn("unable to create nmap scanner", zap.Error(err))
		return dev
	}

	result, warnings, err := scanner.Run()
	if err != nil {
		s.log.Warn("unable to run nmap scan", zap.Error(err))
		return dev
	}
	for _, warning := range warnings {
		s.log.Warn(warning)
	}

	devices := []Device{}
	for _, host := range result.Hosts {
		dev := s.role.newDevice()
		if len(host.Hostnames) > 0 {
			dev.Hostname = host.Hostnames[0].String()
		}
		for _, addr := range host.Addresses {
			if addr.AddrType == "mac" {
				dev.MAC = addr.Addr
			} else {
				dev.IP = addr.Addr
			}
		}
		devices = append(devices, *dev)
		err := dev.put(ctx, int64(s.DiscoveryTTL))
		if err != nil {
			s.log.Warn("ignoring device", zap.Error(err))
		}
	}
	return devices
}

func (s *Subnet) put(ctx context.Context, opts ...clientv3.OpOption) error {
	key := s.inst.KV().Key(types.KeyRole, types.KeySubnets, s.Identifier)
	raw, err := json.Marshal(&s)
	if err != nil {
		return err
	}
	_, err = s.inst.KV().Put(
		ctx,
		key.String(),
		string(raw),
		opts...,
	)
	if err != nil {
		return err
	}
	return nil
}
