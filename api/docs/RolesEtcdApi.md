# \RolesEtcdApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**EtcdGetMembers**](RolesEtcdApi.md#EtcdGetMembers) | **Get** /api/v1/etcd/members | Etcd members
[**EtcdJoinMember**](RolesEtcdApi.md#EtcdJoinMember) | **Post** /api/v1/etcd/join | Etcd join
[**EtcdRemoveMember**](RolesEtcdApi.md#EtcdRemoveMember) | **Delete** /api/v1/etcd/remove | Etcd remove



## EtcdGetMembers

> EtcdAPIMembersOutput EtcdGetMembers(ctx).Execute()

Etcd members

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
    resp, r, err := apiClient.RolesEtcdApi.EtcdGetMembers(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesEtcdApi.EtcdGetMembers``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `EtcdGetMembers`: EtcdAPIMembersOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesEtcdApi.EtcdGetMembers`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiEtcdGetMembersRequest struct via the builder pattern


### Return type

[**EtcdAPIMembersOutput**](EtcdAPIMembersOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## EtcdJoinMember

> EtcdAPIMemberJoinOutput EtcdJoinMember(ctx).EtcdAPIMemberJoinInput(etcdAPIMemberJoinInput).Execute()

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
    etcdAPIMemberJoinInput := *openapiclient.NewEtcdAPIMemberJoinInput() // EtcdAPIMemberJoinInput |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesEtcdApi.EtcdJoinMember(context.Background()).EtcdAPIMemberJoinInput(etcdAPIMemberJoinInput).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesEtcdApi.EtcdJoinMember``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `EtcdJoinMember`: EtcdAPIMemberJoinOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesEtcdApi.EtcdJoinMember`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiEtcdJoinMemberRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **etcdAPIMemberJoinInput** | [**EtcdAPIMemberJoinInput**](EtcdAPIMemberJoinInput.md) |  | 

### Return type

[**EtcdAPIMemberJoinOutput**](EtcdAPIMemberJoinOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## EtcdRemoveMember

> EtcdRemoveMember(ctx).PeerID(peerID).Execute()

Etcd remove

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
    peerID := int32(56) // int32 | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.RolesEtcdApi.EtcdRemoveMember(context.Background()).PeerID(peerID).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesEtcdApi.EtcdRemoveMember``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiEtcdRemoveMemberRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **peerID** | **int32** |  | 

### Return type

 (empty response body)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

