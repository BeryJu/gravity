package dhcp

import (
	"context"
	"errors"
	"fmt"
	"net/netip"

	"beryju.io/gravity/pkg/convert/ms_dhcp"
	"beryju.io/gravity/pkg/roles/api/utils"
	"beryju.io/gravity/pkg/roles/dhcp/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type APIScopesGetInput struct {
	Name string `query:"name" description:"Optionally get DHCP Scope by name"`
}

type APIScopeStatistics struct {
	Usable uint64 `json:"usable" required:"true"`
	Used   uint64 `json:"used" required:"true"`
}
type APIScope struct {
	IPAM       map[string]string   `json:"ipam" required:"true"`
	DNS        *ScopeDNS           `json:"dns"`
	Name       string              `json:"scope" required:"true"`
	SubnetCIDR string              `json:"subnetCidr" required:"true"`
	Options    []*types.DHCPOption `json:"options" required:"true"`
	TTL        int64               `json:"ttl" required:"true"`
	Default    bool                `json:"default" required:"true"`
	Hook       string              `json:"hook" required:"true"`
	Statistics APIScopeStatistics  `json:"statistics" required:"true"`
}
type APIScopesGetOutput struct {
	Scopes     []*APIScope        `json:"scopes" required:"true"`
	Statistics APIScopeStatistics `json:"statistics" required:"true"`
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
		// Fetch all leases for statistics
		rawLeases, err := r.i.KV().Get(
			ctx,
			r.i.KV().Key(
				types.KeyRole,
				types.KeyLeases,
			).String(),
			clientv3.WithPrefix(),
		)
		if err != nil {
			r.log.Warn("failed to get leases", zap.Error(err))
			return status.Wrap(errors.New("failed to get leases"), status.Internal)
		}
		leases := []*Lease{}
		for _, rl := range rawLeases.Kvs {
			l, err := r.leaseFromKV(rl)
			if err != nil {
				r.log.Warn("failed to parse lease", zap.Error(err))
				continue
			}
			leases = append(leases, l)
		}

		// Generate summarized statistics
		sum := APIScopeStatistics{}

		for _, rawScope := range rawScopes.Kvs {
			sc, err := r.scopeFromKV(rawScope)
			if err != nil {
				r.log.Warn("failed to parse scope", zap.Error(err))
				continue
			}
			stat := APIScopeStatistics{
				Usable: sc.ipam.UsableSize().Uint64(),
				Used:   0,
			}
			for _, l := range leases {
				if l.ScopeKey != sc.Name {
					continue
				}
				stat.Used += 1
			}
			sum.Usable += stat.Usable
			sum.Used += stat.Used
			output.Scopes = append(output.Scopes, &APIScope{
				Name:       sc.Name,
				SubnetCIDR: sc.SubnetCIDR,
				Default:    sc.Default,
				Options:    sc.Options,
				TTL:        sc.TTL,
				IPAM:       sc.IPAM,
				DNS:        sc.DNS,
				Hook:       sc.Hook,
				Statistics: stat,
			})
		}
		output.Statistics = sum
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
	Hook       string              `json:"hook" required:"true"`
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
		s.Hook = input.Hook

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

type APIScopesImporterType string

func (APIScopesImporterType) Enum() []interface{} {
	return []interface{}{
		"ms_dhcp",
	}
}

type APIScopesImportInput struct {
	Type    APIScopesImporterType `json:"type"`
	Payload string                `json:"payload"`
	Scope   string                `query:"scope"`
}

type APIScopesImportOutput struct {
	Successful bool `json:"successful"`
}

type DHCPImporter interface {
	Run(ctx context.Context) error
}

func (r *Role) APIScopesImport() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIScopesImportInput, output *APIScopesImportOutput) error {
		var converter DHCPImporter
		var err error
		ac := utils.APIClientFromRequest(ctx)
		if ac == nil {
			return status.Wrap(errors.New("failed to get API Client from context"), status.Internal)
		}
		switch input.Type {
		case "ms_dhcp":
			converter, err = ms_dhcp.New(ac, input.Payload, ms_dhcp.WithExistingScope(input.Scope))
		default:
			err = status.WithDescription(status.InvalidArgument, fmt.Sprintf("invalid converter type specified: %s", input.Type))
		}
		if err != nil {
			return status.Wrap(err, status.InvalidArgument)
		}
		err = converter.Run(ctx)
		if err != nil {
			return status.Wrap(err, status.InvalidArgument)
		}
		output.Successful = true
		return nil
	})
	u.SetName("dhcp.import_scopes")
	u.SetTitle("DHCP Scopes")
	u.SetTags("roles/dhcp")
	u.SetExpectedErrors(status.Internal, status.InvalidArgument)
	return u
}
