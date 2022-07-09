package discovery

import (
	"encoding/json"
	"fmt"

	"beryju.io/ddet/pkg/roles"
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
	scanner, err := nmap.NewScanner(
		nmap.WithTargets(s.CIDR),
		nmap.WithPingScan(),
		nmap.WithForcedDNSResolution(),
	)
	if err != nil {
		log.Fatalf("unable to create nmap scanner: %v", err)
	}

	result, warnings, err := scanner.Run()
	if err != nil {
		log.Fatalf("unable to run nmap scan: %v", err)
	}

	if warnings != nil {
		log.Printf("Warnings: \n %v", warnings)
	}

	// Use the results to print an example output
	for _, host := range result.Hosts {
		if len(host.Ports) == 0 || len(host.Addresses) == 0 {
			continue
		}

		fmt.Printf("Host %q:\n", host.Addresses[0])

		for _, port := range host.Ports {
			fmt.Printf("\tPort %d/%s %s %s\n", port.ID, port.Protocol, port.State, port.Service.Name)
		}
	}

	fmt.Printf("Nmap done: %d hosts up scanned in %3f seconds\n", len(result.Hosts), result.Stats.Finished.Elapsed)

}
