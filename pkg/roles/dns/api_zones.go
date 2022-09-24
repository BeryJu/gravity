package dns

import (
	"context"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func (r *Role) apiHandlerZonesGet() usecase.Interactor {
	type zone struct {
		Name           string              `json:"name" required:"true"`
		Authoritative  bool                `json:"authoritative" required:"true"`
		HandlerConfigs []map[string]string `json:"handlerConfigs" required:"true"`
		DefaultTTL     uint32              `json:"defaultTTL" required:"true"`
	}
	type zonesOutput struct {
		Zones []zone `json:"zones" required:"true"`
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

func (r *Role) apiHandlerZonesPut() usecase.Interactor {
	type zoneInput struct {
		Name           string              `query:"zone" required:"true"`
		Authoritative  bool                `json:"authoritative" required:"true"`
		HandlerConfigs []map[string]string `json:"handlerConfigs" required:"true"`
		DefaultTTL     uint32              `json:"defaultTTL" required:"true"`
	}
	u := usecase.NewIOI(new(zoneInput), new(struct{}), func(ctx context.Context, input, output interface{}) error {
		var (
			in = input.(*zoneInput)
		)
		z := r.newZone(in.Name)
		z.Name = in.Name
		z.Authoritative = in.Authoritative
		z.HandlerConfigs = in.HandlerConfigs
		z.DefaultTTL = in.DefaultTTL
		err := z.put(ctx)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("dns.put_zones")
	u.SetTitle("DNS Zones")
	u.SetTags("roles/dns")
	u.SetExpectedErrors(status.Internal, status.InvalidArgument)
	return u
}

func (r *Role) apiHandlerZonesDelete() usecase.Interactor {
	type zoneInput struct {
		Zone string `query:"zone"`
	}
	u := usecase.NewIOI(new(zoneInput), new(struct{}), func(ctx context.Context, input, output interface{}) error {
		var (
			in = input.(*zoneInput)
		)
		z, ok := r.zones[in.Zone]
		if !ok {
			return status.InvalidArgument
		}
		_, err := r.i.KV().Delete(ctx, z.etcdKey)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("dns.delete_zones")
	u.SetTitle("DNS Zones")
	u.SetTags("roles/dns")
	u.SetExpectedErrors(status.Internal, status.InvalidArgument)
	return u
}
