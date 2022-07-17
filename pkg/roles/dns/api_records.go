package dns

import (
	"context"
	"errors"

	"beryju.io/gravity/pkg/roles/dns/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func (r *Role) apiHandlerZoneRecords() usecase.Interactor {
	type zoneRecordsInput struct {
		Zone string `path:"zone"`
	}
	type record struct {
		FQDN     string `json:"fqdn"`
		Hostname string `json:"hostname"`
		Type     string `json:"type"`
		Data     string `json:"data"`

		MXPreference uint16 `json:"mxPreference,omitempty"`
		SRVPort      uint16 `json:"srvPort,omitempty"`
		SRVPriority  uint16 `json:"srvPriority,omitempty"`
		SRVWeight    uint16 `json:"srvWeight,omitempty"`
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
		).Prefix(true).String())
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
				Hostname:     rec.Name,
				FQDN:         rec.Name + "." + zone.Name,
				Type:         rec.Type,
				Data:         rec.Data,
				MXPreference: rec.MXPreference,
				SRVPort:      rec.SRVPort,
				SRVPriority:  rec.SRVPriority,
				SRVWeight:    rec.SRVWeight,
			})
		}
		return nil
	})
	u.SetName("dns.get_records")
	u.SetTitle("DNS Records")
	u.SetTags("roles/dns")
	u.SetExpectedErrors(status.InvalidArgument, status.NotFound, status.Internal)
	return u
}
