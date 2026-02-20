package auth

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"golang.org/x/oauth2"
)

func (ap *AuthProvider) ConfigureOpenIDConnect(ctx context.Context, config *types.OIDCConfig) error {
	c := &http.Client{Transport: extconfig.Transport()}
	ctx = context.WithValue(ctx, oauth2.HTTPClient, c)
	provider, err := oidc.NewProvider(ctx, config.Issuer)
	if err != nil {
		return errors.Wrap(err, "failed to initialise oidc")
	}
	ap.oidc = config
	red := strings.ReplaceAll(config.RedirectURL, "$INSTANCE_IDENTIFIER", extconfig.Get().Instance.Identifier)
	red = strings.ReplaceAll(red, "$INSTANCE_IP", extconfig.Get().Instance.IP)
	ap.oidcConfig = oauth2.Config{
		ClientID:     config.ClientID,
		ClientSecret: config.ClientSecret,
		RedirectURL:  red,
		Endpoint:     provider.Endpoint(),
		Scopes:       config.Scopes,
	}
	ap.oidcVerifier = provider.Verifier(&oidc.Config{ClientID: config.ClientID})
	return nil
}

func (ap *AuthProvider) oidcInit(w http.ResponseWriter, r *http.Request) {
	if ap.oidc == nil {
		w.WriteHeader(400)
		_, _ = w.Write([]byte("OIDC not configured"))
		return
	}
	newState := base64.RawURLEncoding.EncodeToString(securecookie.GenerateRandomKey(32))
	session := r.Context().Value(types.RequestSession).(*sessions.Session)
	session.Values[types.SessionKeyOIDCState] = newState
	session.Values[types.SessionKeyDirty] = true
	http.Redirect(w, r, ap.oidcConfig.AuthCodeURL(newState), http.StatusFound)
}

func (ap *AuthProvider) oidcCallback(w http.ResponseWriter, r *http.Request) {
	if ap.oidc == nil {
		w.WriteHeader(400)
		_, _ = w.Write([]byte("OIDC not configured"))
		return
	}
	session := r.Context().Value(types.RequestSession).(*sessions.Session)
	state, ok := session.Values[types.SessionKeyOIDCState]
	if !ok || state.(string) != r.URL.Query().Get("state") {
		http.Error(w, "invalid state", http.StatusBadRequest)
		return
	}

	oauth2Token, err := ap.oidcConfig.Exchange(r.Context(), r.URL.Query().Get("code"))
	if err != nil {
		ap.log.Warn("failed to exchange code", zap.Error(err))
		http.Error(w, "failed to authenticate", http.StatusBadRequest)
		return
	}

	rawIDToken, ok := oauth2Token.Extra("id_token").(string)
	if !ok {
		ap.log.Warn("no id_token")
		http.Error(w, "failed to authenticate", http.StatusBadRequest)
		return
	}

	idToken, err := ap.oidcVerifier.Verify(r.Context(), rawIDToken)
	if err != nil {
		ap.log.Warn("failed to verify id_token", zap.Error(err))
		http.Error(w, "failed to authenticate", http.StatusBadRequest)
		return
	}

	var claims struct {
		Email string `json:"email"`
	}
	if err := idToken.Claims(&claims); err != nil {
		ap.log.Warn("failed to get claims", zap.Error(err))
		http.Error(w, "failed to authenticate", http.StatusBadRequest)
		return
	}
	user := &types.User{
		Username: claims.Email,
		Password: "",
		Permissions: []*types.Permission{
			{
				Path:    "/*",
				Methods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodHead, http.MethodDelete},
			},
		},
	}
	session.Values[types.SessionKeyUser] = user
	session.Values[types.SessionKeyDirty] = true
	http.Redirect(w, r, "/ui/#/overview", http.StatusFound)
}

func (ap *AuthProvider) checkJWTToken(r *http.Request) bool {
	header := r.Header.Get(AuthorizationHeader)
	if header == "" {
		return false
	}
	parts := strings.SplitN(header, " ", 2)
	if len(parts) < 2 {
		return false
	}
	if !strings.EqualFold(parts[0], BearerType) {
		return false
	}
	t, err := ap.oidcVerifier.Verify(r.Context(), parts[1])
	if err != nil {
		ap.log.Warn("failed to verify JWT token", zap.Error(err))
		return false
	}
	rt := map[string]interface{}{}
	err = t.Claims(&rt)
	if err != nil {
		return false
	}
	// Get token's user
	rawUsers, err := ap.inst.KV().Get(
		r.Context(),
		ap.inst.KV().Key(
			types.KeyRole,
			types.KeyUsers,
			rt[ap.oidc.TokenUsernameField].(string),
		).String(),
	)
	if err != nil {
		ap.log.Warn("failed to check token", zap.Error(err))
		return false
	}
	if len(rawUsers.Kvs) < 1 {
		return false
	}
	user, err := ap.userFromKV(rawUsers.Kvs[0])
	if err != nil {
		return false
	}
	session := r.Context().Value(types.RequestSession).(*sessions.Session)
	session.Values[types.SessionKeyUser] = user
	return false
}
