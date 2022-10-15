package dns

import (
	"context"
	"errors"
	"strings"

	"beryju.io/gravity/pkg/roles/dns/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type APIZone struct {
	Name           string              `json:"name" required:"true"`
	Authoritative  bool                `json:"authoritative" required:"true"`
	HandlerConfigs []map[string]string `json:"handlerConfigs" required:"true"`
	DefaultTTL     uint32              `json:"defaultTTL" required:"true"`
}
type APIZonesGetOutput struct {
	Zones []APIZone `json:"zones" required:"true"`
}

func (r *Role) APIZonesGet() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *APIZonesGetOutput) error {
		rawZones, err := r.i.KV().Get(
			ctx,
			r.i.KV().Key(
				types.KeyRole,
				types.KeyZones,
			).Prefix(true).String(),
			clientv3.WithPrefix(),
		)
		if err != nil {
			r.log.Warn("failed to get zones", zap.Error(err))
			return status.Wrap(errors.New("failed to get zones"), status.Internal)
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
			output.Zones = append(output.Zones, APIZone{
				Name:           _zone.Name,
				Authoritative:  _zone.Authoritative,
				DefaultTTL:     _zone.DefaultTTL,
				HandlerConfigs: _zone.HandlerConfigs,
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
	Name           string              `query:"zone" required:"true" maxLength:"255"`
	Authoritative  bool                `json:"authoritative" required:"true"`
	HandlerConfigs []map[string]string `json:"handlerConfigs" required:"true"`
	DefaultTTL     uint32              `json:"defaultTTL" required:"true"`
}

func (r *Role) APIZonesPut() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIZonesPutInput, output *struct{}) error {
		z := r.newZone(input.Name)
		z.Name = input.Name
		if !strings.HasSuffix(z.Name, ".") {
			z.Name += "."
		}
		z.Authoritative = input.Authoritative
		z.HandlerConfigs = input.HandlerConfigs
		z.DefaultTTL = input.DefaultTTL
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
