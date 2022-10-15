package dns

import (
	"context"
	"errors"

	"beryju.io/gravity/pkg/roles/dns/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type APIRecordsGetInput struct {
	Zone string `query:"zone"`
}
type APIRecord struct {
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
type APIRecordsGetOutput struct {
	Records []APIRecord `json:"records" required:"true"`
}

func (r *Role) APIRecordsGet() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIRecordsGetInput, output *APIRecordsGetOutput) error {
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
				r.log.Warn("failed to parse record", zap.Error(err))
				continue
			}
			output.Records = append(output.Records, APIRecord{
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

type APIRecordsPutInput struct {
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

func (r *Role) APIRecordsPut() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIRecordsPutInput, output *struct{}) error {
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

type APIRecordsDeleteInput struct {
	Zone     string `query:"zone" required:"true"`
	Hostname string `query:"hostname" required:"true"`
	UID      string `query:"uid" required:"true"`
	Type     string `query:"type" required:"true"`
}

func (r *Role) APIRecordsDelete() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIRecordsDeleteInput, output *struct{}) error {
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
