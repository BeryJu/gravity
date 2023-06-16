# \RolesEtcdApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**EtcdJoinMember**](RolesEtcdApi.md#EtcdJoinMember) | **Post** /api/v1/etcd/join | Etcd join



## EtcdJoinMember

> ApiAPIMemberJoinOutput EtcdJoinMember(ctx).ApiAPIMemberJoinInput(apiAPIMemberJoinInput).Execute()

Etcd join

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
    apiAPIMemberJoinInput := *openapiclient.NewApiAPIMemberJoinInput() // ApiAPIMemberJoinInput |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesEtcdApi.EtcdJoinMember(context.Background()).ApiAPIMemberJoinInput(apiAPIMemberJoinInput).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesEtcdApi.EtcdJoinMember``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `EtcdJoinMember`: ApiAPIMemberJoinOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesEtcdApi.EtcdJoinMember`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiEtcdJoinMemberRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **apiAPIMemberJoinInput** | [**ApiAPIMemberJoinInput**](ApiAPIMemberJoinInput.md) |  | 

### Return type

[**ApiAPIMemberJoinOutput**](ApiAPIMemberJoinOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

