package monitoring

import (
	"context"
	"net/http"
	_ "unsafe"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api"
	apitypes "beryju.io/gravity/pkg/roles/api/types"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"github.com/swaggest/rest/web"
)

type Role struct {
	m      *mux.Router
	log    *log.Entry
	i      roles.Instance
	ctx    context.Context
	cfg    *RoleConfig
	server *http.Server
}

//go:linkname blockyReg github.com/0xERR0R/blocky/metrics.reg
var blockyReg = prometheus.NewRegistry()

func New(instance roles.Instance) *Role {
	mux := mux.NewRouter()
	r := &Role{
		log: instance.Log(),
		i:   instance,
		m:   mux,
	}
	r.m.Use(api.NewRecoverMiddleware(r.log))
	r.m.Use(api.NewLoggingMiddleware(r.log, nil))
	r.m.Path("/healthz/live").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
	})
	r.m.Path("/metrics").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		promhttp.InstrumentMetricHandler(prometheus.DefaultRegisterer, promhttp.HandlerFor(prometheus.DefaultGatherer, promhttp.HandlerOpts{
			DisableCompression: true,
		})).ServeHTTP(w, r)
		promhttp.HandlerFor(blockyReg, promhttp.HandlerOpts{
			DisableCompression: true,
		}).ServeHTTP(w, r)
	})
	r.i.AddEventListener(apitypes.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		svc.Get("/api/v1/roles/monitoring", r.apiHandlerRoleConfigGet())
		svc.Post("/api/v1/roles/monitoring", r.apiHandlerRoleConfigPut())
	})
	return r
}

func (r *Role) Start(ctx context.Context, config []byte) error {
	r.ctx = ctx
	r.cfg = r.decodeRoleConfig(config)
	listen := extconfig.Get().Listen(r.cfg.Port)
	r.log.WithField("listen", listen).Info("starting monitoring Server")
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
	r.server.Shutdown(r.ctx)
}
