package monitoring

import (
	"context"
	"encoding/json"

	instanceTypes "beryju.io/gravity/pkg/instance/types"
	"beryju.io/gravity/pkg/roles/api/types"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

const (
	KeyRole = "monitoring"
)

type RoleConfig struct {
	Port int32 `json:"port"`
}

func (r *Role) decodeRoleConfig(raw []byte) *RoleConfig {
	def := RoleConfig{
		Port: 8009,
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
	type roleMonitoringConfigOutput struct {
		Config RoleConfig `json:"config" required:"true"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *roleMonitoringConfigOutput) error {
		output.Config = *r.cfg
		return nil
	})
	u.SetName("monitoring.get_role_config")
	u.SetTitle("Monitoring role config")
	u.SetTags("roles/monitoring")
	return u
}

func (r *Role) apiHandlerRoleConfigPut() usecase.Interactor {
	type roleMonitoringConfigInput struct {
		Config RoleConfig `json:"config" required:"true"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input roleMonitoringConfigInput, output *struct{}) error {
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
	u.SetName("monitoring.put_role_config")
	u.SetTitle("Monitoring role config")
	u.SetTags("roles/monitoring")
	u.SetExpectedErrors(status.InvalidArgument, status.Internal)
	return u
}
