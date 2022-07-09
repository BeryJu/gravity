package discovery

import (
	"encoding/json"

	"beryju.io/ddet/pkg/roles"
	"beryju.io/ddet/pkg/roles/discovery/types"
	"github.com/Ullaakut/nmap/v2"
	log "github.com/sirupsen/logrus"
	"go.etcd.io/etcd/api/v3/mvccpb"
)

type Subnet struct {
	CIDR string `json:"cidr"`

	inst roles.Instance
	log  *log.Entry
}

func (r *DiscoveryRole) subnetFromKV(raw *mvccpb.KeyValue) (*Subnet, error) {
	sub := &Subnet{
		inst: r.i,
	}
	err := json.Unmarshal(raw.Value, &sub)
	if err != nil {
		return nil, err
	}
	sub.log = log.WithField("subnet", sub.CIDR)
	return sub, nil
}

func (s *Subnet) RunDiscovery() {
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
		nmap.WithServiceInfo(),
		nmap.WithSystemDNS(),
	)
	if err != nil {
		s.log.Fatalf("unable to create nmap scanner: %v", err)
		return
	}

	result, warnings, err := scanner.Run()
	if err != nil {
		s.log.Fatalf("unable to run nmap scan: %v", err)
		return
	}
	if warnings != nil {
		s.log.Printf("Warnings: \n %v", warnings)
	}

	// Use the results to print an example output
	for _, host := range result.Hosts {
		s.inst.DispatchEvent(types.EventTopicDiscoveryDeviceFound, roles.NewEvent(map[string]interface{}{
			"device": host,
		}))
		s.log.WithFields(log.Fields{
			"address":  host.Addresses,
			"hostname": host.Hostnames,
			"host":     host,
		}).Debug("found device")
	}
}
