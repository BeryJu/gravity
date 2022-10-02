package auth

import (
	"context"
	"encoding/gob"
	"encoding/json"
	"net/http"

	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/swaggest/rest/web"
	"golang.org/x/crypto/bcrypt"
)

type AuthProvider struct {
	role         roles.Role
	inst         roles.Instance
	log          *log.Entry
	oidc         *types.OIDCConfig
	inner        http.Handler
	allowedPaths []string
}

func NewAuthProvider(r roles.Role, inst roles.Instance, oidc *types.OIDCConfig) *AuthProvider {
	ap := &AuthProvider{
		role: r,
		inst: inst,
		log:  inst.Log().WithField("mw", "auth"),
		oidc: oidc,
		allowedPaths: []string{
			"/api/v1/auth/me",
			"/api/v1/auth/config",
			"/api/v1/auth/login",
		},
	}
	if ap.oidc != nil {
		ap.InitOIDC()
	}
	gob.Register(User{})
	inst.AddEventListener(types.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		mux := ev.Payload.Data["mux"].(*mux.Router)

		svc.Get("/api/v1/auth/me", ap.APIMe())
		svc.Get("/api/v1/auth/config", ap.APIConfig())

		svc.Post("/api/v1/auth/login", ap.APILogin())
		mux.Path("/auth/logout").HandlerFunc(ap.APILogout)

		svc.Get("/api/v1/auth/users", ap.APIUsersGet())
		svc.Post("/api/v1/auth/users", ap.APIUsersPut())
		svc.Delete("/api/v1/auth/users", ap.APIUsersDelete())
		svc.Get("/api/v1/auth/tokens", ap.APITokensGet())
		svc.Post("/api/v1/auth/tokens", ap.APITokensPut())
		svc.Delete("/api/v1/auth/tokens", ap.APITokensDelete())
	})
	ap.createDefaultUser()
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

func (ap *AuthProvider) AsMiddleware() func(http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		ap.inner = h
		return ap
	}
}
