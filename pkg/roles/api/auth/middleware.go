package auth

import (
	"net/http"

	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/sessions"
)

func (ap *AuthProvider) isAllowedPath(r *http.Request) bool {
	for _, path := range ap.allowedPaths {
		if r.URL.Path == path {
			return true
		}
	}
	return false
}

func (ap *AuthProvider) isRequestAllowed(r *http.Request) bool {
	ap.checkStaticToken(r)
	if ap.oidc != nil {
		ap.checkJWTToken(r)
	}
	if ap.isAllowedPath(r) {
		return true
	}
	session := r.Context().Value(types.RequestSession).(*sessions.Session)
	u, ok := session.Values[types.SessionKeyUser]
	if u == nil || !ok {
		return false
	}
	hub := sentry.GetHubFromContext(r.Context())
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	uu := u.(*types.User)
	hub.Scope().SetUser(sentry.User{
		Username: uu.Username,
	})
	return ap.checkPermission(r, uu)
}

func (ap *AuthProvider) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if ap.isRequestAllowed(r) {
		ap.inner.ServeHTTP(rw, r)
		return
	}
	http.Error(rw, "unauthenticated", http.StatusUnauthorized)
}
