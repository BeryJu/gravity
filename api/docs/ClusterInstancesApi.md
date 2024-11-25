# \ClusterInstancesApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ClusterGetInstanceInfo**](ClusterInstancesApi.md#ClusterGetInstanceInfo) | **Get** /api/v1/cluster/instance | Instance
[**ClusterInstanceRoleRestart**](ClusterInstancesApi.md#ClusterInstanceRoleRestart) | **Post** /api/v1/cluster/roles/restart | Instance roles



## ClusterGetInstanceInfo

> InstanceAPIInstanceInfo ClusterGetInstanceInfo(ctx).Execute()

Instance

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
    resp, r, err := apiClient.ClusterInstancesApi.ClusterGetInstanceInfo(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ClusterInstancesApi.ClusterGetInstanceInfo``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ClusterGetInstanceInfo`: InstanceAPIInstanceInfo
    fmt.Fprintf(os.Stdout, "Response from `ClusterInstancesApi.ClusterGetInstanceInfo`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiClusterGetInstanceInfoRequest struct via the builder pattern


### Return type

[**InstanceAPIInstanceInfo**](InstanceAPIInstanceInfo.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ClusterInstanceRoleRestart

> ClusterInstanceRoleRestart(ctx).InstanceAPIRoleRestartInput(instanceAPIRoleRestartInput).Execute()

Instance roles

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
    instanceAPIRoleRestartInput := *openapiclient.NewInstanceAPIRoleRestartInput() // InstanceAPIRoleRestartInput |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.ClusterInstancesApi.ClusterInstanceRoleRestart(context.Background()).InstanceAPIRoleRestartInput(instanceAPIRoleRestartInput).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ClusterInstancesApi.ClusterInstanceRoleRestart``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiClusterInstanceRoleRestartRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **instanceAPIRoleRestartInput** | [**InstanceAPIRoleRestartInput**](InstanceAPIRoleRestartInput.md) |  | 

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

