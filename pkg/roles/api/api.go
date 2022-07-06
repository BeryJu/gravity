package api

import (
	"fmt"
	"net/http"

	"beryju.io/ddet/pkg/extconfig"
	"beryju.io/ddet/pkg/roles"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	RoleAPIPrefix = "api"
)

type APIRole struct {
	log *log.Entry
	i   roles.Instance
}

func New(instance roles.Instance) *APIRole {
	return &APIRole{
		log: log.WithField("role", "api"),
		i:   instance,
	}
}

func (r *APIRole) Start(config []byte) error {
	cfg := r.decodeAPIRoleConfig(config)

	mux := mux.NewRouter()
	mux.Use(NewLoggingHandler(r.log, nil))
	mux.Name("ddet.api.v0.test").Path("/api/v0/test").Methods("GET").HandlerFunc(r.handleAPIGet)
	mux.Name("ddet.api.v0.test").Path("/api/v0/test").Methods("POST").HandlerFunc(r.handleAPIPost)
	mux.Name("ddet.api.v0.test").Path("/api/v0/test").Methods("DELETE").HandlerFunc(r.handleAPIDel)

	r.log.WithField("port", cfg.Port).Info("Starting API Server")
	return http.ListenAndServe(fmt.Sprintf("%s:%d", extconfig.Get().Instance.IP, cfg.Port), mux)
}

func (r *APIRole) Stop() {
}
