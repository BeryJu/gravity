# \RolesApiApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**ApiAuthConfig**](RolesApiApi.md#ApiAuthConfig) | **Get** /api/v1/auth/config | API Users
[**ApiDeleteTokens**](RolesApiApi.md#ApiDeleteTokens) | **Delete** /api/v1/auth/tokens | Tokens
[**ApiDeleteUsers**](RolesApiApi.md#ApiDeleteUsers) | **Delete** /api/v1/auth/users | API Users
[**ApiGetMembers**](RolesApiApi.md#ApiGetMembers) | **Get** /api/v1/etcd/members | Etcd members
[**ApiGetMetricsMemory**](RolesApiApi.md#ApiGetMetricsMemory) | **Get** /api/v1/system/metrics/memory | System Metrics
[**ApiGetRoleConfig**](RolesApiApi.md#ApiGetRoleConfig) | **Get** /api/v1/roles/api | API role config
[**ApiGetTokens**](RolesApiApi.md#ApiGetTokens) | **Get** /api/v1/auth/tokens | Tokens
[**ApiGetUsers**](RolesApiApi.md#ApiGetUsers) | **Get** /api/v1/auth/users | API Users
[**ApiLoginUser**](RolesApiApi.md#ApiLoginUser) | **Post** /api/v1/auth/login | API Users
[**ApiPutRoleConfig**](RolesApiApi.md#ApiPutRoleConfig) | **Post** /api/v1/roles/api | API role config
[**ApiPutTokens**](RolesApiApi.md#ApiPutTokens) | **Post** /api/v1/auth/tokens | Tokens
[**ApiPutUsers**](RolesApiApi.md#ApiPutUsers) | **Post** /api/v1/auth/users | API Users
[**ApiUsersMe**](RolesApiApi.md#ApiUsersMe) | **Get** /api/v1/auth/me | API Users



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
    openapiclient "./openapi"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesApiApi.ApiAuthConfig(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesApiApi.ApiAuthConfig``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ApiAuthConfig`: AuthAPIConfigOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesApiApi.ApiAuthConfig`: %v\n", resp)
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
    openapiclient "./openapi"
)

func main() {
    key := "key_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesApiApi.ApiDeleteTokens(context.Background()).Key(key).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesApiApi.ApiDeleteTokens``: %v\n", err)
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
    openapiclient "./openapi"
)

func main() {
    username := "username_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesApiApi.ApiDeleteUsers(context.Background()).Username(username).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesApiApi.ApiDeleteUsers``: %v\n", err)
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


## ApiGetMembers

> ApiAPIMembersOutput ApiGetMembers(ctx).Execute()

Etcd members

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
    resp, r, err := apiClient.RolesApiApi.ApiGetMembers(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesApiApi.ApiGetMembers``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ApiGetMembers`: ApiAPIMembersOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesApiApi.ApiGetMembers`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiApiGetMembersRequest struct via the builder pattern


### Return type

[**ApiAPIMembersOutput**](ApiAPIMembersOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## ApiGetMetricsMemory

> TypesAPIMetricsGetOutput ApiGetMetricsMemory(ctx).Execute()

System Metrics

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
    resp, r, err := apiClient.RolesApiApi.ApiGetMetricsMemory(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesApiApi.ApiGetMetricsMemory``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ApiGetMetricsMemory`: TypesAPIMetricsGetOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesApiApi.ApiGetMetricsMemory`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiApiGetMetricsMemoryRequest struct via the builder pattern


### Return type

[**TypesAPIMetricsGetOutput**](TypesAPIMetricsGetOutput.md)

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
    openapiclient "./openapi"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesApiApi.ApiGetRoleConfig(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesApiApi.ApiGetRoleConfig``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ApiGetRoleConfig`: ApiAPIRoleConfigOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesApiApi.ApiGetRoleConfig`: %v\n", resp)
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
    openapiclient "./openapi"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesApiApi.ApiGetTokens(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesApiApi.ApiGetTokens``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ApiGetTokens`: AuthAPITokensGetOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesApiApi.ApiGetTokens`: %v\n", resp)
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

> AuthAPIUsersGetOutput ApiGetUsers(ctx).Execute()

API Users

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
    resp, r, err := apiClient.RolesApiApi.ApiGetUsers(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesApiApi.ApiGetUsers``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ApiGetUsers`: AuthAPIUsersGetOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesApiApi.ApiGetUsers`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiApiGetUsersRequest struct via the builder pattern


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
    openapiclient "./openapi"
)

func main() {
    authAPILoginInput := *openapiclient.NewAuthAPILoginInput() // AuthAPILoginInput |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesApiApi.ApiLoginUser(context.Background()).AuthAPILoginInput(authAPILoginInput).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesApiApi.ApiLoginUser``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ApiLoginUser`: AuthAPILoginOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesApiApi.ApiLoginUser`: %v\n", resp)
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
    openapiclient "./openapi"
)

func main() {
    apiAPIRoleConfigInput := *openapiclient.NewApiAPIRoleConfigInput(*openapiclient.NewApiRoleConfig()) // ApiAPIRoleConfigInput |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesApiApi.ApiPutRoleConfig(context.Background()).ApiAPIRoleConfigInput(apiAPIRoleConfigInput).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesApiApi.ApiPutRoleConfig``: %v\n", err)
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
    openapiclient "./openapi"
)

func main() {
    username := "username_example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesApiApi.ApiPutTokens(context.Background()).Username(username).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesApiApi.ApiPutTokens``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ApiPutTokens`: AuthAPITokensPutOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesApiApi.ApiPutTokens`: %v\n", resp)
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
    openapiclient "./openapi"
)

func main() {
    username := "username_example" // string | 
    authAPIUsersPutInput := *openapiclient.NewAuthAPIUsersPutInput("Password_example") // AuthAPIUsersPutInput |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesApiApi.ApiPutUsers(context.Background()).Username(username).AuthAPIUsersPutInput(authAPIUsersPutInput).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesApiApi.ApiPutUsers``: %v\n", err)
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
    openapiclient "./openapi"
)

func main() {

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesApiApi.ApiUsersMe(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesApiApi.ApiUsersMe``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `ApiUsersMe`: AuthAPIMeOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesApiApi.ApiUsersMe`: %v\n", resp)
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

