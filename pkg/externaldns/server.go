package externaldns

import (
	"context"
	"net/http"
	"sync"

	"beryju.io/gravity/api"
	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/externaldns/generated/externaldnsapi"
	"beryju.io/gravity/pkg/roles/api/middleware"
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
	s.apiRouter.Use(middleware.NewRecoverMiddleware(s.log))
	s.apiRouter.Use(middleware.NewLoggingMiddleware(s.log, nil))
	// This endpoint is "required" but not defined in the API specs
	s.apiRouter.Path("/healthz").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	})

	s.metricsRouter = mux.NewRouter()
	s.metricsRouter.Path("/metrics").Handler(promhttp.Handler())
	s.metricsRouter.Use(middleware.NewLoggingMiddleware(s.log, nil))
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

func (s *Server) errorResponse(err error) (externaldnsapi.ImplResponse, error) {
	return externaldnsapi.Response(500, struct{}{}), err
}

func (s *Server) Negotiate(ctx context.Context) (externaldnsapi.ImplResponse, error) {
	return externaldnsapi.Response(200, externaldnsapi.Filters{}), nil
}

func (s *Server) GetRecords(ctx context.Context) (externaldnsapi.ImplResponse, error) {
	zones, _, err := s.api.RolesDnsApi.DnsGetZones(ctx).Execute()
	if err != nil {
		return s.errorResponse(err)
	}
	endpoints := []externaldnsapi.Endpoint{}
	// TODO: Pagination
	for _, zone := range zones.Zones {
		records, _, err := s.api.RolesDnsApi.DnsGetRecords(ctx).Zone(zone.Name).Execute()
		if err != nil {
			return s.errorResponse(err)
		}
		for _, record := range records.Records {
			endpoints = append(endpoints, externaldnsapi.Endpoint{
				DnsName:    record.Hostname,
				Targets:    []string{record.Data},
				RecordType: record.Type,
				ProviderSpecific: []externaldnsapi.ProviderSpecificProperty{
					{
						Name:  "gravity_uid",
						Value: record.Uid,
					},
				},
			})
		}
	}
	return externaldnsapi.Response(200, endpoints), nil
}

func (s *Server) SetRecords(ctx context.Context, changes externaldnsapi.Changes) (externaldnsapi.ImplResponse, error) {
	return externaldnsapi.Response(200, struct{}{}), nil
}

func (s *Server) AdjustRecords(ctx context.Context, endpoints []externaldnsapi.Endpoint) (externaldnsapi.ImplResponse, error) {
	return externaldnsapi.Response(200, struct{}{}), nil
}
