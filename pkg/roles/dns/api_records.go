package dns

import (
	"context"
	"errors"

	"beryju.io/gravity/pkg/roles/dns/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func (r *DNSRole) apiHandlerZoneRecords() usecase.Interactor {
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
		zone, ok := r.zones[in.Zone]
		if !ok {
			return status.Wrap(errors.New("not found"), status.NotFound)
		}
		rawRecords, err := r.i.KV().Get(ctx, r.i.KV().Key(
			types.KeyRole,
			types.KeyZones,
			zone.Name,
			"",
		))
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		for _, rec := range rawRecords.Kvs {
			rec, err := zone.recordFromKV(rec)
			if err != nil {
				r.log.WithError(err).Warning("failed to parse record")
				continue
			}
			out.Records = append(out.Records, record{
				FQDN: rec.Name,
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
