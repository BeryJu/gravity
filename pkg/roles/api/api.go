package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"beryju.io/ddet/pkg/extconfig"
	"beryju.io/ddet/pkg/roles"
	"beryju.io/ddet/pkg/roles/api/types"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

type APIRole struct {
	log *log.Entry
	i   roles.Instance
}

func New(instance roles.Instance) *APIRole {
	return &APIRole{
		log: log.WithField("role", types.KeyRole),
		i:   instance,
	}
}

func (r *APIRole) Start(config []byte) error {
	cfg := r.decodeAPIRoleConfig(config)

	mux := mux.NewRouter()
	mux.Use(NewLoggingHandler(r.log, nil))
	mux.Use(NewAuthMiddleware(r))
	mux.Name("ddet.api.v0.test").Path("/api/v0/test").Methods("GET").HandlerFunc(r.handleAPIGet)
	mux.Name("ddet.api.v0.test").Path("/api/v0/test").Methods("POST").HandlerFunc(r.handleAPIPost)
	mux.Name("ddet.api.v0.test").Path("/api/v0/test").Methods("DELETE").HandlerFunc(r.handleAPIDel)

	r.log.WithField("port", cfg.Port).Info("Starting API Server")
	return http.ListenAndServe(fmt.Sprintf("%s:%d", extconfig.Get().Instance.IP, cfg.Port), mux)
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

	_, err = r.i.GetKV().Put(
		context.TODO(),
		r.i.GetKV().Key(
			types.KeyRole,
			types.KeyUsers,
			username,
		),
		string(userJson),
	)
	if err != nil {
		return err
	}
	return nil
}

func (r *APIRole) Stop() {
}
