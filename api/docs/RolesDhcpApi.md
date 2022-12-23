# \RolesDhcpApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DhcpDeleteLeases**](RolesDhcpApi.md#DhcpDeleteLeases) | **Delete** /api/v1/dhcp/scopes/leases | DHCP Leases
[**DhcpDeleteScopes**](RolesDhcpApi.md#DhcpDeleteScopes) | **Delete** /api/v1/dhcp/scopes | DHCP Scopes
[**DhcpGetLeases**](RolesDhcpApi.md#DhcpGetLeases) | **Get** /api/v1/dhcp/scopes/leases | DHCP Leases
[**DhcpGetRoleConfig**](RolesDhcpApi.md#DhcpGetRoleConfig) | **Get** /api/v1/roles/dhcp | DHCP role config
[**DhcpGetScopes**](RolesDhcpApi.md#DhcpGetScopes) | **Get** /api/v1/dhcp/scopes | DHCP Scopes
[**DhcpPutLeases**](RolesDhcpApi.md#DhcpPutLeases) | **Post** /api/v1/dhcp/scopes/leases | DHCP Leases
[**DhcpPutRoleConfig**](RolesDhcpApi.md#DhcpPutRoleConfig) | **Post** /api/v1/roles/dhcp | DHCP role config
[**DhcpPutScopes**](RolesDhcpApi.md#DhcpPutScopes) | **Post** /api/v1/dhcp/scopes | DHCP Scopes
[**DhcpWolLeases**](RolesDhcpApi.md#DhcpWolLeases) | **Post** /api/v1/dhcp/scopes/leases/wol | DHCP Leases



## DhcpDeleteLeases

> DhcpDeleteLeases(ctx).Identifier(identifier).Scope(scope).Execute()

DHCP Leases

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
    scope := "scope_example" // string |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesDhcpApi.DhcpDeleteLeases(context.Background()).Identifier(identifier).Scope(scope).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDhcpApi.DhcpDeleteLeases``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDhcpDeleteLeasesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **identifier** | **string** |  | 
 **scope** | **string** |  | 

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


## DhcpDeleteScopes

> DhcpDeleteScopes(ctx).Scope(scope).Execute()

DHCP Scopes

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
    scope := "scope_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesDhcpApi.DhcpDeleteScopes(context.Background()).Scope(scope).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDhcpApi.DhcpDeleteScopes``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDhcpDeleteScopesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **scope** | **string** |  | 

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


## DhcpGetLeases

> DhcpAPILeasesGetOutput DhcpGetLeases(ctx).Scope(scope).Execute()

DHCP Leases

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
    scope := "scope_example" // string |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesDhcpApi.DhcpGetLeases(context.Background()).Scope(scope).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDhcpApi.DhcpGetLeases``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DhcpGetLeases`: DhcpAPILeasesGetOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesDhcpApi.DhcpGetLeases`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDhcpGetLeasesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **scope** | **string** |  | 

### Return type

[**DhcpAPILeasesGetOutput**](DhcpAPILeasesGetOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DhcpGetRoleConfig

> DhcpAPIRoleConfigOutput DhcpGetRoleConfig(ctx).Execute()

DHCP role config

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
    resp, r, err := apiClient.RolesDhcpApi.DhcpGetRoleConfig(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDhcpApi.DhcpGetRoleConfig``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DhcpGetRoleConfig`: DhcpAPIRoleConfigOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesDhcpApi.DhcpGetRoleConfig`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiDhcpGetRoleConfigRequest struct via the builder pattern


### Return type

[**DhcpAPIRoleConfigOutput**](DhcpAPIRoleConfigOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DhcpGetScopes

> DhcpAPIScopesGetOutput DhcpGetScopes(ctx).Execute()

DHCP Scopes

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
    resp, r, err := apiClient.RolesDhcpApi.DhcpGetScopes(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDhcpApi.DhcpGetScopes``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DhcpGetScopes`: DhcpAPIScopesGetOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesDhcpApi.DhcpGetScopes`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiDhcpGetScopesRequest struct via the builder pattern


### Return type

[**DhcpAPIScopesGetOutput**](DhcpAPIScopesGetOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DhcpPutLeases

> DhcpPutLeases(ctx).Identifier(identifier).Scope(scope).DhcpAPILeasesPutInput(dhcpAPILeasesPutInput).Execute()

DHCP Leases

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
    scope := "scope_example" // string | 
    dhcpAPILeasesPutInput := *openapiclient.NewDhcpAPILeasesPutInput("Address_example", "AddressLeaseTime_example", "Hostname_example") // DhcpAPILeasesPutInput |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesDhcpApi.DhcpPutLeases(context.Background()).Identifier(identifier).Scope(scope).DhcpAPILeasesPutInput(dhcpAPILeasesPutInput).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDhcpApi.DhcpPutLeases``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDhcpPutLeasesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **identifier** | **string** |  | 
 **scope** | **string** |  | 
 **dhcpAPILeasesPutInput** | [**DhcpAPILeasesPutInput**](DhcpAPILeasesPutInput.md) |  | 

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


## DhcpPutRoleConfig

> DhcpPutRoleConfig(ctx).DhcpAPIRoleConfigInput(dhcpAPIRoleConfigInput).Execute()

DHCP role config

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
    dhcpAPIRoleConfigInput := *openapiclient.NewDhcpAPIRoleConfigInput(*openapiclient.NewDhcpRoleConfig()) // DhcpAPIRoleConfigInput |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesDhcpApi.DhcpPutRoleConfig(context.Background()).DhcpAPIRoleConfigInput(dhcpAPIRoleConfigInput).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDhcpApi.DhcpPutRoleConfig``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDhcpPutRoleConfigRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **dhcpAPIRoleConfigInput** | [**DhcpAPIRoleConfigInput**](DhcpAPIRoleConfigInput.md) |  | 

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


## DhcpPutScopes

> DhcpPutScopes(ctx).Scope(scope).DhcpAPIScopesPutInput(dhcpAPIScopesPutInput).Execute()

DHCP Scopes

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
    scope := "scope_example" // string | 
    dhcpAPIScopesPutInput := *openapiclient.NewDhcpAPIScopesPutInput(false, []openapiclient.TypesDHCPOption{*openapiclient.NewTypesDHCPOption()}, "SubnetCidr_example", int32(123)) // DhcpAPIScopesPutInput |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesDhcpApi.DhcpPutScopes(context.Background()).Scope(scope).DhcpAPIScopesPutInput(dhcpAPIScopesPutInput).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDhcpApi.DhcpPutScopes``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDhcpPutScopesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **scope** | **string** |  | 
 **dhcpAPIScopesPutInput** | [**DhcpAPIScopesPutInput**](DhcpAPIScopesPutInput.md) |  | 

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


## DhcpWolLeases

> DhcpWolLeases(ctx).Identifier(identifier).Scope(scope).Execute()

DHCP Leases

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
    scope := "scope_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesDhcpApi.DhcpWolLeases(context.Background()).Identifier(identifier).Scope(scope).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDhcpApi.DhcpWolLeases``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDhcpWolLeasesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **identifier** | **string** |  | 
 **scope** | **string** |  | 

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

