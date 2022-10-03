package debug

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/pprof"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/roles/debug/types"
	"github.com/getsentry/sentry-go"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

type Role struct {
	m      *mux.Router
	log    *log.Entry
	i      roles.Instance
	ctx    context.Context
	server *http.Server
}

func New(instance roles.Instance) *Role {
	mux := mux.NewRouter()
	r := &Role{
		log: instance.Log(),
		i:   instance,
		m:   mux,
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
	r.m.HandleFunc("/debug/sentry", func(w http.ResponseWriter, r *http.Request) {
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetTag("gravity.Testerror", "true")
			sentry.CaptureException(errors.New("debug test error"))
		})
	})
	return r
}

func (r *Role) Start(ctx context.Context, config []byte) error {
	r.ctx = ctx
	r.i.DispatchEvent(types.EventTopicDebugMuxSetup, roles.NewEvent(
		ctx,
		map[string]interface{}{
			"mux": r.m,
		},
	))
	listen := extconfig.Get().Listen(8010)
	if !extconfig.Get().Debug {
		return roles.ErrRoleNotConfigured
	}
	r.log.WithField("listen", listen).Info("starting debug Server")
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

func (r *Role) Stop() {
	if r.server != nil {
		r.server.Shutdown(r.ctx)
	}
}

func (r *Role) Index(w http.ResponseWriter, re *http.Request) {
	r.m.Walk(func(route *mux.Route, router *mux.Router, ancestors []*mux.Route) error {
		tpl, err := route.GetPathTemplate()
		if err != nil {
			return nil
		}
		w.Write([]byte(fmt.Sprintf("<a href='%[1]s'>%[1]s</a><br>", tpl)))
		return nil
	})
}
