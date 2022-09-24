package auth

import (
	"context"
	"encoding/gob"
	"encoding/json"

	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api/types"
	log "github.com/sirupsen/logrus"
	"github.com/swaggest/rest/web"
	"golang.org/x/crypto/bcrypt"
)

type AuthProvider struct {
	role roles.Role
	inst roles.Instance
	log  *log.Entry
	oidc *types.OIDCConfig
}

func NewAuthProvider(r roles.Role, inst roles.Instance, oidc *types.OIDCConfig) *AuthProvider {
	ap := &AuthProvider{
		role: r,
		inst: inst,
		log:  inst.Log().WithField("mw", "auth"),
		oidc: oidc,
	}
	if ap.oidc != nil {
		ap.InitOIDC()
	}
	gob.Register(User{})
	inst.AddEventListener(types.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		svc.Get("/api/v1/auth/me", ap.apiHandlerAuthMe())
		svc.Get("/api/v1/auth/config", ap.apiHandlerAuthConfig())
		svc.Post("/api/v1/auth/login", ap.apiHandlerAuthLogin())
		svc.Get("/api/v1/auth/users", ap.apiHandlerAuthUserRead())
		svc.Post("/api/v1/auth/users", ap.apiHandlerAuthUserPut())
		svc.Delete("/api/v1/auth/users", ap.apiHandlerAuthUserDelete())
	})
	return ap
}

func (ap *AuthProvider) CreateUser(ctx context.Context, username, password string) error {
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := User{
		Password: string(hashedPw),
	}
	userJson, err := json.Marshal(user)
	if err != nil {
		return err
	}

	_, err = ap.inst.KV().Put(
		ctx,
		ap.inst.KV().Key(
			types.KeyRole,
			types.KeyUsers,
			username,
		).String(),
		string(userJson),
	)
	if err != nil {
		return err
	}
	return nil
}
