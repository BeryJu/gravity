package dns

import (
	"context"

	"github.com/swaggest/usecase"
)

func (r *DNSRole) apiHandlerZones() usecase.Interactor {
	type zone struct {
		Name           string              `json:"name"`
		Authoritative  bool                `json:"authoritative"`
		HandlerConfigs []map[string]string `json:"handlerConfigs"`
		DefaultTTL     uint32              `json:"defaultTTL"`
	}
	type zonesOutput struct {
		Zones []zone `json:"zones"`
	}
	u := usecase.NewIOI(new(struct{}), new(zonesOutput), func(ctx context.Context, input, output interface{}) error {
		var (
			out = output.(*zonesOutput)
		)
		for name, _zone := range r.zones {
			out.Zones = append(out.Zones, zone{
				Name:          name,
				Authoritative: _zone.Authoritative,
			})
		}
		return nil
	})
	u.SetName("dns.get_zones")
	u.SetTitle("DNS Zones")
	u.SetTags("roles/dns")
	return u
}
