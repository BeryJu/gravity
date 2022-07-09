package api

import (
	"encoding/json"
	"net/http"

	"beryju.io/ddet/pkg/roles/api/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Password string `json:"password"`
}

type AuthMiddleware struct {
	role  *APIRole
	log   *log.Entry
	inner http.Handler
}

func NewAuthMiddleware(r *APIRole) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &AuthMiddleware{
			role:  r,
			log:   log.WithField("role", "api").WithField("mw", "auth"),
			inner: h,
		}
	}
}

func (am *AuthMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		http.Error(rw, "missing http basic authentication", http.StatusUnauthorized)
		return
	}

	rawUsers, err := am.role.i.GetKV().Get(
		r.Context(),
		am.role.i.GetKV().Key(
			types.KeyRole,
			types.KeyUsers,
			username,
		),
	)
	if len(rawUsers.Kvs) < 1 || err != nil {
		bcrypt.CompareHashAndPassword([]byte{}, []byte(password))
		am.log.WithField("user", username).WithError(err).Warning("failed to get users")
		rw.WriteHeader(500)
		return
	}
	user := User{}
	err = json.Unmarshal(rawUsers.Kvs[0].Value, &user)
	if err != nil {
		bcrypt.CompareHashAndPassword([]byte{}, []byte(password))
		am.log.WithField("user", username).WithError(err).Warning("failed to parse user")
		rw.WriteHeader(500)
		return
	}
	if bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)) != nil {
		am.log.WithField("user", username).Warning("invalid credentials")
		http.Error(rw, "invalid http basic authentication", http.StatusUnauthorized)
		return
	}

	am.inner.ServeHTTP(rw, r)
}
