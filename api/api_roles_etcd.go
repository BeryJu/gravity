/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.27.1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
)

// RolesEtcdApiService RolesEtcdApi service
type RolesEtcdApiService service

type ApiEtcdJoinMemberRequest struct {
	ctx                   context.Context
	ApiService            *RolesEtcdApiService
	apiAPIMemberJoinInput *ApiAPIMemberJoinInput
}

func (r ApiEtcdJoinMemberRequest) ApiAPIMemberJoinInput(apiAPIMemberJoinInput ApiAPIMemberJoinInput) ApiEtcdJoinMemberRequest {
	r.apiAPIMemberJoinInput = &apiAPIMemberJoinInput
	return r
}

func (r ApiEtcdJoinMemberRequest) Execute() (*ApiAPIMemberJoinOutput, *http.Response, error) {
	return r.ApiService.EtcdJoinMemberExecute(r)
}

/*
EtcdJoinMember Etcd join

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@return ApiEtcdJoinMemberRequest
*/
func (a *RolesEtcdApiService) EtcdJoinMember(ctx context.Context) ApiEtcdJoinMemberRequest {
	return ApiEtcdJoinMemberRequest{
		ApiService: a,
		ctx:        ctx,
	}
}

// Execute executes the request
//
//	@return ApiAPIMemberJoinOutput
func (a *RolesEtcdApiService) EtcdJoinMemberExecute(r ApiEtcdJoinMemberRequest) (*ApiAPIMemberJoinOutput, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *ApiAPIMemberJoinOutput
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "RolesEtcdApiService.EtcdJoinMember")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/api/v1/etcd/join"

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.apiAPIMemberJoinInput
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
		if localVarHTTPResponse.StatusCode == 500 {
			var v RestErrResponse
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
			newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
			newErr.model = v
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}
