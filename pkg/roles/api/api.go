package api

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	"beryju.io/ddet/pkg/extconfig"
	"beryju.io/ddet/pkg/roles"
	"beryju.io/ddet/pkg/roles/api/types"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/swaggest/rest/nethttp"
	"github.com/swaggest/rest/web"
	"golang.org/x/crypto/bcrypt"
)

type APIRole struct {
	m   *mux.Router
	log *log.Entry
	i   roles.Instance
	ctx context.Context
}

func New(instance roles.Instance) *APIRole {
	r := &APIRole{
		log: instance.GetLogger().WithField("role", types.KeyRole),
		i:   instance,
		m:   mux.NewRouter(),
	}
	go r.i.AddEventListener(types.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		if !extconfig.Get().Debug {
			return
		}
		mux := ev.Payload.Data["mux"].(*mux.Router).Name("roles.api").Subrouter()
		mux.Name("v0.debug").Path("/v0/debug").Methods(http.MethodGet).HandlerFunc(r.apiHandlerDebugGet)
		mux.Name("v0.debug").Path("/v0/debug").Methods(http.MethodPost).HandlerFunc(r.apiHandlerDebugPost)
		mux.Name("v0.debug").Path("/v0/debug").Methods(http.MethodDelete).HandlerFunc(r.apiHandlerDebugDel)
	})
	go r.setupUI()
	return r
}

func (r *APIRole) Start(ctx context.Context, config []byte) error {
	r.ctx = ctx
	cfg := r.decodeAPIRoleConfig(config)

	r.m.Use(NewLoggingHandler(r.log, nil))

	apiRouter := r.m.PathPrefix("/api").Name("api").Subrouter()
	apiRouter.Use(NewAuthMiddleware(r))
	apiRouter.Use(func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Accept", "application/json")
			w.Header().Set("Content-Type", "application/json")
			h.ServeHTTP(w, r)
		})
	})
	// Auto-set some common headers
	service := web.DefaultService()
	service.OpenAPI.Info.Title = "ye"
	service.OpenAPI.Info.Version = "v1.0.0"
	adminSecuritySchema := nethttp.HTTPBasicSecurityMiddleware(service.OpenAPICollector, "Admin", "Admin access")
	service.Use(adminSecuritySchema)
	// apiRouter.Handle("/", service)
	r.i.DispatchEvent(types.EventTopicAPIMuxSetup, roles.NewEvent(map[string]interface{}{
		"mux": apiRouter,
		"svc": service,
	}))
	r.log.Debug("Registered routes:")
	r.m.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
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

	listen := extconfig.Get().Listen(cfg.Port)
	r.log.WithField("listen", listen).Info("Starting API Server")
	return http.ListenAndServe(listen, r.m)
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
