package dhcp

import (
	"context"
	"errors"
	"net/netip"
	"strings"

	"beryju.io/gravity/pkg/roles/dhcp/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type APIScopesGetInput struct {
	Name string `query:"name" description:"Optionally get DHCP Scope by name"`
}
type APIScope struct {
	IPAM       map[string]string   `json:"ipam" required:"true"`
	DNS        *ScopeDNS           `json:"dns"`
	Name       string              `json:"scope" required:"true"`
	SubnetCIDR string              `json:"subnetCidr" required:"true"`
	Options    []*types.DHCPOption `json:"options" required:"true"`
	TTL        int64               `json:"ttl" required:"true"`
	Default    bool                `json:"default" required:"true"`
}
type APIScopesGetOutput struct {
	Scopes []*APIScope `json:"scopes" required:"true"`
}

func (r *Role) APIScopesGet() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIScopesGetInput, output *APIScopesGetOutput) error {
		key := r.i.KV().Key(
			types.KeyRole,
			types.KeyScopes,
		)
		if input.Name == "" {
			key = key.Prefix(true)
		} else {
			key = key.Add(input.Name)
		}
		rawScopes, err := r.i.KV().Get(
			ctx,
			key.String(),
			clientv3.WithPrefix(),
		)
		if err != nil {
			r.log.Warn("failed to get scopes", zap.Error(err))
			return status.Wrap(errors.New("failed to get scopes"), status.Internal)
		}
		for _, rawScope := range rawScopes.Kvs {
			sc, err := r.scopeFromKV(rawScope)
			if err != nil {
				r.log.Warn("failed to parse scope", zap.Error(err))
				continue
			}
			if strings.Contains(sc.Name, "/") {
				continue
			}
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
	IPAM map[string]string `json:"ipam"`
	DNS  *ScopeDNS         `json:"dns"`
	Name string            `query:"scope" required:"true" maxLength:"255"`

	SubnetCIDR string              `json:"subnetCidr" required:"true" maxLength:"40"`
	Options    []*types.DHCPOption `json:"options" required:"true"`
	TTL        int64               `json:"ttl" required:"true"`
	Default    bool                `json:"default" required:"true"`
}

func (r *Role) APIScopesPut() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIScopesPutInput, output *struct{}) error {
		if input.Name == "" {
			return status.InvalidArgument
		}
		s := r.NewScope(input.Name)
		s.SubnetCIDR = input.SubnetCIDR
		s.Default = input.Default
		s.Options = input.Options
		// validate options
		for _, opt := range s.Options {
			if opt.Tag != nil && opt.TagName == "" {
				continue
			}
			_, ok := types.TagMap[types.OptionTagName(opt.TagName)]
			if !ok {
				return status.InvalidArgument
			}
		}

		s.TTL = input.TTL
		s.IPAM = input.IPAM
		s.DNS = input.DNS

		cidr, err := netip.ParsePrefix(s.SubnetCIDR)
		if err != nil {
			return status.Wrap(err, status.InvalidArgument)
		}
		s.cidr = cidr

		_, err = s.ipamType(nil)
		if err != nil {
			return status.Wrap(err, status.InvalidArgument)
		}

		err = s.Put(ctx, -1)
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
