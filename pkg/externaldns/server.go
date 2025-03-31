package externaldns

import (
	"context"
	"net/http"
	"sync"

	"beryju.io/gravity/api"
	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/externaldns/generated/externaldnsapi"
	apiRole "beryju.io/gravity/pkg/roles/api"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

type Server struct {
	apiRouter     *mux.Router
	metricsRouter *mux.Router
	api           *api.APIClient
	log           *zap.Logger
}

func New(api *api.APIClient) *Server {
	s := &Server{
		api: api,
		log: extconfig.Get().Logger().Named("external-dns"),
	}

	s.apiRouter = externaldnsapi.NewRouter(
		externaldnsapi.NewInitializationAPIController(s),
		externaldnsapi.NewListingAPIController(s),
		externaldnsapi.NewUpdateAPIController(s),
	)
	s.apiRouter.Use(apiRole.NewRecoverMiddleware(s.log))
	s.apiRouter.Use(apiRole.NewLoggingMiddleware(s.log, nil))

	s.metricsRouter = mux.NewRouter()
	s.metricsRouter.Path("/metrics").Handler(promhttp.Handler())
	s.metricsRouter.Use(apiRole.NewLoggingMiddleware(s.log, nil))
	return s
}

func (s *Server) Run() {
	// https://kubernetes-sigs.github.io/external-dns/v0.14.2/tutorials/webhook-provider/
	apiListen := "localhost:8888"
	metricsListen := "localhost:8080"

	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		s.log.Info("Serving external-dns API", zap.String("listen", apiListen))
		http.ListenAndServe(apiListen, s.apiRouter)
	}()
	go func() {
		defer wg.Done()
		s.log.Info("Serving metrics", zap.String("listen", metricsListen))
		http.ListenAndServe(metricsListen, s.metricsRouter)
	}()
	wg.Wait()
}

func (s *Server) Negotiate(ctx context.Context) (externaldnsapi.ImplResponse, error) {
	return externaldnsapi.Response(200, externaldnsapi.Filters{}), nil
}

func (s *Server) GetRecords(context.Context) (externaldnsapi.ImplResponse, error) {
	return externaldnsapi.Response(200, struct{}{}), nil
}

func (s *Server) SetRecords(context.Context, externaldnsapi.Changes) (externaldnsapi.ImplResponse, error) {
	return externaldnsapi.Response(200, struct{}{}), nil
}

func (s *Server) AdjustRecords(context.Context, []externaldnsapi.Endpoint) (externaldnsapi.ImplResponse, error) {
	return externaldnsapi.Response(200, struct{}{}), nil
}
