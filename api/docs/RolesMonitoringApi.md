# \RolesMonitoringApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**MonitoringGetRoleConfig**](RolesMonitoringApi.md#MonitoringGetRoleConfig) | **Get** /api/v1/roles/monitoring | Monitoring role config
[**MonitoringPutRoleConfig**](RolesMonitoringApi.md#MonitoringPutRoleConfig) | **Post** /api/v1/roles/monitoring | Monitoring role config



## MonitoringGetRoleConfig

> MonitoringAPIRoleConfigOutput MonitoringGetRoleConfig(ctx).Execute()

Monitoring role config

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesMonitoringApi.MonitoringGetRoleConfig(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesMonitoringApi.MonitoringGetRoleConfig``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `MonitoringGetRoleConfig`: MonitoringAPIRoleConfigOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesMonitoringApi.MonitoringGetRoleConfig`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiMonitoringGetRoleConfigRequest struct via the builder pattern


### Return type

[**MonitoringAPIRoleConfigOutput**](MonitoringAPIRoleConfigOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## MonitoringPutRoleConfig

> MonitoringPutRoleConfig(ctx).MonitoringAPIRoleConfigInput(monitoringAPIRoleConfigInput).Execute()

Monitoring role config

### Example

```go
package main

import (
    "context"
    "fmt"
    "os"
    openapiclient "./openapi"
)

func main() {
    monitoringAPIRoleConfigInput := *openapiclient.NewMonitoringAPIRoleConfigInput(*openapiclient.NewMonitoringRoleConfig()) // MonitoringAPIRoleConfigInput |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesMonitoringApi.MonitoringPutRoleConfig(context.Background()).MonitoringAPIRoleConfigInput(monitoringAPIRoleConfigInput).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesMonitoringApi.MonitoringPutRoleConfig``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiMonitoringPutRoleConfigRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **monitoringAPIRoleConfigInput** | [**MonitoringAPIRoleConfigInput**](MonitoringAPIRoleConfigInput.md) |  | 

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

