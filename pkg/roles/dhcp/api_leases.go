package dhcp

import (
	"context"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func (r *Role) apiHandlerLeasesGet() usecase.Interactor {
	type leasesInput struct {
		ScopeName string `query:"scope"`
	}
	type lease struct {
		Identifier       string `json:"identifier" required:"true"`
		Address          string `json:"address" required:"true"`
		Hostname         string `json:"hostname" required:"true"`
		AddressLeaseTime string `json:"addressLeaseTime" required:"true"`
		ScopeKey         string `json:"scopeKey" required:"true"`
		DNSZone          string `json:"dnsZone"`
	}
	type leasesOutput struct {
		Leases []*lease `json:"leases" required:"true"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input leasesInput, output *leasesOutput) error {
		for _, l := range r.leases {
			if l.ScopeKey == input.ScopeName {
				output.Leases = append(output.Leases, &lease{
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

func (r *Role) apiHandlerLeasesPut() usecase.Interactor {
	type leasesInput struct {
		Identifier string `query:"identifier" required:"true"`
		Scope      string `query:"scope" required:"true"`

		Address          string `json:"address" required:"true"`
		Hostname         string `json:"hostname" required:"true"`
		AddressLeaseTime string `json:"addressLeaseTime" required:"true"`
		DNSZone          string `json:"dnsZone"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input leasesInput, output *struct{}) error {
		l := r.newLease(input.Identifier)
		l.Address = input.Address
		l.Hostname = input.Hostname
		l.AddressLeaseTime = input.AddressLeaseTime
		l.ScopeKey = input.Scope
		l.DNSZone = input.DNSZone
		scope, ok := r.scopes[input.Scope]
		if !ok {
			return status.InvalidArgument
		}
		l.scope = scope
		err := l.put(ctx, -1)
		if err != nil {
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

func (r *Role) apiHandlerLeasesWOL() usecase.Interactor {
	type leasesInput struct {
		Identifier string `query:"identifier" required:"true"`
		Scope      string `query:"scope" required:"true"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input leasesInput, output *struct{}) error {
		l, ok := r.leases[input.Identifier]
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

func (r *Role) apiHandlerLeasesDelete() usecase.Interactor {
	type leasesInput struct {
		Identifier string `query:"identifier"`
		Scope      string `query:"scope"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input leasesInput, output *struct{}) error {
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
