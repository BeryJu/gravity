package dhcp

import (
	"context"
	"encoding/json"

	instanceTypes "beryju.io/gravity/pkg/instance/types"
	"beryju.io/gravity/pkg/roles/dhcp/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type RoleConfig struct {
	Port                  int `json:"port"`
	LeaseNegotiateTimeout int `json:"leaseNegotiateTimeout"`
}

func (r *Role) decodeRoleConfig(raw []byte) *RoleConfig {
	def := RoleConfig{
		Port:                  67,
		LeaseNegotiateTimeout: 30,
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
	type roleDHCPConfigOutput struct {
		Config *RoleConfig `json:"config"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *roleDHCPConfigOutput) error {
		output.Config = r.cfg
		return nil
	})
	u.SetName("dhcp.get_role_config")
	u.SetTitle("DHCP role config")
	u.SetTags("roles/dhcp")
	return u
}

func (r *Role) apiHandlerRoleConfigPut() usecase.Interactor {
	type roleDHCPConfigInput struct {
		Config *RoleConfig `json:"config"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input roleDHCPConfigInput, output *struct{}) error {
		jc, err := json.Marshal(input.Config)
		if err != nil {
			return status.Wrap(err, status.InvalidArgument)
		}
		_, err = r.i.KV().Put(ctx, r.i.KV().Key(instanceTypes.KeyRole, types.KeyRole).String(), string(jc))
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("dhcp.put_role_config")
	u.SetTitle("DHCP role config")
	u.SetTags("roles/dhcp")
	u.SetExpectedErrors(status.InvalidArgument, status.Internal)
	return u
}
