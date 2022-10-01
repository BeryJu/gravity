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
	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Subnet struct {
	Identifier string `json:"-"`

	CIDR         string `json:"cidr"`
	DiscoveryTTL int    `json:"discoveryTTL"`

	etcdKey string
	inst    roles.Instance
	log     *log.Entry
	role    *Role
}

func (r *Role) NewSubnet(name string) *Subnet {
	return &Subnet{
		DiscoveryTTL: int((24 * time.Hour).Seconds()),
		inst:         r.i,
		Identifier:   name,
		log:          r.log.WithField("subnet", name),
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

func (s *Subnet) RunDiscovery() []Device {
	s.log.Trace("starting scan for subnet")
	s.inst.DispatchEvent(types.EventTopicDiscoveryStarted, roles.NewEvent(s.role.ctx, map[string]interface{}{
		"subnet": s,
	}))
	defer s.inst.DispatchEvent(types.EventTopicDiscoveryEnded, roles.NewEvent(s.role.ctx, map[string]interface{}{
		"subnet": s,
	}))

	scanner, err := nmap.NewScanner(
		nmap.WithTargets(s.CIDR),
		nmap.WithPingScan(),
		nmap.WithForcedDNSResolution(),
		nmap.WithCustomDNSServers(extconfig.Get().FallbackDNS),
	)
	s.log.WithField("args", scanner.Args()).Trace("nmap args")
	if err != nil {
		s.log.WithError(err).Warning("unable to create nmap scanner")
		return []Device{}
	}

	progress := make(chan float32, 1)

	// Function to listen and print the progress
	go func() {
		for p := range progress {
			s.log.WithField("progress", p).Debug("scan progress")
		}
	}()

	result, warnings, err := scanner.RunWithProgress(progress)
	if err != nil {
		s.log.WithError(err).Warning("unable to run nmap scan")
		return []Device{}
	}
	for _, warning := range warnings {
		s.log.Warning(warning)
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
		err := dev.put(s.role.ctx, int64(s.DiscoveryTTL))
		if err != nil {
			s.log.WithError(err).Warning("ignoring device")
		}
	}
	return devices
}

func (s *Subnet) put(opts ...clientv3.OpOption) error {
	key := s.inst.KV().Key(types.KeyRole, types.KeySubnets, s.Identifier)
	raw, err := json.Marshal(&s)
	if err != nil {
		return err
	}
	_, err = s.inst.KV().Put(
		context.Background(),
		key.String(),
		string(raw),
		opts...,
	)
	if err != nil {
		return err
	}
	return nil
}
