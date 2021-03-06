package api

import (
	"context"
	"encoding/json"

	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type RoleConfig struct {
	Port int32 `json:"port"`
}

func (r *Role) decodeRoleConfig(raw []byte) *RoleConfig {
	def := RoleConfig{
		Port: 8008,
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
	type roleAPIConfigOutput struct {
		Config *RoleConfig `json:"config"`
	}
	u := usecase.NewIOI(new(struct{}), new(roleAPIConfigOutput), func(ctx context.Context, input, output interface{}) error {
		var (
			out = output.(*roleAPIConfigOutput)
		)
		out.Config = r.cfg
		return nil
	})
	u.SetName("api.get_role_config")
	u.SetTitle("API role config")
	u.SetTags("roles/api")
	return u
}

func (r *Role) apiHandlerRoleConfigPut() usecase.IOInteractor {
	type roleAPIConfigInput struct {
		Config *RoleConfig `json:"config"`
	}
	u := usecase.NewIOI(new(roleAPIConfigInput), new(struct{}), func(ctx context.Context, input, output interface{}) error {
		var (
			in = input.(*roleAPIConfigInput)
		)
		jc, err := json.Marshal(in.Config)
		if err != nil {
			return status.Wrap(err, status.InvalidArgument)
		}
		_, err = r.i.KV().Put(ctx, r.i.KV().Key(types.KeyRole, "api").String(), string(jc))
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("api.put_role_config")
	u.SetTitle("API role config")
	u.SetTags("roles/api")
	u.SetExpectedErrors(status.InvalidArgument, status.Internal)
	return u
}
