package api

import (
	"fmt"
	"net/http"

	"github.com/beryju/dns-dhcp-etcd-thingy/internal/roles"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

const (
	RoleAPIPrefix = "api"
)

type APIServerRole struct {
	log *log.Entry
	i   roles.Instance
}

func New(instance roles.Instance) *APIServerRole {
	return &APIServerRole{
		log: log.WithField("role", "api"),
		i:   instance,
	}
}

func (r *APIServerRole) Start(config []byte) error {
	cfg := r.decodeAPIRoleConfig(config)

	mux := mux.NewRouter()
	mux.Use(NewLoggingHandler(r.log, nil))
	mux.Name("ddet.api.v0.test").Path("/api/v0/test").Methods("GET").HandlerFunc(r.handleAPIGet)
	mux.Name("ddet.api.v0.test").Path("/api/v0/test").Methods("POST").HandlerFunc(r.handleAPIPost)
	mux.Name("ddet.api.v0.test").Path("/api/v0/test").Methods("DELETE").HandlerFunc(r.handleAPIDel)

	r.log.WithField("port", cfg.Port).Info("Starting API Server")
	return http.ListenAndServe(fmt.Sprintf(":%d", cfg.Port), mux)
}

func (r *APIServerRole) Stop() {
}

func (r *APIServerRole) HandleEvent(ev *roles.Event[any]) {
}
