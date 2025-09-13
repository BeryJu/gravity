package externaldns

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"sync"

	"log"

	"beryju.io/gravity/api"
	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/externaldns/generated/externaldnsapi"
	"beryju.io/gravity/pkg/roles/api/middleware"
	"beryju.io/gravity/pkg/roles/dns/types"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.uber.org/zap"
)

const (
	ProviderSpecificUid = "gravity_uid"
)

type Server struct {
	apiRouter     *mux.Router
	metricsRouter *mux.Router
	api           *api.APIClient
	log           *zap.Logger

	zones []api.DnsAPIZone
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
	config.UserAgent = fmt.Sprintf("gravity-external-dns/%s", extconfig.FullVersion())
	config.Debug = extconfig.Get().Debug
	apiClient := api.NewAPIClient(config)

	s := &Server{
		api:   apiClient,
		log:   extconfig.Get().Logger().Named("external-dns"),
		zones: []api.DnsAPIZone{},
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

func (s *Server) Run() error {
	s.log.Info("Pre-fetching zones...")
	err := s.prefetchZones(context.Background())
	if err != nil {
		s.log.Warn("failed to pre-fetch zones", zap.Error(err))
		return err
	}
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
	return nil
}

func (s *Server) errorResponse(err error) (externaldnsapi.ImplResponse, error) {
	s.log.Warn("Error", zap.Error(err))
	return externaldnsapi.Response(http.StatusInternalServerError, struct{}{}), err
}

func (s *Server) apiError(r *http.Response, err error) error {
	if err == nil {
		return nil
	}
	if r == nil {
		return fmt.Errorf("HTTP Error without http response: %v", err)
	}
	buff := &bytes.Buffer{}
	_, er := io.Copy(buff, r.Body)
	if er != nil {
		s.log.Warn("Gravity: failed to read response", zap.Error(er))
	}
	s.log.Warn("Gravity: error response", zap.String("body", buff.String()))
	return fmt.Errorf("HTTP Error '%s' during request '%s %s': \"%s\"", err.Error(), r.Request.Method, r.Request.URL.Path, buff.String())
}

func (s *Server) Negotiate(ctx context.Context) (externaldnsapi.ImplResponse, error) {
	return externaldnsapi.Response(http.StatusOK, externaldnsapi.Filters{
		Filters: Get().DomainFilter,
	}), nil
}

func (s *Server) GetRecords(ctx context.Context) (externaldnsapi.ImplResponse, error) {
	endpoints := []externaldnsapi.Endpoint{}
	for _, zone := range s.zones {
		records, hr, err := s.api.RolesDnsAPI.DnsGetRecords(ctx).Zone(zone.Name).Execute()
		if err != nil {
			return s.errorResponse(s.apiError(hr, err))
		}
		for _, record := range records.Records {
			endpoints = append(endpoints, externaldnsapi.Endpoint{
				DnsName:    strings.TrimSuffix(record.Fqdn, types.DNSSep),
				Targets:    []string{record.Data},
				RecordType: record.Type,
				RecordTTL:  int64(zone.DefaultTTL),
				ProviderSpecific: []externaldnsapi.ProviderSpecificProperty{
					{
						Name:  ProviderSpecificUid,
						Value: record.Uid,
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
		if _, err := s.endpointToWrite(ctx, endpoint); err != nil {
			return s.errorResponse(err)
		}
	}
	return externaldnsapi.Response(http.StatusOK, struct{}{}), nil
}

func (s *Server) AdjustRecords(ctx context.Context, endpoints []externaldnsapi.Endpoint) (externaldnsapi.ImplResponse, error) {
	for _, endpoint := range endpoints {
		if _, err := s.endpointToWrite(ctx, endpoint); err != nil {
			return s.errorResponse(err)
		}
	}
	return externaldnsapi.Response(http.StatusOK, []externaldnsapi.Endpoint{}), nil
}

func (s *Server) recordUID(uid string) string {
	return fmt.Sprintf("external-dns-%s", uid)
}

func (s *Server) recordUIDEndpointTarget(endpoint externaldnsapi.Endpoint, target string) string {
	return s.recordUID(fmt.Sprintf("%s-%s", endpoint.DnsName, target))
}

func (s *Server) endpointToDelete(ctx context.Context, endpoint externaldnsapi.Endpoint) error {
	zone, hostname := s.findZoneForRecord(endpoint.DnsName)
	if zone == nil {
		return fmt.Errorf("zone not found for record: %s", endpoint.DnsName)
	}
	for _, target := range endpoint.Targets {
		hr, err := s.api.RolesDnsAPI.
			DnsDeleteRecords(ctx).
			Zone(zone.Name).
			Hostname(hostname).
			Type_(endpoint.RecordType).
			Uid(s.recordUIDEndpointTarget(endpoint, target)).
			Execute()
		if err != nil {
			return s.apiError(hr, err)
		}
	}
	return nil
}

func (s *Server) endpointToWrite(ctx context.Context, endpoint externaldnsapi.Endpoint) (*externaldnsapi.Endpoint, error) {
	zone, hostname := s.findZoneForRecord(endpoint.DnsName)
	if zone == nil {
		return nil, fmt.Errorf("zone not found for record: %s", endpoint.DnsName)
	}
	for _, target := range endpoint.Targets {
		hr, err := s.api.RolesDnsAPI.
			DnsPutRecords(ctx).
			Zone(zone.Name).
			Hostname(hostname).
			Uid(s.recordUIDEndpointTarget(endpoint, target)).
			DnsAPIRecordsPutInput(api.DnsAPIRecordsPutInput{
				Type: endpoint.RecordType,
				Data: target,
			}).
			Execute()
		if err != nil {
			return nil, s.apiError(hr, err)
		}
	}
	return &endpoint, nil
}
