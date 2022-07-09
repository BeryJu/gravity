package discovery

import (
	"encoding/json"
	"time"

	"beryju.io/ddet/pkg/roles"
	"beryju.io/ddet/pkg/roles/discovery/types"
	"github.com/Ullaakut/nmap/v2"
	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

type Subnet struct {
	CIDR         string `json:"cidr"`
	DiscoveryTTL int    `json:"discoveryTTL"`

	inst roles.Instance
	log  *log.Entry
}

func (r *DiscoveryRole) subnetFromKV(raw *mvccpb.KeyValue) (*Subnet, error) {
	sub := &Subnet{
		DiscoveryTTL: int((24 * time.Hour).Seconds()),
		inst:         r.i,
	}
	err := json.Unmarshal(raw.Value, &sub)
	if err != nil {
		return nil, err
	}
	sub.log = r.log.WithField("subnet", sub.CIDR)
	return sub, nil
}

func (s *Subnet) RunDiscovery() {
	s.log.Trace("Starting scan for subnet")
	s.inst.DispatchEvent(types.EventTopicDiscoveryStarted, roles.NewEvent(map[string]interface{}{
		"subnet": s,
	}))
	defer s.inst.DispatchEvent(types.EventTopicDiscoveryEnded, roles.NewEvent(map[string]interface{}{
		"subnet": s,
	}))

	scanner, err := nmap.NewScanner(
		nmap.WithTargets(s.CIDR),
		nmap.WithPingScan(),
		nmap.WithForcedDNSResolution(),
		nmap.WithSystemDNS(),
	)
	s.log.WithField("args", scanner.Args()).Trace("nmap args")
	if err != nil {
		s.log.Fatalf("unable to create nmap scanner: %v", err)
		return
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
		s.log.Fatalf("unable to run nmap scan: %v", err)
		return
	}
	if warnings != nil {
		s.log.Printf("Warnings: \n %v", warnings)
	}

	// Use the results to print an example output
	for _, host := range result.Hosts {
		dev := NewDevice(s.inst)
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
		err := dev.put(int64(s.DiscoveryTTL))
		if err != nil {
			s.log.WithError(err).Warning("ignoring device")
		}
	}
}
