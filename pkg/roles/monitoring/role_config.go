package monitoring

import (
	"context"
	"encoding/json"

	instanceTypes "beryju.io/gravity/pkg/instance/types"

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

type APIRoleMonitoringConfigOutput struct {
	Config RoleConfig `json:"config" required:"true"`
}

func (r *Role) APIRoleConfigGet() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *APIRoleMonitoringConfigOutput) error {
		output.Config = *r.cfg
		return nil
	})
	u.SetName("monitoring.get_role_config")
	u.SetTitle("Monitoring role config")
	u.SetTags("roles/monitoring")
	return u
}

type APIRoleMonitoringConfigInput struct {
	Config RoleConfig `json:"config" required:"true"`
}

func (r *Role) APIRoleConfigPut() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIRoleMonitoringConfigInput, output *struct{}) error {
		jc, err := json.Marshal(input.Config)
		if err != nil {
			return status.Wrap(err, status.InvalidArgument)
		}
		_, err = r.i.KV().Put(
			ctx,
			r.i.KV().Key(
				instanceTypes.KeyInstance,
				instanceTypes.KeyRole,
				KeyRole,
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
