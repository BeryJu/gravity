package dns

import (
	"context"
	"errors"
	"fmt"
	"slices"
	"strings"

	"beryju.io/gravity/pkg/convert/bind"
	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/roles/dns/utils"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type APIZonesGetInput struct {
	Name string `query:"name"  description:"Optionally get DNS Zone by name"`
}
type APIZone struct {
	Name           string                   `json:"name" required:"true"`
	HandlerConfigs []map[string]interface{} `json:"handlerConfigs" required:"true"`
	DefaultTTL     uint32                   `json:"defaultTTL" required:"true"`
	Authoritative  bool                     `json:"authoritative" required:"true"`
	Hook           string                   `json:"hook" required:"true"`
	RecordCount    int                      `json:"recordCount" required:"true"`
}
type APIZonesGetOutput struct {
	Zones []APIZone `json:"zones" required:"true"`
}

func (r *Role) APIZonesGet() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIZonesGetInput, output *APIZonesGetOutput) error {
		key := r.i.KV().Key(
			types.KeyRole,
			types.KeyZones,
		)
		if input.Name == "" {
			key = key.Prefix(true)
		} else {
			key = key.Add(input.Name)
		}
		rawZones, err := r.i.KV().Get(
			ctx,
			key.String(),
			clientv3.WithPrefix(),
		)
		if err != nil {
			r.log.Warn("failed to get zones", zap.Error(err))
			return status.Wrap(errors.New("failed to get zones"), status.Internal)
		}
		rawRecords, err := r.i.KV().Get(
			ctx,
			r.i.KV().Key(types.KeyRole, types.KeyZones).Prefix(true).String(),
			clientv3.WithPrefix(),
			clientv3.WithKeysOnly(),
		)
		if err != nil {
			return status.Wrap(errors.New("failed to get records"), status.Internal)
		}
		for _, rawZone := range rawZones.Kvs {
			_zone, err := r.zoneFromKV(rawZone)
			if err != nil {
				r.log.Warn("failed to parse zone", zap.Error(err))
				continue
			}
			if strings.Contains(_zone.Name, "/") {
				continue
			}
			recordCount := slices.Collect(func(yield func(int) bool) {
				for _, v := range rawRecords.Kvs {
					if strings.HasPrefix(string(v.Key), _zone.etcdKey+"/") {
						ok := yield(1)
						if !ok {
							return
						}
					}
				}
			})
			output.Zones = append(output.Zones, APIZone{
				Name:           _zone.Name,
				Authoritative:  _zone.Authoritative,
				DefaultTTL:     _zone.DefaultTTL,
				HandlerConfigs: _zone.HandlerConfigs,
				Hook:           _zone.Hook,
				RecordCount:    len(recordCount),
			})
		}
		return nil
	})
	u.SetName("dns.get_zones")
	u.SetTitle("DNS Zones")
	u.SetTags("roles/dns")
	u.SetExpectedErrors(status.Internal)
	return u
}

type APIZonesPutInput struct {
	Name           string                   `query:"zone" required:"true" maxLength:"255"`
	HandlerConfigs []map[string]interface{} `json:"handlerConfigs" required:"true"`
	DefaultTTL     uint32                   `json:"defaultTTL" required:"true"`
	Authoritative  bool                     `json:"authoritative" required:"true"`
	Hook           string                   `json:"hook" required:"true"`
}

func (r *Role) APIZonesPut() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIZonesPutInput, output *struct{}) error {
		z := r.newZone(input.Name)
		z.Name = utils.EnsureTrailingPeriod(input.Name)
		z.Authoritative = input.Authoritative
		z.HandlerConfigs = input.HandlerConfigs
		z.DefaultTTL = input.DefaultTTL
		z.Hook = input.Hook
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

type APIZonesDeleteInput struct {
	Zone string `query:"zone"`
}

func (r *Role) APIZonesDelete() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIZonesDeleteInput, output *struct{}) error {
		_, err := r.i.KV().Delete(
			ctx,
			r.i.KV().Key(
				types.KeyRole,
				types.KeyZones,
				input.Zone,
			).String(),
			clientv3.WithPrefix(),
		)
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

type APIZonesImportInput struct {
	Type    string `json:"type" enum:"bind"`
	Payload string `json:"payload"`
	Zone    string `query:"zone"`
}

type APIZonesImportOutput struct{}

type DNSImporter interface {
	Run(ctx context.Context) error
}

func (r *Role) APIZonesImport() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIZonesImportInput, output *APIZonesImportOutput) error {
		var converter DNSImporter
		var err error
		switch input.Type {
		case "bind":
			converter, err = bind.New(nil, input.Payload)
		default:
			err = status.WithDescription(status.InvalidArgument, fmt.Sprintf("invalid converter type specified: %s", input.Type))
		}
		if err != nil {
			return status.Wrap(err, status.InvalidArgument)
		}
		err = converter.Run(ctx)
		if err != nil {
			return status.Wrap(err, status.InvalidArgument)
		}
		return nil
	})
	u.SetName("dns.import_zones")
	u.SetTitle("DNS Zones")
	u.SetTags("roles/dns")
	u.SetExpectedErrors(status.Internal, status.InvalidArgument)
	return u
}
