package api

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	instanceTypes "beryju.io/gravity/pkg/instance/types"
	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/gorilla/securecookie"
	"github.com/swaggest/usecase"
	"github.com/swaggest/usecase/status"
	"go.uber.org/zap"
)

type RoleConfig struct {
	OIDC *types.OIDCConfig `json:"oidc"`
	// Override listen address temporarily, must by JSON accessible as config is passed as JSON
	ListenOverride  string `json:"listenOverride,omitempty"`
	CookieSecret    string `json:"cookieSecret"`
	Port            int32  `json:"port"`
	SessionDuration string `json:"sessionDuration"`
}

func (r *Role) checkCookieSecret(cfg *RoleConfig, fallback string, ctx context.Context) {
	if cfg.CookieSecret != "" {
		return
	}
	cfg.CookieSecret = fallback
	r.log.Info("cookie secret not in config, generating one")
	go func(cfg *RoleConfig) {
		jc, err := json.Marshal(cfg)
		if err != nil {
			r.log.Warn("failed to json parse config", zap.Error(err))
			return
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
			r.log.Warn("failed to save config", zap.Error(err))
			return
		}
	}(cfg)
}

func (r *Role) decodeRoleConfig(ctx context.Context, raw []byte) *RoleConfig {
	fallbackSecret := base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32))
	def := RoleConfig{
		Port: 8008,
	}
	if len(raw) < 1 {
		r.checkCookieSecret(&def, fallbackSecret, ctx)
		return &def
	}
	err := json.Unmarshal(raw, &def)
	if err != nil {
		r.log.Warn("failed to decode role config", zap.Error(err))
	}
	r.checkCookieSecret(&def, fallbackSecret, ctx)
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
	u.SetName("api.get_role_config")
	u.SetTitle("API role config")
	u.SetTags("roles/api")
	return u
}

type APIRoleConfigInput struct {
	Config RoleConfig `json:"config" required:"true"`
}

func (r *Role) APIRoleConfigPut() usecase.Interactor {
	u := usecase.NewInteractor(func(ctx context.Context, input APIRoleConfigInput, output *struct{}) error {
		if input.Config.CookieSecret == "" {
			input.Config.CookieSecret = r.cfg.CookieSecret
		}
		if input.Config.SessionDuration != "" {
			if _, err := time.ParseDuration(input.Config.SessionDuration); err != nil {
				return status.InvalidArgument
			}
		}
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
	u.SetName("api.put_role_config")
	u.SetTitle("API role config")
	u.SetTags("roles/api")
	u.SetExpectedErrors(status.InvalidArgument, status.Internal)
	return u
}
