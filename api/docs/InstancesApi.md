# \InstancesApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**RootGetInfo**](InstancesApi.md#RootGetInfo) | **Get** /api/v1/info | Instances
[**RootGetInstances**](InstancesApi.md#RootGetInstances) | **Get** /api/v1/instances | Instances



## RootGetInfo

> InstanceAPIInstanceInfo RootGetInfo(ctx).Execute()

Instances

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
    resp, r, err := apiClient.InstancesApi.RootGetInfo(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `InstancesApi.RootGetInfo``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `RootGetInfo`: InstanceAPIInstanceInfo
    fmt.Fprintf(os.Stdout, "Response from `InstancesApi.RootGetInfo`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiRootGetInfoRequest struct via the builder pattern


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


## RootGetInstances

> InstanceAPIInstancesOutput RootGetInstances(ctx).Execute()

Instances

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
    resp, r, err := apiClient.InstancesApi.RootGetInstances(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `InstancesApi.RootGetInstances``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `RootGetInstances`: InstanceAPIInstancesOutput
    fmt.Fprintf(os.Stdout, "Response from `InstancesApi.RootGetInstances`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiRootGetInstancesRequest struct via the builder pattern


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

