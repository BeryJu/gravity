package api

import (
	"context"
	"encoding/base64"
	"net/http"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api/auth"
	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/api7/etcdstore"
	"github.com/gorilla/mux"
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
	server   *http.Server
	auth     *auth.AuthProvider
}

func New(instance roles.Instance) *Role {
	mux := mux.NewRouter()
	r := &Role{
		log: instance.Log(),
		i:   instance,
		m:   mux,
		cfg: &RoleConfig{},
	}
	r.m.Use(NewRecoverMiddleware(r.log))
	r.m.Use(r.SessionMiddleware)
	r.m.Use(NewLoggingMiddleware(r.log, nil))
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

	cookieSecret, err := base64.StdEncoding.DecodeString(r.cfg.CookieSecret)
	if err != nil {
		return err
	}
	sess, err := etcdstore.NewEtcdStore(
		r.i.KV().Config(),
		ctx,
		r.i.KV().Key(
			extconfig.Get().Etcd.Prefix,
			types.KeyRole,
			types.KeySessions,
		).String(),
		cookieSecret,
	)
	if err != nil {
		return err
	}
	r.sessions = sess

	r.auth = auth.NewAuthProvider(r, r.i, r.cfg.OIDC)
	r.prepareOpenAPI(ctx)
	listen := extconfig.Get().Listen(r.cfg.Port)
	r.log.WithField("listen", listen).Info("starting API Server")
	r.server = &http.Server{
		Addr:    listen,
		Handler: r.m,
	}
	go func() {
		err := r.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			r.log.WithError(err).Warning("failed to listen")
		}
	}()
	return nil
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
	apiRouter.Use(r.auth.AsMiddleware())
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
	if r.server != nil {
		r.server.Shutdown(r.ctx)
	}
}
