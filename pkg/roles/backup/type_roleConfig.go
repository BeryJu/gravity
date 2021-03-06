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

func (r *Role) apiHandlerRoleConfigGet() usecase.IOInteractor {
	type roleBackupConfigOutput struct {
		Config *RoleConfig `json:"config"`
	}
	u := usecase.NewIOI(new(struct{}), new(roleBackupConfigOutput), func(ctx context.Context, input, output interface{}) error {
		var (
			out = output.(*roleBackupConfigOutput)
		)
		out.Config = r.cfg
		return nil
	})
	u.SetName("backup.get_role_config")
	u.SetTitle("Backup role config")
	u.SetTags("roles/backup")
	return u
}

func (r *Role) apiHandlerRoleConfigPut() usecase.IOInteractor {
	type roleBackupConfigInput struct {
		Config *RoleConfig `json:"config"`
	}
	u := usecase.NewIOI(new(roleBackupConfigInput), new(struct{}), func(ctx context.Context, input, output interface{}) error {
		var (
			in = input.(*roleBackupConfigInput)
		)
		jc, err := json.Marshal(in.Config)
		if err != nil {
			return status.Wrap(err, status.InvalidArgument)
		}
		_, err = r.i.KV().Put(ctx, r.i.KV().Key(types.KeyRole, "backup").String(), string(jc))
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
