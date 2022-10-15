package tsdb

import (
	"context"
	"encoding/json"

	instanceTypes "beryju.io/gravity/pkg/instance/types"
	"go.uber.org/zap"

	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

const (
	KeyRole = "tsdb"
)

type RoleConfig struct {
	Enabled bool  `json:"enabled"`
	Expire  int64 `json:"expire"`
	Scrape  int64 `json:"scrape"`
}

func (r *Role) decodeRoleConfig(raw []byte) *RoleConfig {
	def := RoleConfig{
		Enabled: true,
		Expire:  60 * 30,
		Scrape:  30,
	}
	if len(raw) < 1 {
		return &def
	}
	err := json.Unmarshal(raw, &def)
	if err != nil {
		r.log.Warn("failed to decode role config", zap.Error(err))
	}
	return &def
}

type APIRoleConfigOutput struct {
	Config RoleConfig `json:"config" required:"true"`
}

func (r *Role) APIRoleConfigGet() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *APIRoleConfigOutput) error {
		output.Config = *r.cfg
		return nil
	})
	u.SetName("tsdb.get_role_config")
	u.SetTitle("TSDB role config")
	u.SetTags("roles/tsdb")
	return u
}

type APIRoleConfigInput struct {
	Config RoleConfig `json:"config" required:"true"`
}

func (r *Role) APIRoleConfigPut() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIRoleConfigInput, output *struct{}) error {
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
	u.SetName("tsdb.put_role_config")
	u.SetTitle("TSDB role config")
	u.SetTags("roles/tsdb")
	u.SetExpectedErrors(status.InvalidArgument, status.Internal)
	return u
}
