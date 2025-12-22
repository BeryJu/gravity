package api

import (
	"context"
	"encoding/base64"
	"errors"
	"io/fs"
	"net"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/avast/retry-go/v4"
	sentryhttp "github.com/getsentry/sentry-go/http"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api/auth"
	"beryju.io/gravity/pkg/roles/api/middleware"
	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/api7/etcdstore"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/swaggest/openapi-go/openapi3"
	"github.com/swaggest/rest/web"
	"go.uber.org/zap"
)

const (
	VAR_RUN      = "/var/run"
	GRAVITY_SOCK = "gravity.sock"
)

type Role struct {
	i            roles.Instance
	ctx          context.Context
	sessions     sessions.Store
	m            *mux.Router
	oapi         *web.Service
	log          *zap.Logger
	cfg          *RoleConfig
	auth         *auth.AuthProvider
	httpServer   *http.Server
	socketServer *http.Server
}

func init() {
	roles.Register("api", func(i roles.Instance) roles.Role {
		return New(i)
	})
}

func New(instance roles.Instance) *Role {
	mux := mux.NewRouter()
	r := &Role{
		log:          instance.Log(),
		i:            instance,
		m:            mux,
		cfg:          &RoleConfig{},
		httpServer:   &http.Server{},
		socketServer: &http.Server{},
		ctx:          instance.Context(),
	}
	r.auth = auth.NewAuthProvider(r, r.i)
	r.m.Use(middleware.NewRecoverMiddleware(r.log))
	r.m.Use(sentryhttp.New(sentryhttp.Options{
		Repanic: true,
	}).Handle)
	r.m.Use(r.SessionMiddleware)
	r.m.Use(middleware.NewLoggingMiddleware(r.log, nil))
	r.m.Use(NewAPIConfigMiddleware())
	r.setupUI()
	r.i.AddEventListener(types.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		svc.Get("/api/v1/cluster/node/logs", r.APIClusterNodeLogMessages())
		svc.Post("/api/v1/cluster/export", r.APIClusterExport())
		svc.Post("/api/v1/cluster/import", r.APIClusterImport())
		svc.Get("/api/v1/roles/api", r.APIRoleConfigGet())
		svc.Post("/api/v1/roles/api", r.APIRoleConfigPut())
		svc.Post("/api/v1/tools/ping", r.APIToolPing())
		svc.Post("/api/v1/tools/traceroute", r.APIToolTraceroute())
		svc.Post("/api/v1/tools/portmap", r.APIToolPortmap())
	})
	return r
}

func (r *Role) Mux() *mux.Router {
	return r.m
}

func (r *Role) SessionStore() sessions.Store {
	return r.sessions
}

func (r *Role) Start(ctx context.Context, config []byte) error {
	r.cfg = r.decodeRoleConfig(ctx, config)

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
	sessDur := time.Hour * 24
	if d, err := time.ParseDuration(r.cfg.SessionDuration); err == nil {
		sessDur = d
	}
	sess.Options.MaxAge = int(sessDur.Seconds())
	if err != nil {
		return err
	}
	r.sessions = sess

	if r.cfg.OIDC != nil {
		go func() {
			_ = retry.Do(
				func() error {
					return r.auth.ConfigureOpenIDConnect(ctx, r.cfg.OIDC)
				},
				retry.DelayType(retry.BackOffDelay),
				retry.Attempts(50),
				retry.OnRetry(func(attempt uint, err error) {
					r.log.Warn("failed to setup OpenID Connect, retrying", zap.Uint("attempt", attempt), zap.Error(err))
				}),
			)
		}()
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
	r.oapi = web.NewService(openapi3.NewReflector())
	r.oapi.OpenAPISchema().SetTitle("gravity")
	r.oapi.OpenAPISchema().SetVersion(extconfig.Version)
	r.oapi.Method(http.MethodGet, "/api/v1/openapi.json", r.oapi.OpenAPICollector)

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
	r.httpServer = &http.Server{}
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
	r.socketServer = &http.Server{}
	r.socketServer.Handler = http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		ctx := context.WithValue(req.Context(), types.RequestSession, &sessions.Session{
			Values: map[interface{}]interface{}{
				types.SessionKeyUser: &types.User{
					Username: "gravity-socket",
					Permissions: []*types.Permission{
						{
							Path:    "/*",
							Methods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodHead, http.MethodDelete},
						},
					},
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
	if r.httpServer != nil {
		err := r.httpServer.Shutdown(r.ctx)
		if err != nil {
			r.log.Warn("failed to shutdown http server", zap.Error(err))
		}
	}
	if r.socketServer != nil {
		err := r.socketServer.Shutdown(r.ctx)
		if err != nil {
			r.log.Warn("failed to shutdown socket server", zap.Error(err))
		}
	}
	socketPath := path.Join(VAR_RUN, GRAVITY_SOCK)
	if extconfig.Get().Debug {
		socketPath = path.Join("./", GRAVITY_SOCK)
	}
	err := os.Remove(socketPath)
	if err != nil && !errors.Is(err, fs.ErrNotExist) {
		r.log.Warn("failed to remove socket", zap.Error(err))
	}
}
