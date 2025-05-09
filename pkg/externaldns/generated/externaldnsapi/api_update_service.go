// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

/*
 * External DNS Webhook Server
 *
 * Implements the external DNS webhook endpoints.
 *
 * API version: v0.15.0
 */

package externaldnsapi

import (
	"context"
	"errors"
	"net/http"
)

// UpdateAPIService is a service that implements the logic for the UpdateAPIServicer
// This service should implement the business logic for every endpoint for the UpdateAPI API.
// Include any external packages or services that will be required by this service.
type UpdateAPIService struct {
}

// NewUpdateAPIService creates a default api service
func NewUpdateAPIService() *UpdateAPIService {
	return &UpdateAPIService{}
}

// SetRecords - Applies the changes.
func (s *UpdateAPIService) SetRecords(ctx context.Context, changes Changes) (ImplResponse, error) {
	// TODO - update SetRecords with the required logic for this service method.
	// Add api_update_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(204, {}) or use other options such as http.Ok ...
	// return Response(204, nil),nil

	// TODO: Uncomment the next line to return response Response(500, {}) or use other options such as http.Ok ...
	// return Response(500, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("SetRecords method not implemented")
}

// AdjustRecords - Executes the AdjustEndpoints method.
func (s *UpdateAPIService) AdjustRecords(ctx context.Context, endpoint []Endpoint) (ImplResponse, error) {
	// TODO - update AdjustRecords with the required logic for this service method.
	// Add api_update_service.go to the .openapi-generator-ignore to avoid overwriting this service implementation when updating open api generation.

	// TODO: Uncomment the next line to return response Response(200, []Endpoint{}) or use other options such as http.Ok ...
	// return Response(200, []Endpoint{}), nil

	// TODO: Uncomment the next line to return response Response(500, {}) or use other options such as http.Ok ...
	// return Response(500, nil),nil

	return Response(http.StatusNotImplemented, nil), errors.New("AdjustRecords method not implemented")
}
