# \RolesTftpApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**TftpDeleteFiles**](RolesTftpApi.md#TftpDeleteFiles) | **Delete** /api/v1/tftp/files | TFTP Files
[**TftpDownloadFiles**](RolesTftpApi.md#TftpDownloadFiles) | **Get** /api/v1/tftp/files/download | TFTP Files
[**TftpGetFiles**](RolesTftpApi.md#TftpGetFiles) | **Get** /api/v1/tftp/files | TFTP Files
[**TftpGetRoleConfig**](RolesTftpApi.md#TftpGetRoleConfig) | **Get** /api/v1/roles/tftp | TFTP role config
[**TftpPutFiles**](RolesTftpApi.md#TftpPutFiles) | **Post** /api/v1/tftp/files | TFTP Files
[**TftpPutRoleConfig**](RolesTftpApi.md#TftpPutRoleConfig) | **Post** /api/v1/roles/tftp | TFTP role config



## TftpDeleteFiles

> TftpDeleteFiles(ctx).Host(host).Name(name).Execute()

TFTP Files

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
    host := "host_example" // string |  (optional)
    name := "name_example" // string |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.RolesTftpApi.TftpDeleteFiles(context.Background()).Host(host).Name(name).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesTftpApi.TftpDeleteFiles``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiTftpDeleteFilesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **host** | **string** |  | 
 **name** | **string** |  | 

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


## TftpDownloadFiles

> TftpAPIFilesDownloadOutput TftpDownloadFiles(ctx).Host(host).Name(name).Execute()

TFTP Files

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
    host := "host_example" // string |  (optional)
    name := "name_example" // string |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesTftpApi.TftpDownloadFiles(context.Background()).Host(host).Name(name).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesTftpApi.TftpDownloadFiles``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `TftpDownloadFiles`: TftpAPIFilesDownloadOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesTftpApi.TftpDownloadFiles`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiTftpDownloadFilesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **host** | **string** |  | 
 **name** | **string** |  | 

### Return type

[**TftpAPIFilesDownloadOutput**](TftpAPIFilesDownloadOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## TftpGetFiles

> TftpAPIFilesGetOutput TftpGetFiles(ctx).Execute()

TFTP Files

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
    resp, r, err := apiClient.RolesTftpApi.TftpGetFiles(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesTftpApi.TftpGetFiles``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `TftpGetFiles`: TftpAPIFilesGetOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesTftpApi.TftpGetFiles`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiTftpGetFilesRequest struct via the builder pattern


### Return type

[**TftpAPIFilesGetOutput**](TftpAPIFilesGetOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## TftpGetRoleConfig

> TftpAPIRoleConfigOutput TftpGetRoleConfig(ctx).Execute()

TFTP role config

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
    resp, r, err := apiClient.RolesTftpApi.TftpGetRoleConfig(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesTftpApi.TftpGetRoleConfig``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `TftpGetRoleConfig`: TftpAPIRoleConfigOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesTftpApi.TftpGetRoleConfig`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiTftpGetRoleConfigRequest struct via the builder pattern


### Return type

[**TftpAPIRoleConfigOutput**](TftpAPIRoleConfigOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## TftpPutFiles

> TftpPutFiles(ctx).TftpAPIFilesPutInput(tftpAPIFilesPutInput).Execute()

TFTP Files

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
    tftpAPIFilesPutInput := *openapiclient.NewTftpAPIFilesPutInput("Data_example", "Host_example", "Name_example") // TftpAPIFilesPutInput |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.RolesTftpApi.TftpPutFiles(context.Background()).TftpAPIFilesPutInput(tftpAPIFilesPutInput).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesTftpApi.TftpPutFiles``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiTftpPutFilesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **tftpAPIFilesPutInput** | [**TftpAPIFilesPutInput**](TftpAPIFilesPutInput.md) |  | 

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


## TftpPutRoleConfig

> TftpPutRoleConfig(ctx).TftpAPIRoleConfigInput(tftpAPIRoleConfigInput).Execute()

TFTP role config

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
    tftpAPIRoleConfigInput := *openapiclient.NewTftpAPIRoleConfigInput(*openapiclient.NewTftpRoleConfig()) // TftpAPIRoleConfigInput |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.RolesTftpApi.TftpPutRoleConfig(context.Background()).TftpAPIRoleConfigInput(tftpAPIRoleConfigInput).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesTftpApi.TftpPutRoleConfig``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiTftpPutRoleConfigRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **tftpAPIRoleConfigInput** | [**TftpAPIRoleConfigInput**](TftpAPIRoleConfigInput.md) |  | 

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

