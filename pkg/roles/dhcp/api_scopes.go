package dhcp

import (
	"context"
	"net/netip"

	"beryju.io/gravity/pkg/roles/dhcp/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type APIScope struct {
	Name       string            `json:"scope" required:"true"`
	SubnetCIDR string            `json:"subnetCidr" required:"true"`
	Default    bool              `json:"default" required:"true"`
	Options    []*types.Option   `json:"options" required:"true"`
	TTL        int64             `json:"ttl" required:"true"`
	IPAM       map[string]string `json:"ipam" required:"true"`
	DNS        *ScopeDNS         `json:"dns"`
}
type APIScopesGetOutput struct {
	Scopes []*APIScope `json:"scopes" required:"true"`
}

func (r *Role) APIScopesGet() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *APIScopesGetOutput) error {
		for _, sc := range r.scopes {
			output.Scopes = append(output.Scopes, &APIScope{
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

type APIScopesPutInput struct {
	Name string `query:"scope" required:"true" maxLength:"255"`

	SubnetCIDR string            `json:"subnetCidr" required:"true" maxLength:"40"`
	Default    bool              `json:"default" required:"true"`
	Options    []*types.Option   `json:"options" required:"true"`
	TTL        int64             `json:"ttl" required:"true"`
	IPAM       map[string]string `json:"ipam"`
	DNS        *ScopeDNS         `json:"dns"`
}

func (r *Role) APIScopesPut() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIScopesPutInput, output *struct{}) error {
		if input.Name == "" {
			return status.InvalidArgument
		}
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

type APIScopesDeleteInput struct {
	Scope string `query:"scope" required:"true"`
}

func (r *Role) APIScopesDelete() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIScopesDeleteInput, output *struct{}) error {
		_, err := r.i.KV().Delete(
			ctx,
			r.i.KV().Key(
				types.KeyRole,
				types.KeyScopes,
				input.Scope,
			).String(),
			clientv3.WithPrefix(),
		)
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
