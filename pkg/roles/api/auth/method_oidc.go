package auth

import (
	"context"
	"encoding/base64"
	"net/http"
	"strings"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"golang.org/x/oauth2"
)

func (ap *AuthProvider) InitOIDC() {
	provider, err := oidc.NewProvider(context.Background(), ap.oidc.Issuer)
	if err != nil {
		ap.log.WithError(err).Warning("failed to initialise oidc")
		ap.oidc = nil
		return
	}
	red := strings.ReplaceAll(ap.oidc.RedirectURL, "$INSTANCE_IDENTIFIER", extconfig.Get().Instance.Identifier)
	red = strings.ReplaceAll(red, "$INSTANCE_IP", extconfig.Get().Instance.IP)
	oauth2Config := oauth2.Config{
		ClientID:     ap.oidc.ClientID,
		ClientSecret: ap.oidc.ClientSecret,
		RedirectURL:  red,
		Endpoint:     provider.Endpoint(),
		Scopes:       ap.oidc.Scopes,
	}
	verifier := provider.Verifier(&oidc.Config{ClientID: ap.oidc.ClientID})

	ap.inst.AddEventListener(types.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		mux := ev.Payload.Data["mux"].(*mux.Router)
		mux.Path("/auth/oidc").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			newState := base64.RawURLEncoding.EncodeToString(securecookie.GenerateRandomKey(32))
			session := r.Context().Value(types.RequestSession).(*sessions.Session)
			session.Values[types.SessionKeyOIDCState] = newState
			http.Redirect(w, r, oauth2Config.AuthCodeURL(newState), http.StatusFound)
		})
		mux.Path("/auth/oidc/callback").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session := r.Context().Value(types.RequestSession).(*sessions.Session)
			state, ok := session.Values[types.SessionKeyOIDCState]
			if !ok || state.(string) != r.URL.Query().Get("state") {
				http.Error(w, "invalid state", http.StatusBadRequest)
				return
			}

			oauth2Token, err := oauth2Config.Exchange(r.Context(), r.URL.Query().Get("code"))
			if err != nil {
				ap.log.WithError(err).Warning("failed to exchange code")
				http.Error(w, "failed to authenticate", http.StatusBadRequest)
				return
			}

			rawIDToken, ok := oauth2Token.Extra("id_token").(string)
			if !ok {
				ap.log.Warning("no id_token")
				http.Error(w, "failed to authenticate", http.StatusBadRequest)
				return
			}

			idToken, err := verifier.Verify(r.Context(), rawIDToken)
			if err != nil {
				ap.log.WithError(err).Warning("failed to verify id_token")
				http.Error(w, "failed to authenticate", http.StatusBadRequest)
				return
			}

			var claims struct {
				Email string `json:"email"`
			}
			if err := idToken.Claims(&claims); err != nil {
				ap.log.WithError(err).Warning("failed to get claims")
				http.Error(w, "failed to authenticate", http.StatusBadRequest)
				return
			}
			user := User{
				Username: claims.Email,
				Password: "",
			}
			session.Values[types.SessionKeyUser] = user
			http.Redirect(w, r, "/", http.StatusFound)
		})
	})
}
