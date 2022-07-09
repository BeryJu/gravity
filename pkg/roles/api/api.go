package api

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"beryju.io/ddet/pkg/extconfig"
	"beryju.io/ddet/pkg/roles"
	"beryju.io/ddet/pkg/roles/api/types"
	"github.com/go-chi/chi/v5"
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
		mux := ev.Payload.Data["mux"].(*chi.Mux)
		mux.Get("/api/v0/debug", r.apiHandlerDebugGet)
		mux.Post("/api/v0/debug", r.apiHandlerDebugPost)
		mux.Delete("/api/v0/debug", r.apiHandlerDebugDel)
	})
	return r
}

func (r *APIRole) Start(config []byte) error {
	cfg := r.decodeAPIRoleConfig(config)

	m := chi.NewRouter()
	m.Use(NewLoggingHandler(r.log, nil))
	m.Use(NewAuthMiddleware(r))

	r.i.DispatchEvent(types.EventTopicAPIMuxSetup, roles.NewEvent(map[string]interface{}{
		"mux": m,
	}))
	r.log.Debug("Registered routes:")
	chi.Walk(m, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		r.log.WithFields(log.Fields{
			// "routeName": ,
			"method": method,
		}).Debug(route)
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
