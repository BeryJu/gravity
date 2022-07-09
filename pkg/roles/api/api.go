package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

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
	r := &APIRole{
		log: instance.GetLogger().WithField("role", types.KeyRole),
		i:   instance,
	}
	r.i.AddEventListener(types.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		if !extconfig.Get().Debug {
			return
		}
		mux := ev.Payload.Data["mux"].(*mux.Router).Name("roles.api").Subrouter()
		mux.Name("v0.debug").Path("/api/v0/debug").Methods("GET").HandlerFunc(r.apiHandlerDebugGet)
		mux.Name("v0.debug").Path("/api/v0/debug").Methods("POST").HandlerFunc(r.apiHandlerDebugPost)
		mux.Name("v0.debug").Path("/api/v0/debug").Methods("DELETE").HandlerFunc(r.apiHandlerDebugDel)
	})
	return r
}

func (r *APIRole) Start(config []byte) error {
	cfg := r.decodeAPIRoleConfig(config)

	m := mux.NewRouter()
	m.Use(NewLoggingHandler(r.log, nil))
	m.Use(NewAuthMiddleware(r))

	r.i.DispatchEvent(types.EventTopicAPIMuxSetup, roles.NewEvent(map[string]interface{}{
		"mux": m,
	}))
	r.log.Debug("Registered routes:")
	m.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		fullName := []string{route.GetName()}
		for _, anc := range ancestors {
			fullName = append([]string{anc.GetName()}, fullName...)
		}
		var methods []string
		var err error
		if methods, err = route.GetMethods(); err != nil {
			methods = []string{}
		}
		var path string
		if path, err = route.GetPathTemplate(); err != nil {
			return nil
		}
		r.log.WithFields(log.Fields{
			"routeName": strings.Join(fullName, "."),
			"method":    methods,
		}).Debug(path)
		return nil
	})

	r.log.WithField("port", cfg.Port).Info("Starting API Server")
	return http.ListenAndServe(fmt.Sprintf("%s:%d", extconfig.Get().Instance.IP, cfg.Port), m)
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
