package externaldns

import (
	"context"
	"net/http"

	"beryju.io/gravity/pkg/externaldns/generated/externaldnsapi"
)

type Server struct {
}

func (s *Server) Run() {
	r := externaldnsapi.NewRouter(
		externaldnsapi.NewInitializationAPIController(s),
		externaldnsapi.NewListingAPIController(s),
		externaldnsapi.NewUpdateAPIController(s),
	)

	http.ListenAndServe("localhost:8888", r)
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
