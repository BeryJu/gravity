package api

import (
	"context"
	"encoding/json"
	"net/http"

	"beryju.io/gravity/pkg/roles/api/types"
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
			log:   r.log.WithField("mw", "auth"),
			inner: h,
		}
	}
}

func (am *AuthMiddleware) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	username, password, ok := r.BasicAuth()
	if !ok {
		rw.Header().Set("WWW-Authenticate", "Basic realm=gravity")
		http.Error(rw, "missing http basic authentication", http.StatusUnauthorized)
		return
	}

	rawUsers, err := am.role.i.KV().Get(
		r.Context(),
		am.role.i.KV().Key(
			types.KeyRole,
			types.KeyUsers,
			username,
		).String(),
	)
	if err != nil {
		bcrypt.CompareHashAndPassword([]byte{}, []byte(password))
		am.log.WithError(err).Warning("failed to get users")
		rw.WriteHeader(500)
		return
	}
	if len(rawUsers.Kvs) < 1 {
		bcrypt.CompareHashAndPassword([]byte{}, []byte(password))
		http.Error(rw, "invalid http basic authentication", http.StatusUnauthorized)
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

func (r *APIRole) CreateUser(username, password string) error {
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

	_, err = r.i.KV().Put(
		context.TODO(),
		r.i.KV().Key(
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
