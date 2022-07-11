package api

import (
	"context"
	"net/http"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"github.com/swaggest/openapi-go/openapi3"
	"github.com/swaggest/rest/nethttp"
	"github.com/swaggest/rest/web"
	swgui "github.com/swaggest/swgui/v4emb"
)

type APIRole struct {
	m    *mux.Router
	oapi *web.Service
	log  *log.Entry
	i    roles.Instance
	ctx  context.Context
}

func New(instance roles.Instance) *APIRole {
	r := &APIRole{
		log: instance.GetLogger().WithField("role", types.KeyRole),
		i:   instance,
		m:   mux.NewRouter(),
	}
	r.m.Use(NewLoggingHandler(r.log, nil))
	r.setupUI()
	return r
}

func (r *APIRole) Start(ctx context.Context, config []byte) error {
	r.ctx = ctx
	cfg := r.decodeAPIRoleConfig(config)
	r.prepareOpenAPI()
	listen := extconfig.Get().Listen(cfg.Port)
	r.log.WithField("listen", listen).Info("Starting API Server")
	return http.ListenAndServe(listen, r.m)
}

func (r *APIRole) prepareOpenAPI() {
	if r.oapi != nil {
		return
	}
	r.oapi = web.DefaultService()
	r.oapi.OpenAPI.Info.Title = "gravity"
	r.oapi.OpenAPI.Info.Version = extconfig.Version
	r.oapi.Use(nethttp.HTTPBasicSecurityMiddleware(r.oapi.OpenAPICollector, "Admin", "Admin access"))
	// r.oapi.OpenAPICollector.Collect()
	r.oapi.Docs("/api/v1/docs", swgui.New)

	apiRouter := r.m.PathPrefix("/api").Name("api").Subrouter()
	apiRouter.Use(NewAuthMiddleware(r))
	apiRouter.PathPrefix("/v1").Handler(r.oapi)

	r.i.DispatchEvent(types.EventTopicAPIMuxSetup, roles.NewEvent(map[string]interface{}{
		"svc": r.oapi,
	}))
}

func (r *APIRole) Schema() *openapi3.Spec {
	r.prepareOpenAPI()
	return r.oapi.OpenAPICollector.Reflector().Spec
}

func (r *APIRole) Stop() {
}
