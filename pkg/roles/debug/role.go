package debug

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"
	"runtime"

	"github.com/felixge/fgprof"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/roles/debug/types"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	"go.uber.org/zap"
)

type Role struct {
	m      *mux.Router
	log    *zap.Logger
	i      roles.Instance
	ctx    context.Context
	server *http.Server
}

func init() {
	roles.Register("debug", func(i roles.Instance) roles.Role {
		return New(i)
	})
}

func New(instance roles.Instance) *Role {
	mux := mux.NewRouter()
	r := &Role{
		log: instance.Log(),
		i:   instance,
		m:   mux,
		ctx: instance.Context(),
	}
	r.m.Use(api.NewRecoverMiddleware(r.log))
	r.m.Use(api.NewLoggingMiddleware(r.log, nil))
	r.m.HandleFunc("/", r.Index)
	r.m.HandleFunc("/debug/pprof/", pprof.Index)
	r.m.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	r.m.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.m.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.m.HandleFunc("/debug/pprof/trace", pprof.Trace)
	r.m.HandleFunc("/debug/pprof/{cmd}", pprof.Index)
	r.m.HandleFunc("/debug/fgprof", fgprof.Handler().ServeHTTP)
	r.m.HandleFunc("/debug/sentry", func(w http.ResponseWriter, r *http.Request) {
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetTag("gravity.Testerror", "true")
			sentry.CaptureException(errors.New("debug test error"))
		})
	})
	return r
}

func (r *Role) Mux() *mux.Router {
	return r.m
}

func (r *Role) Start(ctx context.Context, config []byte) error {
	r.i.DispatchEvent(types.EventTopicDebugMuxSetup, roles.NewEvent(
		ctx,
		map[string]interface{}{
			"mux": r.m,
		},
	))
	listen := extconfig.Get().Listen(8010)
	runtime.SetBlockProfileRate(5)

	r.log.Info("starting debug server", zap.String("listen", listen))
	r.server = &http.Server{
		Addr:    listen,
		Handler: r.m,
	}
	go func() {
		err := r.server.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			r.log.Warn("failed to listen", zap.Error(err))
		}
	}()
	return nil
}

func (r *Role) Stop() {
	if r.server != nil {
		err := r.server.Shutdown(r.ctx)
		if err != nil {
			r.log.Warn("failed to shutdown server", zap.Error(err))
		}
	}
}

func (r *Role) Index(w http.ResponseWriter, re *http.Request) {
	err := r.m.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, err := route.GetPathTemplate()
		if err != nil {
			return nil
		}
		_, err = fmt.Fprintf(w, "<a href='%[1]s'>%[1]s</a><br>", tpl)
		if err != nil {
			r.log.Warn("failed to write index overview link")
		}
		return nil
	})
	if err != nil {
		r.log.Warn("failed to walk routes for index", zap.Error(err))
	}
}
