package dhcp

import (
	"context"
	"net/netip"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

func (r *DHCPRole) apiHandlerScopes() usecase.Interactor {
	type scope struct {
		Name       string    `json:"scope"`
		SubnetCIDR string    `json:"subnetCidr"`
		Default    bool      `json:"default"`
		Options    []*Option `json:"options"`
		TTL        int64     `json:"ttl"`
		Range      struct {
			Start string `json:"start"`
			End   string `json:"end"`
		} `json:"range"`
		DNS struct {
			Zone              string   `json:"zone"`
			Search            []string `json:"search"`
			AddZoneInHostname bool     `json:"addZoneInHostname"`
		} `json:"dns"`
	}
	type scopesOutput struct {
		Scopes []*scope `json:"scopes"`
	}
	u := usecase.NewIOI(new(struct{}), new(scopesOutput), func(ctx context.Context, input, output interface{}) error {
		var (
			out = output.(*scopesOutput)
		)
		for _, sc := range r.scopes {
			out.Scopes = append(out.Scopes, &scope{
				Name:       sc.Name,
				SubnetCIDR: sc.SubnetCIDR,
				Default:    sc.Default,
				Options:    sc.Options,
				TTL:        sc.TTL,
				Range:      sc.Range,
				DNS:        sc.DNS,
			})
		}
		return nil
	})
	u.SetName("dhcp.get_scopes")
	u.SetTitle("DHCP Scopes")
	u.SetTags("roles/dhcp")
	return u
}

func (r *DHCPRole) apiHandlerScopesPut() usecase.Interactor {
	type scopesInput struct {
		Name string `path:"scope"`

		SubnetCIDR string    `json:"subnetCidr"`
		Default    bool      `json:"default"`
		Options    []*Option `json:"options"`
		TTL        int64     `json:"ttl"`
		Range      struct {
			Start string `json:"start"`
			End   string `json:"end"`
		} `json:"range"`
		DNS struct {
			Zone              string   `json:"zone"`
			Search            []string `json:"search"`
			AddZoneInHostname bool     `json:"addZoneInHostname"`
		} `json:"dns"`
	}
	u := usecase.NewIOI(new(scopesInput), new(struct{}), func(ctx context.Context, input, output interface{}) error {
		var (
			in = input.(*scopesInput)
		)
		s := r.newScope(in.Name)
		s.SubnetCIDR = in.SubnetCIDR
		s.Default = in.Default
		s.Options = in.Options
		s.TTL = in.TTL
		s.Range = in.Range
		s.DNS = in.DNS

		cidr, err := netip.ParsePrefix(s.SubnetCIDR)
		if err != nil {
			return status.Wrap(err, status.InvalidArgument)
		}
		s.cidr = cidr

		err = s.put(-1)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("dhcp.put_scope")
	u.SetTitle("DHCP Scopes")
	u.SetTags("roles/dhcp")
	u.SetExpectedErrors(status.Internal, status.InvalidArgument)
	return u
}

func (r *DHCPRole) apiHandlerScopesDelete() usecase.Interactor {
	type scopesInput struct {
		Scope string `path:"scope"`
	}
	u := usecase.NewIOI(new(scopesInput), new(struct{}), func(ctx context.Context, input, output interface{}) error {
		var (
			in = input.(*scopesInput)
		)
		s, ok := r.leases[in.Scope]
		if !ok {
			return status.InvalidArgument
		}
		_, err := r.i.KV().Delete(ctx, s.etcdKey)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("dhcp.delete_scopes")
	u.SetTitle("DHCP Scopes")
	u.SetTags("roles/dhcp")
	u.SetExpectedErrors(status.Internal, status.InvalidArgument)
	return u
}
