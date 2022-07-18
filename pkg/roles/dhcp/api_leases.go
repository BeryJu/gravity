package dhcp

import (
	"context"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func (r *Role) apiHandlerLeasesGet() usecase.Interactor {
	type leasesInput struct {
		ScopeName string `path:"scope"`
	}
	type lease struct {
		Identifier       string `json:"identifier"`
		Address          string `json:"address"`
		Hostname         string `json:"hostname"`
		AddressLeaseTime string `json:"addressLeaseTime,omitempty"`
		ScopeKey         string `json:"scopeKey"`
		DNSZone          string `json:"dnsZone"`
	}
	type leasesOutput struct {
		Leases []*lease `json:"leases"`
	}
	u := usecase.NewIOI(new(leasesInput), new(leasesOutput), func(ctx context.Context, input, output interface{}) error {
		var (
			in  = input.(*leasesInput)
			out = output.(*leasesOutput)
		)
		for _, l := range r.leases {
			if l.ScopeKey == in.ScopeName {
				out.Leases = append(out.Leases, &lease{
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
		Identifier string `path:"identifier"`
		Scope      string `path:"scope"`

		Address          string `json:"address"`
		Hostname         string `json:"hostname"`
		AddressLeaseTime string `json:"addressLeaseTime"`
		DNSZone          string `json:"dnsZone"`
	}
	u := usecase.NewIOI(new(leasesInput), new(struct{}), func(ctx context.Context, input, output interface{}) error {
		var (
			in = input.(*leasesInput)
		)
		l := r.newLease(in.Identifier)
		l.Address = in.Address
		l.Hostname = in.Hostname
		l.AddressLeaseTime = in.AddressLeaseTime
		l.ScopeKey = in.Scope
		l.DNSZone = in.DNSZone
		scope, ok := r.scopes[in.Scope]
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
			r.log.WithError(err).Warning("failed to WOL")
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
