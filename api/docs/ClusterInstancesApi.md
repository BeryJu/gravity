# \ClusterInstancesApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ClusterGetInfo**](ClusterInstancesApi.md#ClusterGetInfo) | **Get** /api/v1/cluster/info | Instance
[**ClusterGetInstances**](ClusterInstancesApi.md#ClusterGetInstances) | **Get** /api/v1/cluster/instances | Instances
[**ClusterInstanceRoleRestart**](ClusterInstancesApi.md#ClusterInstanceRoleRestart) | **Post** /api/v1/cluster/roles/restart | Instance roles



## ClusterGetInfo

> InstanceAPIInstanceInfo ClusterGetInfo(ctx).Execute()

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
    resp, r, err := apiClient.ClusterInstancesApi.ClusterGetInfo(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ClusterInstancesApi.ClusterGetInfo``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ClusterGetInfo`: InstanceAPIInstanceInfo
    fmt.Fprintf(os.Stdout, "Response from `ClusterInstancesApi.ClusterGetInfo`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiClusterGetInfoRequest struct via the builder pattern


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


## ClusterGetInstances

> InstanceAPIInstancesOutput ClusterGetInstances(ctx).Execute()

Instances

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
    resp, r, err := apiClient.ClusterInstancesApi.ClusterGetInstances(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ClusterInstancesApi.ClusterGetInstances``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ClusterGetInstances`: InstanceAPIInstancesOutput
    fmt.Fprintf(os.Stdout, "Response from `ClusterInstancesApi.ClusterGetInstances`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiClusterGetInstancesRequest struct via the builder pattern


### Return type

[**InstanceAPIInstancesOutput**](InstanceAPIInstancesOutput.md)

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

