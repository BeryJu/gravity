package dhcp

import (
	"context"

	"github.com/swaggest/usecase"
)

func (r *DHCPRole) apiHandlerScopes() usecase.Interactor {
	type scopesOutput struct {
		Scopes []*Scope `json:"scopes"`
	}
	u := usecase.NewIOI(new(struct{}), new(scopesOutput), func(ctx context.Context, input, output interface{}) error {
		var (
			out = output.(*scopesOutput)
		)
		for _, scope := range r.scopes {
			out.Scopes = append(out.Scopes, scope)
		}
		return nil
	})
	u.SetName("dhcp.get_scopes")
	u.SetTitle("DHCP Scopes")
	u.SetTags("roles/dhcp")
	u.SetDescription("List all DHCP scopes.")
	return u
}
