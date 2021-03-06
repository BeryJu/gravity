package discovery

import (
	"context"
	"encoding/json"

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

func (r *Role) apiHandlerRoleConfigGet() usecase.IOInteractor {
	type roleDiscoveryConfigOutput struct {
		Config *RoleConfig `json:"config"`
	}
	u := usecase.NewIOI(new(struct{}), new(roleDiscoveryConfigOutput), func(ctx context.Context, input, output interface{}) error {
		var (
			out = output.(*roleDiscoveryConfigOutput)
		)
		out.Config = r.cfg
		return nil
	})
	u.SetName("discovery.get_role_config")
	u.SetTitle("Discovery role config")
	u.SetTags("roles/discovery")
	return u
}

func (r *Role) apiHandlerRoleConfigPut() usecase.IOInteractor {
	type roleDiscoveryConfigInput struct {
		Config *RoleConfig `json:"config"`
	}
	u := usecase.NewIOI(new(roleDiscoveryConfigInput), new(struct{}), func(ctx context.Context, input, output interface{}) error {
		var (
			in = input.(*roleDiscoveryConfigInput)
		)
		jc, err := json.Marshal(in.Config)
		if err != nil {
			return status.Wrap(err, status.InvalidArgument)
		}
		_, err = r.i.KV().Put(ctx, r.i.KV().Key(types.KeyRole, "discovery").String(), string(jc))
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
