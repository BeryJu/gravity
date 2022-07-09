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
	m   *chi.Mux
	log *log.Entry
	i   roles.Instance
}

func New(instance roles.Instance) *APIRole {
	r := &APIRole{
		log: instance.GetLogger().WithField("role", types.KeyRole),
		i:   instance,
		m:   chi.NewRouter(),
	}
	r.i.AddEventListener(types.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		if !extconfig.Get().Debug {
			return
		}
		mux := ev.Payload.Data["mux"].(*chi.Mux)
		mux.Get("/v0/debug", r.apiHandlerDebugGet)
		mux.Post("/v0/debug", r.apiHandlerDebugPost)
		mux.Delete("/v0/debug", r.apiHandlerDebugDel)
	})
	r.setupUI()
	return r
}

func (r *APIRole) Start(config []byte) error {
	cfg := r.decodeAPIRoleConfig(config)

	r.m.Use(NewLoggingHandler(r.log, nil))

	r.m.Route("/api", func(ro chi.Router) {
		ro.Use(NewAuthMiddleware(r))
		// Auto-set some common headers
		ro.Use(func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Accept", "application/json")
				w.Header().Set("Content-Type", "application/json")
				h.ServeHTTP(w, r)
			})
		})
		r.i.DispatchEvent(types.EventTopicAPIMuxSetup, roles.NewEvent(map[string]interface{}{
			"mux": ro,
		}))
	})
	r.log.Debug("Registered routes:")
	chi.Walk(r.m, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		r.log.WithFields(log.Fields{
			"method": method,
		}).Debug(route)
		return nil
	})

	r.log.WithField("port", cfg.Port).Info("Starting API Server")
	return http.ListenAndServe(fmt.Sprintf("%s:%d", extconfig.Get().Instance.IP, cfg.Port), r.m)
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
