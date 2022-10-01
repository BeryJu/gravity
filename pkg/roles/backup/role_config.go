package backup

import (
	"context"
	"encoding/json"

	"beryju.io/gravity/pkg/instance/types"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
)

type RoleConfig struct {
	AccessKey string `json:"accessKey"`
	SecretKey string `json:"secretKey"`
	Endpoint  string `json:"endpoint"`
	Bucket    string `json:"bucket"`
	Path      string `json:"path"`
	CronExpr  string `json:"cronExpr"`
}

func (r *Role) decodeRoleConfig(raw []byte) *RoleConfig {
	def := RoleConfig{
		CronExpr: "0 */24 * * *",
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

type APIRoleBackupConfigOutput struct {
	Config RoleConfig `json:"config" required:"true"`
}

func (r *Role) APIHandlerRoleConfigGet() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input struct{}, output *APIRoleBackupConfigOutput) error {
		output.Config = *r.cfg
		return nil
	})
	u.SetName("backup.get_role_config")
	u.SetTitle("Backup role config")
	u.SetTags("roles/backup")
	return u
}

type APIRoleBackupConfigInput struct {
	Config RoleConfig `json:"config" required:"true"`
}

func (r *Role) APIHandlerRoleConfigPut() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIRoleBackupConfigInput, output *struct{}) error {
		jc, err := json.Marshal(input.Config)
		if err != nil {
			return status.Wrap(err, status.InvalidArgument)
		}
		_, err = r.i.KV().Put(ctx, r.i.KV().Key(types.KeyRole, KeyRole).String(), string(jc))
		if err != nil {
			return status.Wrap(err, status.Internal)
		}
		return nil
	})
	u.SetName("backup.put_role_config")
	u.SetTitle("Backup role config")
	u.SetTags("roles/backup")
	u.SetExpectedErrors(status.InvalidArgument, status.Internal)
	return u
}
