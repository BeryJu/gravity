package auth

import (
	"net/http"

	"beryju.io/gravity/pkg/roles/api/types"
	"golang.org/x/crypto/bcrypt"
)

func (ap *AuthProvider) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		rw.Header().Set("WWW-Authenticate", "Basic realm=gravity")
		http.Error(rw, "missing http basic authentication", http.StatusUnauthorized)
		return
	}

	rawUsers, err := ap.inst.KV().Get(
		r.Context(),
		ap.inst.KV().Key(
			types.KeyRole,
			types.KeyUsers,
			username,
		).String(),
	)
	if err != nil {
		bcrypt.CompareHashAndPassword([]byte{}, []byte(password))
		ap.log.WithError(err).Warning("failed to get users")
		rw.WriteHeader(500)
		return
	}
	if len(rawUsers.Kvs) < 1 {
		bcrypt.CompareHashAndPassword([]byte{}, []byte(password))
		http.Error(rw, "invalid http basic authentication", http.StatusUnauthorized)
		return
	}
	user, err := ap.userFromKV(rawUsers.Kvs[0])
	if err != nil {
		bcrypt.CompareHashAndPassword([]byte{}, []byte(password))
		ap.log.WithField("user", username).WithError(err).Warning("failed to parse user")
		rw.WriteHeader(500)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		ap.log.WithField("user", username).Warning("invalid credentials")
		http.Error(rw, "invalid http basic authentication", http.StatusUnauthorized)
		return
	}

	ap.inner.ServeHTTP(rw, r)
}
