package dns

import (
	"context"
	"errors"

	"beryju.io/gravity/pkg/roles/dns/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

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
