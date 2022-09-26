package discovery

import (
	"context"
	"encoding/json"

	instanceTypes "beryju.io/gravity/pkg/instance/types"
	"beryju.io/gravity/pkg/roles/discovery/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type RoleConfig struct {
	Enabled bool `json:"enabled"`
}

func (r *Role) decodeRoleConfig(raw []byte) *RoleConfig {
	def := RoleConfig{
		Enabled: true,
	}
	if len(raw) < 1 {
		return &def
	}
	err := json.Unmarshal(raw, &def)
	if err != nil {
		r.log.WithError(err).Warning("failed to decode role config")
	}
	return &def
}

func (r *Role) apiHandlerRoleConfigGet() usecase.Interactor {
	type roleDiscoveryConfigOutput struct {
		Config RoleConfig `json:"config" required:"true"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *roleDiscoveryConfigOutput) error {
		output.Config = *r.cfg
		return nil
	})
	u.SetName("discovery.get_role_config")
	u.SetTitle("Discovery role config")
	u.SetTags("roles/discovery")
	return u
}

func (r *Role) apiHandlerRoleConfigPut() usecase.Interactor {
	type roleDiscoveryConfigInput struct {
		Config RoleConfig `json:"config" required:"true"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input roleDiscoveryConfigInput, output *struct{}) error {
		jc, err := json.Marshal(input.Config)
		if err != nil {
			return status.Wrap(err, status.InvalidArgument)
		}
		_, err = r.i.KV().Put(
			ctx,
			r.i.KV().Key(
				instanceTypes.KeyInstance,
				instanceTypes.KeyRole,
				types.KeyRole,
			).String(),
			string(jc),
		)
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("discovery.put_role_config")
	u.SetTitle("Discovery role config")
	u.SetTags("roles/discovery")
	u.SetExpectedErrors(status.InvalidArgument, status.Internal)
	return u
}
