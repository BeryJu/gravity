package monitoring

import (
	"context"
	"net/http"
	_ "unsafe"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

type Role struct {
	m   *mux.Router
	log *log.Entry
	i   roles.Instance
	ctx context.Context
}

//go:linkname blockyReg github.com/0xERR0R/blocky/metrics.reg
var blockyReg = prometheus.NewRegistry()

func New(instance roles.Instance) *Role {
	r := &Role{
		log: instance.Log(),
		i:   instance,
		m:   mux.NewRouter(),
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
		promhttp.InstrumentMetricHandler(blockyReg, promhttp.HandlerFor(blockyReg, promhttp.HandlerOpts{
			DisableCompression: true,
		})).ServeHTTP(w, r)
	})
	return r
}

func (r *Role) Start(ctx context.Context, config []byte) error {
	r.ctx = ctx
	cfg := r.decodeRoleConfig(config)
	listen := extconfig.Get().Listen(cfg.Port)
	r.log.WithField("listen", listen).Info("starting monitoring Server")
	return http.ListenAndServe(listen, r.m)
}

func (r *Role) Stop() {
}
