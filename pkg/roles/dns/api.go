package dns

import (
	"context"
	"errors"

	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/dns/types"
	"github.com/swaggest/rest/web"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func (ro *DNSRole) apiHandlerZones() usecase.Interactor {
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
		for name, _zone := range ro.zones {
			out.Zones = append(out.Zones, zone{
				Name:          name,
				Authoritative: _zone.Authoritative,
			})
		}
		return nil
	})
	u.SetTitle("DNS Zones")
	u.SetTags("dns")
	u.SetDescription("List all DNS Zones.")
	return u
}

func (ro *DNSRole) apiHandlerZoneRecords() usecase.Interactor {
	type zoneRecordsInput struct {
		Zone string `path:"zone"`
	}
	type record struct {
		FQDN string `json:"fqdn"`
	}
	type zoneRecordsOutput struct {
		Records []record `json:"records"`
	}
	u := usecase.NewIOI(new(zoneRecordsInput), new(zoneRecordsOutput), func(ctx context.Context, input, output interface{}) error {
		var (
			in  = input.(*zoneRecordsInput)
			out = output.(*zoneRecordsOutput)
		)
		zone, ok := ro.zones[in.Zone]
		if !ok {
			return status.Wrap(errors.New("not found"), status.NotFound)
		}
		rawRecords, err := ro.i.KV().Get(ctx, ro.i.KV().Key(
			types.KeyRole,
			types.KeyZones,
			zone.Name,
			"",
		))
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		for _, rec := range rawRecords.Kvs {
			r := zone.recordFromKV(rec)
			out.Records = append(out.Records, record{
				FQDN: r.Name,
			})
		}
		return nil
	})
	u.SetTitle("DNS Records")
	u.SetTags("dns")
	u.SetDescription("List all DNS Records within a zone.")
	u.SetExpectedErrors(status.InvalidArgument, status.NotFound, status.Internal)
	return u
}

func (ro *DNSRole) eventHandlerAPIMux(ev *roles.Event) {
	svc := ev.Payload.Data["svc"].(*web.Service)
	svc.Get("/api/v1/dns/zones", ro.apiHandlerZones())
	svc.Get("/api/v1/dns/zones/{zone}/records", ro.apiHandlerZoneRecords())
}
