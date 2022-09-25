package dhcp

import (
	"context"
	"net/netip"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (r *Role) apiHandlerScopesGet() usecase.Interactor {
	type scope struct {
		Name       string            `json:"scope" required:"true"`
		SubnetCIDR string            `json:"subnetCidr" required:"true"`
		Default    bool              `json:"default" required:"true"`
		Options    []*Option         `json:"options" required:"true"`
		TTL        int64             `json:"ttl" required:"true"`
		IPAM       map[string]string `json:"ipam" required:"true"`
		DNS        struct {
			Zone              string   `json:"zone"`
			Search            []string `json:"search"`
			AddZoneInHostname bool     `json:"addZoneInHostname"`
		} `json:"dns"`
	}
	type scopesOutput struct {
		Scopes []*scope `json:"scopes" required:"true"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *scopesOutput) error {
		for _, sc := range r.scopes {
			output.Scopes = append(output.Scopes, &scope{
				Name:       sc.Name,
				SubnetCIDR: sc.SubnetCIDR,
				Default:    sc.Default,
				Options:    sc.Options,
				TTL:        sc.TTL,
				IPAM:       sc.IPAM,
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

func (r *Role) apiHandlerScopesPut() usecase.Interactor {
	type scopesInput struct {
		Name string `query:"scope" required:"true" maxLength:"255"`

		SubnetCIDR string            `json:"subnetCidr" required:"true" maxLength:"40"`
		Default    bool              `json:"default" required:"true"`
		Options    []*Option         `json:"options" required:"true"`
		TTL        int64             `json:"ttl" required:"true"`
		IPAM       map[string]string `json:"ipam"`
		DNS        struct {
			Zone              string   `json:"zone"`
			Search            []string `json:"search"`
			AddZoneInHostname bool     `json:"addZoneInHostname"`
		} `json:"dns"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input scopesInput, output *struct{}) error {
		s := r.newScope(input.Name)
		s.SubnetCIDR = input.SubnetCIDR
		s.Default = input.Default
		s.Options = input.Options
		s.TTL = input.TTL
		s.IPAM = input.IPAM
		s.DNS = input.DNS

		cidr, err := netip.ParsePrefix(s.SubnetCIDR)
		if err != nil {
			return status.Wrap(err, status.InvalidArgument)
		}
		s.cidr = cidr

		err = s.put(ctx, -1)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("dhcp.put_scopes")
	u.SetTitle("DHCP Scopes")
	u.SetTags("roles/dhcp")
	u.SetExpectedErrors(status.Internal, status.InvalidArgument)
	return u
}

func (r *Role) apiHandlerScopesDelete() usecase.Interactor {
	type scopesInput struct {
		Scope string `query:"scope" required:"true"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input scopesInput, output *struct{}) error {
		s, ok := r.scopes[input.Scope]
		if !ok {
			return status.InvalidArgument
		}
		_, err := r.i.KV().Delete(ctx, s.etcdKey, clientv3.WithPrefix())
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
