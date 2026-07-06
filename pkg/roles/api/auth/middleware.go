package auth

import (
	"net/http"

	"github.com/getsentry/sentry-go"
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
	u := ap.GetUserFromSession(r.Context())
	if u == nil {
		return false
	}
	hub := sentry.GetHubFromContext(r.Context())
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	hub.Scope().SetUser(sentry.User{
		Username: u.Username,
	})
	return ap.checkPermission(r, u)
}

func (ap *AuthProvider) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if ap.isRequestAllowed(r) {
		ap.inner.ServeHTTP(rw, r)
		return
	}
	http.Error(rw, "unauthenticated", http.StatusUnauthorized)
}
