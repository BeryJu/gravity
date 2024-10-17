# \RolesTsdbApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**TsdbGetMetrics**](RolesTsdbApi.md#TsdbGetMetrics) | **Get** /api/v1/tsdb/metrics | Retrieve Metrics
[**TsdbGetRoleConfig**](RolesTsdbApi.md#TsdbGetRoleConfig) | **Get** /api/v1/roles/tsdb | TSDB role config
[**TsdbPutRoleConfig**](RolesTsdbApi.md#TsdbPutRoleConfig) | **Post** /api/v1/roles/tsdb | TSDB role config



## TsdbGetMetrics

> TypesAPIMetricsGetOutput TsdbGetMetrics(ctx).Role(role).Category(category).Node(node).Since(since).Execute()

Retrieve Metrics

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    "time"
    openapiclient "beryju.io/gravity/api"
)

func main() {
    role := openapiclient.TypesAPIMetricsRole("system") // TypesAPIMetricsRole | 
    category := "category_example" // string |  (optional)
    node := "node_example" // string |  (optional)
    since := time.Now() // time.Time | Optionally set a start time for which to return datapoints after (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesTsdbApi.TsdbGetMetrics(context.Background()).Role(role).Category(category).Node(node).Since(since).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesTsdbApi.TsdbGetMetrics``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `TsdbGetMetrics`: TypesAPIMetricsGetOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesTsdbApi.TsdbGetMetrics`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiTsdbGetMetricsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **role** | [**TypesAPIMetricsRole**](TypesAPIMetricsRole.md) |  | 
 **category** | **string** |  | 
 **node** | **string** |  | 
 **since** | **time.Time** | Optionally set a start time for which to return datapoints after | 

### Return type

[**TypesAPIMetricsGetOutput**](TypesAPIMetricsGetOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## TsdbGetRoleConfig

> TsdbAPIRoleConfigOutput TsdbGetRoleConfig(ctx).Execute()

TSDB role config

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "beryju.io/gravity/api"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesTsdbApi.TsdbGetRoleConfig(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesTsdbApi.TsdbGetRoleConfig``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `TsdbGetRoleConfig`: TsdbAPIRoleConfigOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesTsdbApi.TsdbGetRoleConfig`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiTsdbGetRoleConfigRequest struct via the builder pattern


### Return type

[**TsdbAPIRoleConfigOutput**](TsdbAPIRoleConfigOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## TsdbPutRoleConfig

> TsdbPutRoleConfig(ctx).TsdbAPIRoleConfigInput(tsdbAPIRoleConfigInput).Execute()

TSDB role config

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "beryju.io/gravity/api"
)

func main() {
    tsdbAPIRoleConfigInput := *openapiclient.NewTsdbAPIRoleConfigInput(*openapiclient.NewTsdbRoleConfig()) // TsdbAPIRoleConfigInput |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.RolesTsdbApi.TsdbPutRoleConfig(context.Background()).TsdbAPIRoleConfigInput(tsdbAPIRoleConfigInput).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesTsdbApi.TsdbPutRoleConfig``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiTsdbPutRoleConfigRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **tsdbAPIRoleConfigInput** | [**TsdbAPIRoleConfigInput**](TsdbAPIRoleConfigInput.md) |  | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

