package dns

import (
	"context"
	"encoding/json"

	instanceTypes "beryju.io/gravity/pkg/instance/types"

	"beryju.io/gravity/pkg/roles/dns/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type RoleConfig struct {
	Port int32 `json:"port"`
}

func (r *Role) decodeRoleConfig(raw []byte) *RoleConfig {
	def := RoleConfig{
		Port: 53,
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

type APIRoleDNSConfigOutput struct {
	Config RoleConfig `json:"config" required:"true"`
}

func (r *Role) APIRoleConfigGet() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *APIRoleDNSConfigOutput) error {
		output.Config = *r.cfg
		return nil
	})
	u.SetName("dns.get_role_config")
	u.SetTitle("DNS role config")
	u.SetTags("roles/dns")
	return u
}

type APIRoleDNSConfigInput struct {
	Config RoleConfig `json:"config" required:"true"`
}

func (r *Role) APIRoleConfigPut() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIRoleDNSConfigInput, output *struct{}) error {
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
	u.SetName("dns.put_role_config")
	u.SetTitle("DNS role config")
	u.SetTags("roles/dns")
	u.SetExpectedErrors(status.InvalidArgument, status.Internal)
	return u
}
