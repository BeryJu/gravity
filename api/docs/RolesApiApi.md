# \RolesApiAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApiAuthConfig**](RolesApiAPI.md#ApiAuthConfig) | **Get** /api/v1/auth/config | API Users
[**ApiDeleteTokens**](RolesApiAPI.md#ApiDeleteTokens) | **Delete** /api/v1/auth/tokens | Tokens
[**ApiDeleteUsers**](RolesApiAPI.md#ApiDeleteUsers) | **Delete** /api/v1/auth/users | API Users
[**ApiExport**](RolesApiAPI.md#ApiExport) | **Post** /api/v1/cluster/export | Export Cluster
[**ApiGetLogMessages**](RolesApiAPI.md#ApiGetLogMessages) | **Get** /api/v1/cluster/node/logs | Log messages
[**ApiGetRoleConfig**](RolesApiAPI.md#ApiGetRoleConfig) | **Get** /api/v1/roles/api | API role config
[**ApiGetTokens**](RolesApiAPI.md#ApiGetTokens) | **Get** /api/v1/auth/tokens | Tokens
[**ApiGetUsers**](RolesApiAPI.md#ApiGetUsers) | **Get** /api/v1/auth/users | API Users
[**ApiImport**](RolesApiAPI.md#ApiImport) | **Post** /api/v1/cluster/import | Import Cluster
[**ApiLoginUser**](RolesApiAPI.md#ApiLoginUser) | **Post** /api/v1/auth/login | API Users
[**ApiPutRoleConfig**](RolesApiAPI.md#ApiPutRoleConfig) | **Post** /api/v1/roles/api | API role config
[**ApiPutTokens**](RolesApiAPI.md#ApiPutTokens) | **Post** /api/v1/auth/tokens | Tokens
[**ApiPutUsers**](RolesApiAPI.md#ApiPutUsers) | **Post** /api/v1/auth/users | API Users
[**ApiUsersMe**](RolesApiAPI.md#ApiUsersMe) | **Get** /api/v1/auth/me | API Users
[**ToolsPing**](RolesApiAPI.md#ToolsPing) | **Post** /api/v1/tools/ping | Ping tool
[**ToolsPortmap**](RolesApiAPI.md#ToolsPortmap) | **Post** /api/v1/tools/portmap | Portmap tool
[**ToolsTraceroute**](RolesApiAPI.md#ToolsTraceroute) | **Post** /api/v1/tools/traceroute | Traceroute tool



## ApiAuthConfig

> AuthAPIConfigOutput ApiAuthConfig(ctx).Execute()

API Users

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
	resp, r, err := apiClient.RolesApiAPI.ApiAuthConfig(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesApiAPI.ApiAuthConfig``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiAuthConfig`: AuthAPIConfigOutput
	fmt.Fprintf(os.Stdout, "Response from `RolesApiAPI.ApiAuthConfig`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiApiAuthConfigRequest struct via the builder pattern


### Return type

[**AuthAPIConfigOutput**](AuthAPIConfigOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiDeleteTokens

> ApiDeleteTokens(ctx).Key(key).Execute()

Tokens

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
	key := "key_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.RolesApiAPI.ApiDeleteTokens(context.Background()).Key(key).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesApiAPI.ApiDeleteTokens``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiDeleteTokensRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **key** | **string** |  | 

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


## ApiDeleteUsers

> ApiDeleteUsers(ctx).Username(username).Execute()

API Users

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
	username := "username_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.RolesApiAPI.ApiDeleteUsers(context.Background()).Username(username).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesApiAPI.ApiDeleteUsers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiDeleteUsersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **username** | **string** |  | 

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


## ApiExport

> ApiAPIExportOutput ApiExport(ctx).ApiAPIExportInput(apiAPIExportInput).Execute()

Export Cluster

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
	apiAPIExportInput := *openapiclient.NewApiAPIExportInput() // ApiAPIExportInput |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.RolesApiAPI.ApiExport(context.Background()).ApiAPIExportInput(apiAPIExportInput).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesApiAPI.ApiExport``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiExport`: ApiAPIExportOutput
	fmt.Fprintf(os.Stdout, "Response from `RolesApiAPI.ApiExport`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiExportRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **apiAPIExportInput** | [**ApiAPIExportInput**](ApiAPIExportInput.md) |  | 

### Return type

[**ApiAPIExportOutput**](ApiAPIExportOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiGetLogMessages

> ApiAPILogMessages ApiGetLogMessages(ctx).Execute()

Log messages

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
	resp, r, err := apiClient.RolesApiAPI.ApiGetLogMessages(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesApiAPI.ApiGetLogMessages``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiGetLogMessages`: ApiAPILogMessages
	fmt.Fprintf(os.Stdout, "Response from `RolesApiAPI.ApiGetLogMessages`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiApiGetLogMessagesRequest struct via the builder pattern


### Return type

[**ApiAPILogMessages**](ApiAPILogMessages.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiGetRoleConfig

> ApiAPIRoleConfigOutput ApiGetRoleConfig(ctx).Execute()

API role config

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
	resp, r, err := apiClient.RolesApiAPI.ApiGetRoleConfig(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesApiAPI.ApiGetRoleConfig``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiGetRoleConfig`: ApiAPIRoleConfigOutput
	fmt.Fprintf(os.Stdout, "Response from `RolesApiAPI.ApiGetRoleConfig`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiApiGetRoleConfigRequest struct via the builder pattern


### Return type

[**ApiAPIRoleConfigOutput**](ApiAPIRoleConfigOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiGetTokens

> AuthAPITokensGetOutput ApiGetTokens(ctx).Execute()

Tokens

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
	resp, r, err := apiClient.RolesApiAPI.ApiGetTokens(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesApiAPI.ApiGetTokens``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiGetTokens`: AuthAPITokensGetOutput
	fmt.Fprintf(os.Stdout, "Response from `RolesApiAPI.ApiGetTokens`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiApiGetTokensRequest struct via the builder pattern


### Return type

[**AuthAPITokensGetOutput**](AuthAPITokensGetOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiGetUsers

> AuthAPIUsersGetOutput ApiGetUsers(ctx).Username(username).Execute()

API Users

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
	username := "username_example" // string | Optional username of a user to get (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.RolesApiAPI.ApiGetUsers(context.Background()).Username(username).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesApiAPI.ApiGetUsers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiGetUsers`: AuthAPIUsersGetOutput
	fmt.Fprintf(os.Stdout, "Response from `RolesApiAPI.ApiGetUsers`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiGetUsersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **username** | **string** | Optional username of a user to get | 

### Return type

[**AuthAPIUsersGetOutput**](AuthAPIUsersGetOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiImport

> ApiImport(ctx).ApiAPIImportInput(apiAPIImportInput).Execute()

Import Cluster

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
	apiAPIImportInput := *openapiclient.NewApiAPIImportInput() // ApiAPIImportInput |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.RolesApiAPI.ApiImport(context.Background()).ApiAPIImportInput(apiAPIImportInput).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesApiAPI.ApiImport``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiImportRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **apiAPIImportInput** | [**ApiAPIImportInput**](ApiAPIImportInput.md) |  | 

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


## ApiLoginUser

> AuthAPILoginOutput ApiLoginUser(ctx).AuthAPILoginInput(authAPILoginInput).Execute()

API Users

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
	authAPILoginInput := *openapiclient.NewAuthAPILoginInput() // AuthAPILoginInput |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.RolesApiAPI.ApiLoginUser(context.Background()).AuthAPILoginInput(authAPILoginInput).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesApiAPI.ApiLoginUser``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiLoginUser`: AuthAPILoginOutput
	fmt.Fprintf(os.Stdout, "Response from `RolesApiAPI.ApiLoginUser`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiLoginUserRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **authAPILoginInput** | [**AuthAPILoginInput**](AuthAPILoginInput.md) |  | 

### Return type

[**AuthAPILoginOutput**](AuthAPILoginOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiPutRoleConfig

> ApiPutRoleConfig(ctx).ApiAPIRoleConfigInput(apiAPIRoleConfigInput).Execute()

API role config

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
	apiAPIRoleConfigInput := *openapiclient.NewApiAPIRoleConfigInput(*openapiclient.NewApiRoleConfig()) // ApiAPIRoleConfigInput |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.RolesApiAPI.ApiPutRoleConfig(context.Background()).ApiAPIRoleConfigInput(apiAPIRoleConfigInput).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesApiAPI.ApiPutRoleConfig``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiPutRoleConfigRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **apiAPIRoleConfigInput** | [**ApiAPIRoleConfigInput**](ApiAPIRoleConfigInput.md) |  | 

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


## ApiPutTokens

> AuthAPITokensPutOutput ApiPutTokens(ctx).Username(username).Execute()

Tokens

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
	username := "username_example" // string | 

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.RolesApiAPI.ApiPutTokens(context.Background()).Username(username).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesApiAPI.ApiPutTokens``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiPutTokens`: AuthAPITokensPutOutput
	fmt.Fprintf(os.Stdout, "Response from `RolesApiAPI.ApiPutTokens`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiPutTokensRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **username** | **string** |  | 

### Return type

[**AuthAPITokensPutOutput**](AuthAPITokensPutOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiPutUsers

> ApiPutUsers(ctx).Username(username).AuthAPIUsersPutInput(authAPIUsersPutInput).Execute()

API Users

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
	username := "username_example" // string | 
	authAPIUsersPutInput := *openapiclient.NewAuthAPIUsersPutInput("Password_example", []openapiclient.AuthPermission{*openapiclient.NewAuthPermission()}) // AuthAPIUsersPutInput |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	r, err := apiClient.RolesApiAPI.ApiPutUsers(context.Background()).Username(username).AuthAPIUsersPutInput(authAPIUsersPutInput).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesApiAPI.ApiPutUsers``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiApiPutUsersRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **username** | **string** |  | 
 **authAPIUsersPutInput** | [**AuthAPIUsersPutInput**](AuthAPIUsersPutInput.md) |  | 

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


## ApiUsersMe

> AuthAPIMeOutput ApiUsersMe(ctx).Execute()

API Users

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
	resp, r, err := apiClient.RolesApiAPI.ApiUsersMe(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesApiAPI.ApiUsersMe``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ApiUsersMe`: AuthAPIMeOutput
	fmt.Fprintf(os.Stdout, "Response from `RolesApiAPI.ApiUsersMe`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiApiUsersMeRequest struct via the builder pattern


### Return type

[**AuthAPIMeOutput**](AuthAPIMeOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ToolsPing

> ApiAPIToolPingOutput ToolsPing(ctx).ApiAPIToolPingInput(apiAPIToolPingInput).Execute()

Ping tool

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
	apiAPIToolPingInput := *openapiclient.NewApiAPIToolPingInput() // ApiAPIToolPingInput |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.RolesApiAPI.ToolsPing(context.Background()).ApiAPIToolPingInput(apiAPIToolPingInput).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesApiAPI.ToolsPing``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ToolsPing`: ApiAPIToolPingOutput
	fmt.Fprintf(os.Stdout, "Response from `RolesApiAPI.ToolsPing`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiToolsPingRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **apiAPIToolPingInput** | [**ApiAPIToolPingInput**](ApiAPIToolPingInput.md) |  | 

### Return type

[**ApiAPIToolPingOutput**](ApiAPIToolPingOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ToolsPortmap

> ApiAPIToolPortmapOutput ToolsPortmap(ctx).ApiAPIToolPortmapInput(apiAPIToolPortmapInput).Execute()

Portmap tool

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
	apiAPIToolPortmapInput := *openapiclient.NewApiAPIToolPortmapInput() // ApiAPIToolPortmapInput |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.RolesApiAPI.ToolsPortmap(context.Background()).ApiAPIToolPortmapInput(apiAPIToolPortmapInput).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesApiAPI.ToolsPortmap``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ToolsPortmap`: ApiAPIToolPortmapOutput
	fmt.Fprintf(os.Stdout, "Response from `RolesApiAPI.ToolsPortmap`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiToolsPortmapRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **apiAPIToolPortmapInput** | [**ApiAPIToolPortmapInput**](ApiAPIToolPortmapInput.md) |  | 

### Return type

[**ApiAPIToolPortmapOutput**](ApiAPIToolPortmapOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ToolsTraceroute

> ApiAPIToolTracerouteOutput ToolsTraceroute(ctx).ApiAPIToolTracerouteInput(apiAPIToolTracerouteInput).Execute()

Traceroute tool

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
	apiAPIToolTracerouteInput := *openapiclient.NewApiAPIToolTracerouteInput() // ApiAPIToolTracerouteInput |  (optional)

	configuration := openapiclient.NewConfiguration()
	apiClient := openapiclient.NewAPIClient(configuration)
	resp, r, err := apiClient.RolesApiAPI.ToolsTraceroute(context.Background()).ApiAPIToolTracerouteInput(apiAPIToolTracerouteInput).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesApiAPI.ToolsTraceroute``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `ToolsTraceroute`: ApiAPIToolTracerouteOutput
	fmt.Fprintf(os.Stdout, "Response from `RolesApiAPI.ToolsTraceroute`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiToolsTracerouteRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **apiAPIToolTracerouteInput** | [**ApiAPIToolTracerouteInput**](ApiAPIToolTracerouteInput.md) |  | 

### Return type

[**ApiAPIToolTracerouteOutput**](ApiAPIToolTracerouteOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

