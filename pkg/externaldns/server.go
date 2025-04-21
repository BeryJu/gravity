package externaldns

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"

	"log"

	"beryju.io/gravity/api"
	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/externaldns/generated/externaldnsapi"
	"beryju.io/gravity/pkg/roles/api/middleware"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

const (
	ProviderSpecificUid  = "gravity_uid"
	ProviderSpecificZone = "gravity_zone"
)

type Server struct {
	apiRouter     *mux.Router
	metricsRouter *mux.Router
	api           *api.APIClient
	log           *zap.Logger
}

func New() (*Server, error) {
	url, err := url.Parse(Get().Gravity.URL)
	if err != nil {
		return nil, err
	}

	config := api.NewConfiguration()
	config.Host = url.Host
	config.Scheme = url.Scheme
	if Get().Gravity.Token != "" {
		config.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", Get().Gravity.Token))
	}
	config.UserAgent = fmt.Sprintf("gravity-cli/%s", extconfig.FullVersion())
	apiClient := api.NewAPIClient(config)

	s := &Server{
		api: apiClient,
		log: extconfig.Get().Logger().Named("external-dns"),
	}

	// Discard go's inbuilt log as its used by the generated code and can't be disabled...
	log.SetOutput(io.Discard)
	s.apiRouter = externaldnsapi.NewRouter(
		externaldnsapi.NewInitializationAPIController(s),
		externaldnsapi.NewListingAPIController(s),
		externaldnsapi.NewUpdateAPIController(s),
	)
	s.apiRouter.Use(middleware.NewRecoverMiddleware(s.log))
	s.apiRouter.Use(middleware.NewLoggingMiddleware(s.log, nil))

	s.metricsRouter = mux.NewRouter()
	s.metricsRouter.Path("/metrics").Handler(promhttp.Handler())
	s.metricsRouter.Path("/healthz").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		_, _ = w.Write([]byte("ok"))
	})
	s.metricsRouter.Use(middleware.NewLoggingMiddleware(s.log, nil))
	return s, nil
}

func (s *Server) Run() {
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		s.log.Info("Serving external-dns API", zap.String("listen", Get().Listen.API))
		err := http.ListenAndServe(Get().Listen.API, s.apiRouter)
		if err != nil {
			s.log.Warn("Failed to listen", zap.Error(err))
		}
	}()
	go func() {
		defer wg.Done()
		s.log.Info("Serving metrics", zap.String("listen", Get().Listen.Metrics))
		err := http.ListenAndServe(Get().Listen.Metrics, s.metricsRouter)
		if err != nil {
			s.log.Warn("Failed to listen", zap.Error(err))
		}
	}()
	wg.Wait()
}

func (s *Server) errorResponse(err error) (externaldnsapi.ImplResponse, error) {
	s.log.Warn("Error", zap.Error(err))
	return externaldnsapi.Response(http.StatusInternalServerError, struct{}{}), err
}

func (s *Server) Negotiate(ctx context.Context) (externaldnsapi.ImplResponse, error) {
	return externaldnsapi.Response(http.StatusOK, externaldnsapi.Filters{
		Filters: Get().DomainFilter,
	}), nil
}

func (s *Server) GetRecords(ctx context.Context) (externaldnsapi.ImplResponse, error) {
	zones, _, err := s.api.RolesDnsApi.DnsGetZones(ctx).Execute()
	if err != nil {
		return s.errorResponse(err)
	}
	endpoints := []externaldnsapi.Endpoint{}
	// TODO: Pagination
	for _, zone := range zones.Zones {
		if !slices.Contains(Get().DomainFilter, zone.Name) {
			continue
		}
		records, _, err := s.api.RolesDnsApi.DnsGetRecords(ctx).Zone(zone.Name).Execute()
		if err != nil {
			return s.errorResponse(err)
		}
		for _, record := range records.Records {
			endpoints = append(endpoints, externaldnsapi.Endpoint{
				DnsName:    record.Fqdn,
				Targets:    []string{record.Data},
				RecordType: record.Type,
				RecordTTL:  int64(zone.DefaultTTL),
				ProviderSpecific: []externaldnsapi.ProviderSpecificProperty{
					{
						Name:  ProviderSpecificUid,
						Value: record.Uid,
					},
					{
						Name:  ProviderSpecificZone,
						Value: zone.Name,
					},
				},
			})
		}
	}
	return externaldnsapi.Response(http.StatusOK, endpoints), nil
}

func (s *Server) SetRecords(ctx context.Context, changes externaldnsapi.Changes) (externaldnsapi.ImplResponse, error) {
	for _, endpoint := range changes.Delete {
		if err := s.endpointToDelete(ctx, endpoint); err != nil {
			return s.errorResponse(err)
		}
	}
	for _, endpoint := range append(changes.Create, changes.UpdateNew...) {
		if err := s.endpointToWrite(ctx, endpoint); err != nil {
			return s.errorResponse(err)
		}
	}
	return externaldnsapi.Response(http.StatusOK, struct{}{}), nil
}

func (s *Server) AdjustRecords(ctx context.Context, endpoints []externaldnsapi.Endpoint) (externaldnsapi.ImplResponse, error) {
	for _, endpoint := range endpoints {
		if err := s.endpointToWrite(ctx, endpoint); err != nil {
			return s.errorResponse(err)
		}
	}
	return externaldnsapi.Response(http.StatusOK, []externaldnsapi.Endpoint{}), nil
}

func (s *Server) endpointToDelete(ctx context.Context, endpoint externaldnsapi.Endpoint) error {
	uid := ""
	zone := ""
	for _, pv := range endpoint.ProviderSpecific {
		switch pv.Name {
		case ProviderSpecificUid:
			uid = pv.Value
		case ProviderSpecificZone:
			zone = pv.Value
		}
	}
	hostname := strings.TrimSuffix(strings.Replace(endpoint.DnsName, zone, "", 1), ".")
	_, err := s.api.RolesDnsApi.DnsDeleteRecords(ctx).Zone(zone).Hostname(hostname).Type_(endpoint.RecordType).Uid(uid).Execute()
	return err
}

func (s *Server) endpointToWrite(ctx context.Context, endpoint externaldnsapi.Endpoint) error {
	uid := ""
	zone := ""
	for _, pv := range endpoint.ProviderSpecific {
		switch pv.Name {
		case ProviderSpecificUid:
			uid = pv.Value
		case ProviderSpecificZone:
			zone = pv.Value
		}
	}
	hostname := strings.TrimSuffix(strings.Replace(endpoint.DnsName, zone, "", 1), ".")
	_, err := s.api.RolesDnsApi.DnsPutRecords(ctx).Zone(zone).Hostname(hostname).Uid(uid).DnsAPIRecordsPutInput(api.DnsAPIRecordsPutInput{
		Type: endpoint.RecordType,
		Data: endpoint.Targets[0],
	}).Execute()
	return err
}
