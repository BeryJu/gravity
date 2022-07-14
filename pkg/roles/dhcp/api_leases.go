package dhcp

import (
	"context"

	"github.com/swaggest/usecase"
)

func (r *DHCPRole) apiHandlerLeases() usecase.Interactor {
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
	u.SetDescription("List all DHCP leases.")
	return u
}
