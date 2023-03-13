package api

import (
	"context"

	"beryju.io/gravity/pkg/extconfig"
	"github.com/Ullaakut/nmap/v2"
	"github.com/pkg/errors"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func (r *Role) Portmap(addr string) (nmap.Host, error) {
	scanner, err := nmap.NewScanner(
		nmap.WithTargets(addr),
		nmap.WithCustomDNSServers(extconfig.Get().FallbackDNS),
		nmap.WithForcedDNSResolution(),
	)
	if err != nil {
		return nmap.Host{}, err
	}
	result, warnings, err := scanner.Run()
	if err != nil {
		return nmap.Host{}, err
	}
	for _, warning := range warnings {
		r.log.Warn(warning)
	}
	if len(result.Hosts) < 1 {
		return nmap.Host{}, errors.New("no host responded")
	}
	return result.Hosts[0], nil
}

type APIToolPortmapInput struct {
	Host string `json:"host"`
}

type APIToolPortmapOutputPort struct {
	Protocol string `json:"protocol"`
	Name     string `json:"name"`
	Reason   string `json:"reason"`
	Port     uint16 `json:"port"`
}
type APIToolPortmapOutput struct {
	Ports []APIToolPortmapOutputPort `json:"ports"`
}

func (r *Role) APIToolPortmap() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIToolPortmapInput, output *APIToolPortmapOutput) error {
		device, err := r.Portmap(input.Host)
		if err != nil {
			return errors.Wrap(err, "failed to portmap")
		}
		for _, p := range device.Ports {
			ap := APIToolPortmapOutputPort{
				Protocol: p.Protocol,
				Name:     p.Service.Name,
				Reason:   p.State.String(),
				Port:     p.ID,
			}
			output.Ports = append(output.Ports, ap)
		}
		return nil
	})
	u.SetName("tools.portmap")
	u.SetTitle("Portmap tool")
	u.SetTags("roles/api")
	u.SetExpectedErrors(status.Internal)
	return u
}
