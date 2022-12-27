package api

import (
	"context"
	"encoding/base64"
	"errors"
	"net"
	"net/http"
	"os"
	"path"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api/auth"
	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/api7/etcdstore"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/swaggest/openapi-go/openapi3"
	"github.com/swaggest/rest/web"
	swgui "github.com/swaggest/swgui/v4emb"
	"go.uber.org/zap"
)

const (
	VAR_RUN      = "/var/run"
	GRAVITY_SOCK = "gravity.sock"
)

type Role struct {
	m            *mux.Router
	oapi         *web.Service
	log          *zap.Logger
	i            roles.Instance
	ctx          context.Context
	cfg          *RoleConfig
	sessions     sessions.Store
	httpServer   http.Server
	socketServer http.Server
	auth         *auth.AuthProvider
}

func New(instance roles.Instance) *Role {
	mux := mux.NewRouter()
	r := &Role{
		log:          instance.Log(),
		i:            instance,
		m:            mux,
		cfg:          &RoleConfig{},
		httpServer:   http.Server{},
		socketServer: http.Server{},
	}
	r.auth = auth.NewAuthProvider(r, r.i)
	r.m.Use(NewRecoverMiddleware(r.log))
	r.m.Use(r.SessionMiddleware)
	r.m.Use(NewLoggingMiddleware(r.log, nil))
	r.setupUI()
	r.i.AddEventListener(types.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		svc.Get("/api/v1/cluster/export", r.APIClusterExport())
		svc.Post("/api/v1/cluster/import", r.APIClusterImport())
		svc.Get("/api/v1/roles/api", r.APIRoleConfigGet())
		svc.Post("/api/v1/roles/api", r.APIRoleConfigPut())
		svc.Get("/api/v1/etcd/members", r.APIClusterMembers())
		svc.Post("/api/v1/etcd/join", r.APIClusterJoin())
	})
	return r
}

func (r *Role) SessionStore() sessions.Store {
	return r.sessions
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

	if r.cfg.OIDC != nil {
		err := r.auth.ConfigureOpenIDConnect(r.ctx, r.cfg.OIDC)
		if err != nil {
			r.log.Warn("failed to setup OpenID Connect, ignoring", zap.Error(err))
		}
	}
	r.prepareOpenAPI(ctx)
	go r.ListenAndServeHTTP()
	go r.ListenAndServeSocket()
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

func (r *Role) ListenAndServeHTTP() {
	r.httpServer.Handler = r.m
	listen := extconfig.Get().Listen(r.cfg.Port)
	if r.cfg.ListenOverride != "" {
		listen = r.cfg.ListenOverride
	}
	r.log.Info("starting API Server", zap.String("listen", listen))
	r.httpServer.Addr = listen
	err := r.httpServer.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		r.log.Warn("failed to listen", zap.Error(err))
	}
}

func (r *Role) ListenAndServeSocket() {
	stat, err := os.Stat(VAR_RUN)
	if errors.Is(err, os.ErrNotExist) || !stat.IsDir() {
		r.log.Info("/var/run doesn't exist or is not a dir, not starting socket API server")
		return
	}
	r.socketServer.Handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := context.WithValue(req.Context(), types.RequestSession, &sessions.Session{
			Values: map[interface{}]interface{}{
				types.SessionKeyUser: auth.User{
					Username: "gravity-socket",
				},
			},
		})
		reqAuth := req.WithContext(ctx)
		r.m.ServeHTTP(w, reqAuth)
	})
	socketPath := path.Join(VAR_RUN, GRAVITY_SOCK)
	if extconfig.Get().Debug {
		socketPath = path.Join("./", GRAVITY_SOCK)
	}
	unixListener, err := net.Listen("unix", socketPath)
	if err != nil {
		r.log.Warn("failed to listen on socket", zap.Error(err))
		return
	}
	r.log.Info("starting API Server (socket)", zap.String("listen", socketPath))
	err = r.socketServer.Serve(unixListener)
	if err != nil && err != http.ErrServerClosed {
		r.log.Warn("failed to listen", zap.Error(err))
	}
}

func (r *Role) Schema(ctx context.Context) *openapi3.Spec {
	r.prepareOpenAPI(ctx)
	return r.oapi.OpenAPICollector.Reflector().Spec
}

func (r *Role) Stop() {
	r.httpServer.Shutdown(r.ctx)
	r.socketServer.Shutdown(r.ctx)
}
