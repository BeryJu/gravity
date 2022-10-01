package debug

import (
	"context"
	"errors"
	"net/http"
	"net/http/pprof"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api"
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
	r.m.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
	r.m.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
	r.m.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
	r.m.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
	r.m.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
	r.m.Handle("/debug/pprof/{cmd}", http.HandlerFunc(pprof.Index))
	r.m.HandleFunc("/debug/sentry", func(w http.ResponseWriter, r *http.Request) {
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetTag("gravity.test_error", "true")
			sentry.CaptureException(errors.New("debug test error"))
		})
	})
	return r
}

func (r *Role) Start(ctx context.Context, config []byte) error {
	r.ctx = ctx
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
