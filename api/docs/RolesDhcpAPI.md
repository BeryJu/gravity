# \RolesDhcpAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DhcpDeleteLeases**](RolesDhcpAPI.md#DhcpDeleteLeases) | **Delete** /api/v1/dhcp/scopes/leases | DHCP Leases
[**DhcpDeleteScopes**](RolesDhcpAPI.md#DhcpDeleteScopes) | **Delete** /api/v1/dhcp/scopes | DHCP Scopes
[**DhcpGetLeases**](RolesDhcpAPI.md#DhcpGetLeases) | **Get** /api/v1/dhcp/scopes/leases | DHCP Leases
[**DhcpGetRoleConfig**](RolesDhcpAPI.md#DhcpGetRoleConfig) | **Get** /api/v1/roles/dhcp | DHCP role config
[**DhcpGetScopes**](RolesDhcpAPI.md#DhcpGetScopes) | **Get** /api/v1/dhcp/scopes | DHCP Scopes
[**DhcpImportScopes**](RolesDhcpAPI.md#DhcpImportScopes) | **Post** /api/v1/dhcp/scopes/import | DHCP Scopes
[**DhcpPutLeases**](RolesDhcpAPI.md#DhcpPutLeases) | **Post** /api/v1/dhcp/scopes/leases | DHCP Leases
[**DhcpPutRoleConfig**](RolesDhcpAPI.md#DhcpPutRoleConfig) | **Post** /api/v1/roles/dhcp | DHCP role config
[**DhcpPutScopes**](RolesDhcpAPI.md#DhcpPutScopes) | **Post** /api/v1/dhcp/scopes | DHCP Scopes
[**DhcpWolLeases**](RolesDhcpAPI.md#DhcpWolLeases) | **Post** /api/v1/dhcp/scopes/leases/wol | DHCP Leases



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
	openapiclient "beryju.io/gravity/api"
)

func main() {
	identifier := "identifier_example" // string |  (optional)
	scope := "scope_example" // string |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.RolesDhcpAPI.DhcpDeleteLeases(context.Background()).Identifier(identifier).Scope(scope).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesDhcpAPI.DhcpDeleteLeases``: %v\n", err)
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
	openapiclient "beryju.io/gravity/api"
)

func main() {
	scope := "scope_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.RolesDhcpAPI.DhcpDeleteScopes(context.Background()).Scope(scope).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesDhcpAPI.DhcpDeleteScopes``: %v\n", err)
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

> DhcpAPILeasesGetOutput DhcpGetLeases(ctx).Scope(scope).Identifier(identifier).Execute()

DHCP Leases

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
	scope := "scope_example" // string |  (optional)
	identifier := "identifier_example" // string | Optional identifier of a lease to get (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.RolesDhcpAPI.DhcpGetLeases(context.Background()).Scope(scope).Identifier(identifier).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesDhcpAPI.DhcpGetLeases``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DhcpGetLeases`: DhcpAPILeasesGetOutput
	fmt.Fprintf(os.Stdout, "Response from `RolesDhcpAPI.DhcpGetLeases`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDhcpGetLeasesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **scope** | **string** |  | 
 **identifier** | **string** | Optional identifier of a lease to get | 

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
	openapiclient "beryju.io/gravity/api"
)

func main() {

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.RolesDhcpAPI.DhcpGetRoleConfig(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesDhcpAPI.DhcpGetRoleConfig``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DhcpGetRoleConfig`: DhcpAPIRoleConfigOutput
	fmt.Fprintf(os.Stdout, "Response from `RolesDhcpAPI.DhcpGetRoleConfig`: %v\n", resp)
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

> DhcpAPIScopesGetOutput DhcpGetScopes(ctx).Name(name).Execute()

DHCP Scopes

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
	name := "name_example" // string | Optionally get DHCP Scope by name (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.RolesDhcpAPI.DhcpGetScopes(context.Background()).Name(name).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesDhcpAPI.DhcpGetScopes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DhcpGetScopes`: DhcpAPIScopesGetOutput
	fmt.Fprintf(os.Stdout, "Response from `RolesDhcpAPI.DhcpGetScopes`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDhcpGetScopesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string** | Optionally get DHCP Scope by name | 

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


## DhcpImportScopes

> DhcpAPIScopesImportOutput DhcpImportScopes(ctx).Scope(scope).DhcpAPIScopesImportInput(dhcpAPIScopesImportInput).Execute()

DHCP Scopes

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
	scope := "scope_example" // string |  (optional)
	dhcpAPIScopesImportInput := *openapiclient.NewDhcpAPIScopesImportInput() // DhcpAPIScopesImportInput |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.RolesDhcpAPI.DhcpImportScopes(context.Background()).Scope(scope).DhcpAPIScopesImportInput(dhcpAPIScopesImportInput).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesDhcpAPI.DhcpImportScopes``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `DhcpImportScopes`: DhcpAPIScopesImportOutput
	fmt.Fprintf(os.Stdout, "Response from `RolesDhcpAPI.DhcpImportScopes`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDhcpImportScopesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **scope** | **string** |  | 
 **dhcpAPIScopesImportInput** | [**DhcpAPIScopesImportInput**](DhcpAPIScopesImportInput.md) |  | 

### Return type

[**DhcpAPIScopesImportOutput**](DhcpAPIScopesImportOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
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
	openapiclient "beryju.io/gravity/api"
)

func main() {
	identifier := "identifier_example" // string | 
	scope := "scope_example" // string | 
	dhcpAPILeasesPutInput := *openapiclient.NewDhcpAPILeasesPutInput("Address_example", "AddressLeaseTime_example", "Hostname_example") // DhcpAPILeasesPutInput |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.RolesDhcpAPI.DhcpPutLeases(context.Background()).Identifier(identifier).Scope(scope).DhcpAPILeasesPutInput(dhcpAPILeasesPutInput).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesDhcpAPI.DhcpPutLeases``: %v\n", err)
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
	openapiclient "beryju.io/gravity/api"
)

func main() {
	dhcpAPIRoleConfigInput := *openapiclient.NewDhcpAPIRoleConfigInput(*openapiclient.NewDhcpRoleConfig()) // DhcpAPIRoleConfigInput |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.RolesDhcpAPI.DhcpPutRoleConfig(context.Background()).DhcpAPIRoleConfigInput(dhcpAPIRoleConfigInput).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesDhcpAPI.DhcpPutRoleConfig``: %v\n", err)
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
	openapiclient "beryju.io/gravity/api"
)

func main() {
	scope := "scope_example" // string | 
	dhcpAPIScopesPutInput := *openapiclient.NewDhcpAPIScopesPutInput(false, "Hook_example", []openapiclient.TypesDHCPOption{*openapiclient.NewTypesDHCPOption()}, "SubnetCidr_example", int64(123)) // DhcpAPIScopesPutInput |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.RolesDhcpAPI.DhcpPutScopes(context.Background()).Scope(scope).DhcpAPIScopesPutInput(dhcpAPIScopesPutInput).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesDhcpAPI.DhcpPutScopes``: %v\n", err)
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
	openapiclient "beryju.io/gravity/api"
)

func main() {
	identifier := "identifier_example" // string | 
	scope := "scope_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.RolesDhcpAPI.DhcpWolLeases(context.Background()).Identifier(identifier).Scope(scope).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesDhcpAPI.DhcpWolLeases``: %v\n", err)
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

