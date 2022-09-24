package api

import (
	"context"
	"net/http"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api/auth"
	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/gorilla/mux"
	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	log "github.com/sirupsen/logrus"
	"github.com/swaggest/openapi-go/openapi3"
	"github.com/swaggest/rest/web"
	swgui "github.com/swaggest/swgui/v4emb"
)

type Role struct {
	m        *mux.Router
	oapi     *web.Service
	log      *log.Entry
	i        roles.Instance
	ctx      context.Context
	cfg      *RoleConfig
	sessions sessions.Store
}

func New(instance roles.Instance) *Role {
	sess := sessions.NewCookieStore(securecookie.GenerateRandomKey(32))
	r := &Role{
		log:      instance.Log(),
		i:        instance,
		m:        mux.NewRouter(),
		sessions: sess,
		cfg:      &RoleConfig{},
	}
	r.m.Use(NewRecoverMiddleware(r.log))
	r.m.Use(NewLoggingMiddleware(r.log, nil))
	r.m.Use(r.SessionMiddleware)
	r.setupUI()
	r.i.AddEventListener(types.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		svc.Get("/api/v1/roles/api", r.apiHandlerRoleConfigGet())
		svc.Post("/api/v1/roles/api", r.apiHandlerRoleConfigPut())
	})
	return r
}

func (r *Role) Start(ctx context.Context, config []byte) error {
	r.ctx = ctx
	r.cfg = r.decodeRoleConfig(config)
	r.prepareOpenAPI(ctx)
	listen := extconfig.Get().Listen(r.cfg.Port)
	r.log.WithField("listen", listen).Info("starting API Server")
	return http.ListenAndServe(listen, r.m)
}

func (r *Role) prepareOpenAPI(ctx context.Context) {
	if r.oapi != nil {
		return
	}
	r.oapi = web.DefaultService()
	r.oapi.OpenAPI.Info.Title = "gravity"
	r.oapi.OpenAPI.Info.Version = extconfig.Version
	r.oapi.Docs("/api/v1/docs", swgui.New)

	apiRouter := r.m.PathPrefix("/api").Name("api").Subrouter()
	auth.NewAuthProvider(r, r.i, r.cfg.OIDC)
	apiRouter.PathPrefix("/v1").Handler(r.oapi)

	r.i.DispatchEvent(types.EventTopicAPIMuxSetup, roles.NewEvent(ctx, map[string]interface{}{
		"svc":     r.oapi,
		"mux":     r.m,
		"session": r.sessions,
	}))
}

func (r *Role) Schema() *openapi3.Spec {
	r.prepareOpenAPI(context.Background())
	return r.oapi.OpenAPICollector.Reflector().Spec
}

func (r *Role) Stop() {
}
