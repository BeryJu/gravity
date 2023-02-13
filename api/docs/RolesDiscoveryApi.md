# \RolesDiscoveryApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DiscoveryApplyDevice**](RolesDiscoveryApi.md#DiscoveryApplyDevice) | **Post** /api/v1/discovery/devices/apply | Apply Discovered devices
[**DiscoveryDeleteDevices**](RolesDiscoveryApi.md#DiscoveryDeleteDevices) | **Delete** /api/v1/discovery/devices/delete | Discovery devices
[**DiscoveryDeleteSubnets**](RolesDiscoveryApi.md#DiscoveryDeleteSubnets) | **Delete** /api/v1/discovery/subnets | Discovery Subnets
[**DiscoveryGetDevices**](RolesDiscoveryApi.md#DiscoveryGetDevices) | **Get** /api/v1/discovery/devices | Discovery devices
[**DiscoveryGetRoleConfig**](RolesDiscoveryApi.md#DiscoveryGetRoleConfig) | **Get** /api/v1/roles/discovery | Discovery role config
[**DiscoveryGetSubnets**](RolesDiscoveryApi.md#DiscoveryGetSubnets) | **Get** /api/v1/discovery/subnets | Discovery subnets
[**DiscoveryPutRoleConfig**](RolesDiscoveryApi.md#DiscoveryPutRoleConfig) | **Post** /api/v1/roles/discovery | Discovery role config
[**DiscoveryPutSubnets**](RolesDiscoveryApi.md#DiscoveryPutSubnets) | **Post** /api/v1/discovery/subnets | Discovery Subnets
[**DiscoverySubnetStart**](RolesDiscoveryApi.md#DiscoverySubnetStart) | **Post** /api/v1/discovery/subnets/start | Discovery Subnets



## DiscoveryApplyDevice

> DiscoveryApplyDevice(ctx).Identifier(identifier).DiscoveryAPIDevicesApplyInput(discoveryAPIDevicesApplyInput).Execute()

Apply Discovered devices

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
    identifier := "identifier_example" // string | 
    discoveryAPIDevicesApplyInput := *openapiclient.NewDiscoveryAPIDevicesApplyInput("DhcpScope_example", "DnsZone_example", "To_example") // DiscoveryAPIDevicesApplyInput |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesDiscoveryApi.DiscoveryApplyDevice(context.Background()).Identifier(identifier).DiscoveryAPIDevicesApplyInput(discoveryAPIDevicesApplyInput).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDiscoveryApi.DiscoveryApplyDevice``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDiscoveryApplyDeviceRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **identifier** | **string** |  | 
 **discoveryAPIDevicesApplyInput** | [**DiscoveryAPIDevicesApplyInput**](DiscoveryAPIDevicesApplyInput.md) |  | 

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


## DiscoveryDeleteDevices

> DiscoveryDeleteDevices(ctx).Identifier(identifier).Execute()

Discovery devices

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
    identifier := "identifier_example" // string |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesDiscoveryApi.DiscoveryDeleteDevices(context.Background()).Identifier(identifier).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDiscoveryApi.DiscoveryDeleteDevices``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDiscoveryDeleteDevicesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **identifier** | **string** |  | 

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


## DiscoveryDeleteSubnets

> DiscoveryDeleteSubnets(ctx).Identifier(identifier).Execute()

Discovery Subnets

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
    identifier := "identifier_example" // string |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesDiscoveryApi.DiscoveryDeleteSubnets(context.Background()).Identifier(identifier).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDiscoveryApi.DiscoveryDeleteSubnets``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDiscoveryDeleteSubnetsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **identifier** | **string** |  | 

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


## DiscoveryGetDevices

> DiscoveryAPIDevicesGetOutput DiscoveryGetDevices(ctx).Identifier(identifier).Execute()

Discovery devices

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
    identifier := "identifier_example" // string | Optionally get device by identifier (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesDiscoveryApi.DiscoveryGetDevices(context.Background()).Identifier(identifier).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDiscoveryApi.DiscoveryGetDevices``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DiscoveryGetDevices`: DiscoveryAPIDevicesGetOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesDiscoveryApi.DiscoveryGetDevices`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDiscoveryGetDevicesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **identifier** | **string** | Optionally get device by identifier | 

### Return type

[**DiscoveryAPIDevicesGetOutput**](DiscoveryAPIDevicesGetOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DiscoveryGetRoleConfig

> DiscoveryAPIRoleConfigOutput DiscoveryGetRoleConfig(ctx).Execute()

Discovery role config

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
    resp, r, err := apiClient.RolesDiscoveryApi.DiscoveryGetRoleConfig(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDiscoveryApi.DiscoveryGetRoleConfig``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DiscoveryGetRoleConfig`: DiscoveryAPIRoleConfigOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesDiscoveryApi.DiscoveryGetRoleConfig`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiDiscoveryGetRoleConfigRequest struct via the builder pattern


### Return type

[**DiscoveryAPIRoleConfigOutput**](DiscoveryAPIRoleConfigOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DiscoveryGetSubnets

> DiscoveryAPISubnetsGetOutput DiscoveryGetSubnets(ctx).Name(name).Execute()

Discovery subnets

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
    name := "name_example" // string | Optionally get Subnet by name (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesDiscoveryApi.DiscoveryGetSubnets(context.Background()).Name(name).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDiscoveryApi.DiscoveryGetSubnets``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DiscoveryGetSubnets`: DiscoveryAPISubnetsGetOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesDiscoveryApi.DiscoveryGetSubnets`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDiscoveryGetSubnetsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string** | Optionally get Subnet by name | 

### Return type

[**DiscoveryAPISubnetsGetOutput**](DiscoveryAPISubnetsGetOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DiscoveryPutRoleConfig

> DiscoveryPutRoleConfig(ctx).DiscoveryAPIRoleConfigInput(discoveryAPIRoleConfigInput).Execute()

Discovery role config

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
    discoveryAPIRoleConfigInput := *openapiclient.NewDiscoveryAPIRoleConfigInput(*openapiclient.NewDiscoveryRoleConfig()) // DiscoveryAPIRoleConfigInput |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesDiscoveryApi.DiscoveryPutRoleConfig(context.Background()).DiscoveryAPIRoleConfigInput(discoveryAPIRoleConfigInput).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDiscoveryApi.DiscoveryPutRoleConfig``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDiscoveryPutRoleConfigRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **discoveryAPIRoleConfigInput** | [**DiscoveryAPIRoleConfigInput**](DiscoveryAPIRoleConfigInput.md) |  | 

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


## DiscoveryPutSubnets

> DiscoveryPutSubnets(ctx).Identifier(identifier).DiscoveryAPISubnetsPutInput(discoveryAPISubnetsPutInput).Execute()

Discovery Subnets

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
    identifier := "identifier_example" // string | 
    discoveryAPISubnetsPutInput := *openapiclient.NewDiscoveryAPISubnetsPutInput(int32(123), "DnsResolver_example", "SubnetCidr_example") // DiscoveryAPISubnetsPutInput |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesDiscoveryApi.DiscoveryPutSubnets(context.Background()).Identifier(identifier).DiscoveryAPISubnetsPutInput(discoveryAPISubnetsPutInput).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDiscoveryApi.DiscoveryPutSubnets``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDiscoveryPutSubnetsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **identifier** | **string** |  | 
 **discoveryAPISubnetsPutInput** | [**DiscoveryAPISubnetsPutInput**](DiscoveryAPISubnetsPutInput.md) |  | 

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


## DiscoverySubnetStart

> DiscoverySubnetStart(ctx).Identifier(identifier).Execute()

Discovery Subnets

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
    identifier := "identifier_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesDiscoveryApi.DiscoverySubnetStart(context.Background()).Identifier(identifier).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDiscoveryApi.DiscoverySubnetStart``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDiscoverySubnetStartRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **identifier** | **string** |  | 

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

