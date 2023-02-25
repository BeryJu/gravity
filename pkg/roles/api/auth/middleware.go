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
	if ap.isAllowedPath(r) {
		return true
	}
	if ap.checkToken(r) {
		return true
	}
	session := r.Context().Value(types.RequestSession).(*sessions.Session)
	u, ok := session.Values[types.SessionKeyUser]
	if u != nil && ok {
		hub := sentry.GetHubFromContext(r.Context())
		if hub == nil {
			hub = sentry.CurrentHub()
		}
		hub.Scope().SetUser(sentry.User{
			Username: u.(User).Username,
		})
		return true
	}
	return false
}

func (ap *AuthProvider) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	if ap.isRequestAllowed(r) {
		ap.inner.ServeHTTP(rw, r)
		return
	}
	http.Error(rw, "unauthenticated", http.StatusUnauthorized)
}
