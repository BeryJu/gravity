package dns

import (
	"context"
	"errors"

	"beryju.io/gravity/pkg/roles/dns/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (r *Role) apiHandlerZoneRecordsGet() usecase.Interactor {
	type recordsInput struct {
		Zone string `query:"zone"`
	}
	type record struct {
		UID      string `json:"uid" required:"true"`
		FQDN     string `json:"fqdn" required:"true"`
		Hostname string `json:"hostname" required:"true"`
		Type     string `json:"type" required:"true"`

		Data         string `json:"data" required:"true"`
		MXPreference uint16 `json:"mxPreference,omitempty"`
		SRVPort      uint16 `json:"srvPort,omitempty"`
		SRVPriority  uint16 `json:"srvPriority,omitempty"`
		SRVWeight    uint16 `json:"srvWeight,omitempty"`
	}
	type recordsOutput struct {
		Records []record `json:"records" required:"true"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input recordsInput, output *recordsOutput) error {
		zone, ok := r.zones[input.Zone]
		if !ok {
			return status.Wrap(errors.New("not found"), status.NotFound)
		}
		rawRecords, err := r.i.KV().Get(ctx, r.i.KV().Key(
			types.KeyRole,
			types.KeyZones,
			zone.Name,
		).Prefix(true).String(), clientv3.WithPrefix())
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		for _, rec := range rawRecords.Kvs {
			rec, err := zone.recordFromKV(rec)
			if err != nil {
				r.log.WithError(err).Warning("failed to parse record")
				continue
			}
			output.Records = append(output.Records, record{
				UID:          rec.uid,
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

func (r *Role) apiHandlerZoneRecordsPut() usecase.Interactor {
	type recordsInput struct {
		Zone     string `query:"zone" required:"true" maxLength:"255"`
		Hostname string `query:"hostname" required:"true" maxLength:"255"`
		UID      string `query:"uid" maxLength:"255"`

		Type string `json:"type" required:"true"`

		Data         string `json:"data" required:"true"`
		MXPreference uint16 `json:"mxPreference,omitempty"`
		SRVPort      uint16 `json:"srvPort,omitempty"`
		SRVPriority  uint16 `json:"srvPriority,omitempty"`
		SRVWeight    uint16 `json:"srvWeight,omitempty"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input recordsInput, output *struct{}) error {
		zone, ok := r.zones[input.Zone]
		if !ok {
			return status.Wrap(errors.New("zone not found"), status.NotFound)
		}
		rec := zone.newRecord(input.Hostname, input.Type)
		rec.uid = input.UID
		rec.Data = input.Data
		rec.MXPreference = input.MXPreference
		rec.SRVPort = input.SRVPort
		rec.SRVPriority = input.SRVPriority
		rec.SRVWeight = input.SRVWeight
		err := rec.put(ctx, -1)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("dns.put_records")
	u.SetTitle("DNS Records")
	u.SetTags("roles/dns")
	u.SetExpectedErrors(status.Internal, status.InvalidArgument, status.NotFound)
	return u
}

func (r *Role) apiHandlerZoneRecordsDelete() usecase.Interactor {
	type recordsDeleteInput struct {
		Zone     string `query:"zone" required:"true"`
		Hostname string `query:"hostname" required:"true"`
		UID      string `query:"uid" required:"true"`
		Type     string `query:"type" required:"true"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input recordsDeleteInput, output *struct{}) error {
		_, ok := r.zones[input.Zone]
		if !ok {
			return status.Wrap(errors.New("zone not found"), status.NotFound)
		}
		_, err := r.i.KV().Delete(
			ctx,
			r.i.KV().Key(
				types.KeyRole,
				types.KeyZones,
				input.Zone,
				input.Hostname,
				input.Type,
				input.UID,
			).String(),
		)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("dns.delete_records")
	u.SetTitle("DNS Records")
	u.SetTags("roles/dns")
	u.SetExpectedErrors(status.Internal, status.InvalidArgument, status.NotFound)
	return u
}
