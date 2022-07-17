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
		Zone string `path:"zone"`
	}
	type record struct {
		UID      string `json:"uid"`
		FQDN     string `json:"fqdn"`
		Hostname string `json:"hostname"`
		Type     string `json:"type"`

		Data         string `json:"data"`
		MXPreference uint16 `json:"mxPreference,omitempty"`
		SRVPort      uint16 `json:"srvPort,omitempty"`
		SRVPriority  uint16 `json:"srvPriority,omitempty"`
		SRVWeight    uint16 `json:"srvWeight,omitempty"`
	}
	type recordsOutput struct {
		Records []record `json:"records"`
	}
	u := usecase.NewIOI(new(recordsInput), new(recordsOutput), func(ctx context.Context, input, output interface{}) error {
		var (
			in  = input.(*recordsInput)
			out = output.(*recordsOutput)
		)
		zone, ok := r.zones[in.Zone]
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
			out.Records = append(out.Records, record{
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
		Zone     string `path:"zone"`
		Hostname string `path:"hostname"`
		UID      string `query:"uid"`

		Type string `json:"type"`

		Data         string `json:"data"`
		MXPreference uint16 `json:"mxPreference,omitempty"`
		SRVPort      uint16 `json:"srvPort,omitempty"`
		SRVPriority  uint16 `json:"srvPriority,omitempty"`
		SRVWeight    uint16 `json:"srvWeight,omitempty"`
	}
	u := usecase.NewIOI(new(recordsInput), new(struct{}), func(ctx context.Context, input, output interface{}) error {
		var (
			in = input.(*recordsInput)
		)
		zone, ok := r.zones[in.Zone]
		if !ok {
			return status.Wrap(errors.New("zone not found"), status.NotFound)
		}
		rec := zone.newRecord(in.Hostname, in.Type)
		rec.uid = in.UID
		rec.Data = in.Data
		rec.MXPreference = in.MXPreference
		rec.SRVPort = in.SRVPort
		rec.SRVPriority = in.SRVPriority
		rec.SRVWeight = in.SRVWeight
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
	type recordsInput struct {
		Zone     string `path:"zone"`
		Hostname string `path:"hostname"`
		UID      string `query:"uid"`
	}
	u := usecase.NewIOI(new(recordsInput), new(struct{}), func(ctx context.Context, input, output interface{}) error {
		var (
			in = input.(*recordsInput)
		)
		zone, ok := r.zones[in.Zone]
		if !ok {
			return status.Wrap(errors.New("zone not found"), status.NotFound)
		}
		key := r.i.KV().Key(types.KeyRole, types.KeyZones, in.Zone, in.Hostname)
		recs, ok := zone.records[key.String()]
		if !ok {
			return status.Wrap(errors.New("record not found"), status.NotFound)
		}
		rec, ok := recs[in.UID]
		if !ok {
			return status.Wrap(errors.New("record uid not found"), status.NotFound)
		}
		_, err := r.i.KV().Delete(ctx, rec.recordKey)
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
