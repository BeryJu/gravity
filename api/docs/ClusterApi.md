# \ClusterApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ClusterGetClusterInfo**](ClusterApi.md#ClusterGetClusterInfo) | **Get** /api/v1/cluster | Cluster



## ClusterGetClusterInfo

> InstanceAPIClusterInfoOutput ClusterGetClusterInfo(ctx).Execute()

Cluster

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
    resp, r, err := apiClient.ClusterApi.ClusterGetClusterInfo(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `ClusterApi.ClusterGetClusterInfo``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ClusterGetClusterInfo`: InstanceAPIClusterInfoOutput
    fmt.Fprintf(os.Stdout, "Response from `ClusterApi.ClusterGetClusterInfo`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiClusterGetClusterInfoRequest struct via the builder pattern


### Return type

[**InstanceAPIClusterInfoOutput**](InstanceAPIClusterInfoOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

