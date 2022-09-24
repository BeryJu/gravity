package dns

import (
	"context"
	"encoding/json"

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

func (r *Role) apiHandlerRoleConfigGet() usecase.Interactor {
	type roleDNSConfigOutput struct {
		Config *RoleConfig `json:"config"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *roleDNSConfigOutput) error {
		output.Config = r.cfg
		return nil
	})
	u.SetName("dns.get_role_config")
	u.SetTitle("DNS role config")
	u.SetTags("roles/dns")
	return u
}

func (r *Role) apiHandlerRoleConfigPut() usecase.Interactor {
	type roleDNSConfigInput struct {
		Config *RoleConfig `json:"config"`
	}
	u := usecase.NewInteractor(func(ctx context.Context, input roleDNSConfigInput, output *struct{}) error {
		jc, err := json.Marshal(input.Config)
		if err != nil {
			return status.Wrap(err, status.InvalidArgument)
		}
		_, err = r.i.KV().Put(ctx, r.i.KV().Key(types.KeyRole, "dns").String(), string(jc))
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
