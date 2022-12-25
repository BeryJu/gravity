package dhcp

import (
	"context"
	"errors"

	"beryju.io/gravity/pkg/roles/dhcp/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	"go.uber.org/zap"
)

type APILeasesGetInput struct {
	ScopeName string `query:"scope"`
}
type APILease struct {
	Identifier       string `json:"identifier" required:"true"`
	Address          string `json:"address" required:"true"`
	Hostname         string `json:"hostname" required:"true"`
	AddressLeaseTime string `json:"addressLeaseTime" required:"true"`
	ScopeKey         string `json:"scopeKey" required:"true"`
	DNSZone          string `json:"dnsZone"`
}
type APILeasesGetOutput struct {
	Leases []*APILease `json:"leases" required:"true"`
}

func (r *Role) APILeasesGet() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APILeasesGetInput, output *APILeasesGetOutput) error {
		r.leasesM.RLock()
		defer r.leasesM.RUnlock()
		for _, l := range r.leases {
			if l.ScopeKey == input.ScopeName {
				output.Leases = append(output.Leases, &APILease{
					Identifier:       l.Identifier,
					Address:          l.Address,
					Hostname:         l.Hostname,
					AddressLeaseTime: l.AddressLeaseTime,
					ScopeKey:         l.ScopeKey,
					DNSZone:          l.DNSZone,
				})
			}
		}
		return nil
	})
	u.SetName("dhcp.get_leases")
	u.SetTitle("DHCP Leases")
	u.SetTags("roles/dhcp")
	return u
}

type APILeasesPutInput struct {
	Identifier string `query:"identifier" required:"true" maxLength:"255"`
	Scope      string `query:"scope" required:"true" maxLength:"255"`

	Address          string `json:"address" required:"true" maxLength:"40"`
	Hostname         string `json:"hostname" required:"true" maxLength:"255"`
	AddressLeaseTime string `json:"addressLeaseTime" required:"true" maxLength:"40"`
	DNSZone          string `json:"dnsZone" maxLength:"255"`
}

func (r *Role) APILeasesPut() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APILeasesPutInput, output *struct{}) error {
		rawScope, err := r.i.KV().Get(
			ctx,
			r.i.KV().Key(
				types.KeyRole,
				types.KeyScopes,
				input.Scope,
			).String(),
		)
		if err != nil || len(rawScope.Kvs) < 1 {
			r.log.Warn("failed to get scope", zap.Error(err))
			return status.Wrap(errors.New("failed to get scope"), status.Internal)
		}
		scope, err := r.scopeFromKV(rawScope.Kvs[0])
		if err != nil {
			r.log.Warn("failed to construct scope", zap.Error(err))
			return status.Wrap(errors.New("failed to construct scope"), status.Internal)
		}

		l := r.NewLease(input.Identifier)
		l.Address = input.Address
		l.Hostname = input.Hostname
		l.AddressLeaseTime = input.AddressLeaseTime
		l.ScopeKey = input.Scope
		l.DNSZone = input.DNSZone
		l.scope = scope
		err = l.Put(ctx, -1)
		if err != nil {
			r.log.Warn("failed to put lease", zap.Error(err))
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("dhcp.put_leases")
	u.SetTitle("DHCP Leases")
	u.SetTags("roles/dhcp")
	u.SetExpectedErrors(status.Internal, status.InvalidArgument)
	return u
}

type APILeasesWOLInput struct {
	Identifier string `query:"identifier" required:"true"`
	Scope      string `query:"scope" required:"true"`
}

func (r *Role) APILeasesWOL() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APILeasesWOLInput, output *struct{}) error {
		r.leasesM.RLock()
		l, ok := r.leases[input.Identifier]
		r.leasesM.RUnlock()
		if !ok {
			return status.InvalidArgument
		}
		err := l.sendWOL()
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("dhcp.wol_leases")
	u.SetTitle("DHCP Leases")
	u.SetTags("roles/dhcp")
	u.SetExpectedErrors(status.Internal, status.InvalidArgument)
	return u
}

type APILeasesDeleteInput struct {
	Identifier string `query:"identifier"`
	Scope      string `query:"scope"`
}

func (r *Role) APILeasesDelete() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APILeasesDeleteInput, output *struct{}) error {
		l, ok := r.leases[input.Identifier]
		if !ok {
			return status.InvalidArgument
		}
		_, err := r.i.KV().Delete(ctx, l.etcdKey)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("dhcp.delete_leases")
	u.SetTitle("DHCP Leases")
	u.SetTags("roles/dhcp")
	u.SetExpectedErrors(status.Internal, status.InvalidArgument)
	return u
}
