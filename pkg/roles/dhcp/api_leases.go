package dhcp

import (
	"context"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func (r *Role) apiHandlerLeases() usecase.Interactor {
	type leasesInput struct {
		ScopeName string `path:"scope"`
	}
	type leasesOutput struct {
		Leases []*Lease `json:"leases"`
	}
	u := usecase.NewIOI(new(leasesInput), new(leasesOutput), func(ctx context.Context, input, output interface{}) error {
		var (
			in  = input.(*leasesInput)
			out = output.(*leasesOutput)
		)
		for _, lease := range r.leases {
			if lease.ScopeKey == in.ScopeName {
				out.Leases = append(out.Leases, lease)
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
		Identifier string `path:"identifier"`
		Scope      string `path:"scope"`

		Address          string `json:"address"`
		Hostname         string `json:"hostname"`
		AddressLeaseTime string `json:"addressLeaseTime"`
	}
	u := usecase.NewIOI(new(leasesInput), new(struct{}), func(ctx context.Context, input, output interface{}) error {
		var (
			in = input.(*leasesInput)
		)
		l := r.newLease()
		l.Address = in.Address
		l.Hostname = in.Hostname
		l.AddressLeaseTime = in.AddressLeaseTime
		l.ScopeKey = in.Scope
		scope, ok := r.scopes[in.Scope]
		if !ok {
			return status.InvalidArgument
		}
		l.scope = scope
		err := l.put(-1)
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
		Identifier string `path:"identifier"`
		Scope      string `path:"scope"`
	}
	u := usecase.NewIOI(new(leasesInput), new(struct{}), func(ctx context.Context, input, output interface{}) error {
		var (
			in = input.(*leasesInput)
		)
		l, ok := r.leases[in.Identifier]
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
		Identifier string `path:"identifier"`
		Scope      string `path:"scope"`
	}
	u := usecase.NewIOI(new(leasesInput), new(struct{}), func(ctx context.Context, input, output interface{}) error {
		var (
			in = input.(*leasesInput)
		)
		l, ok := r.leases[in.Identifier]
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
