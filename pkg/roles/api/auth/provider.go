package auth

import (
	"context"
	"encoding/gob"
	"net/http"

	instanceTypes "beryju.io/gravity/pkg/instance/types"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/mux"
	"github.com/swaggest/rest/web"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/oauth2"
)

type AuthProvider struct {
	role         roles.Role
	inst         roles.Instance
	log          *zap.Logger
	oidc         *types.OIDCConfig
	oidcConfig   oauth2.Config
	oidcVerifier *oidc.IDTokenVerifier
	inner        http.Handler
	allowedPaths []string
}

func NewAuthProvider(r roles.Role, inst roles.Instance) *AuthProvider {
	ap := &AuthProvider{
		role: r,
		inst: inst,
		log:  inst.Log().Named("role.api.auth"),
		allowedPaths: []string{
			"/api/v1/auth/me",
			"/api/v1/auth/config",
			"/api/v1/auth/login",
			"/api/v1/openapi.json",
		},
	}
	gob.Register(types.User{})
	inst.AddEventListener(types.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		mux := ev.Payload.Data["mux"].(*mux.Router)

		svc.Get("/api/v1/auth/me", ap.APIMe())
		svc.Get("/api/v1/auth/config", ap.APIConfig())

		svc.Post("/api/v1/auth/login", ap.APILogin())
		mux.Path("/auth/logout").HandlerFunc(ap.APILogout)
		mux.Path("/auth/oidc").HandlerFunc(ap.oidcInit)
		mux.Path("/auth/oidc/callback").HandlerFunc(ap.oidcCallback)

		svc.Get("/api/v1/auth/users", ap.APIUsersGet())
		svc.Post("/api/v1/auth/users", ap.APIUsersPut())
		svc.Delete("/api/v1/auth/users", ap.APIUsersDelete())
		svc.Get("/api/v1/auth/tokens", ap.APITokensGet())
		svc.Post("/api/v1/auth/tokens", ap.APITokensPut())
		svc.Delete("/api/v1/auth/tokens", ap.APITokensDelete())
	})
	inst.AddEventListener(instanceTypes.EventTopicInstanceFirstStart, ap.FirstStart)
	return ap
}

func (ap *AuthProvider) CreateUser(ctx context.Context, username, password string) error {
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := types.User{
		Password: string(hashedPw),
		Permissions: []*types.Permission{
			{
				Path:    "/*",
				Methods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodHead, http.MethodDelete},
			},
		},
	}

	_, err = ap.inst.KV().PutObj(
		ctx,
		ap.inst.KV().Key(
			types.KeyRole,
			types.KeyUsers,
			username,
		).String(),
		&user,
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
