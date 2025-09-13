# \RolesTftpAPI

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**TftpDeleteFiles**](RolesTftpAPI.md#TftpDeleteFiles) | **Delete** /api/v1/tftp/files | TFTP Files
[**TftpDownloadFiles**](RolesTftpAPI.md#TftpDownloadFiles) | **Get** /api/v1/tftp/files/download | TFTP Files
[**TftpGetFiles**](RolesTftpAPI.md#TftpGetFiles) | **Get** /api/v1/tftp/files | TFTP Files
[**TftpGetRoleConfig**](RolesTftpAPI.md#TftpGetRoleConfig) | **Get** /api/v1/roles/tftp | TFTP role config
[**TftpPutFiles**](RolesTftpAPI.md#TftpPutFiles) | **Post** /api/v1/tftp/files | TFTP Files
[**TftpPutRoleConfig**](RolesTftpAPI.md#TftpPutRoleConfig) | **Post** /api/v1/roles/tftp | TFTP role config



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
	r, err := apiClient.RolesTftpAPI.TftpDeleteFiles(context.Background()).Host(host).Name(name).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesTftpAPI.TftpDeleteFiles``: %v\n", err)
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
	resp, r, err := apiClient.RolesTftpAPI.TftpDownloadFiles(context.Background()).Host(host).Name(name).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesTftpAPI.TftpDownloadFiles``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `TftpDownloadFiles`: TftpAPIFilesDownloadOutput
	fmt.Fprintf(os.Stdout, "Response from `RolesTftpAPI.TftpDownloadFiles`: %v\n", resp)
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
	resp, r, err := apiClient.RolesTftpAPI.TftpGetFiles(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesTftpAPI.TftpGetFiles``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `TftpGetFiles`: TftpAPIFilesGetOutput
	fmt.Fprintf(os.Stdout, "Response from `RolesTftpAPI.TftpGetFiles`: %v\n", resp)
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
	resp, r, err := apiClient.RolesTftpAPI.TftpGetRoleConfig(context.Background()).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesTftpAPI.TftpGetRoleConfig``: %v\n", err)
		fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
	}
	// response from `TftpGetRoleConfig`: TftpAPIRoleConfigOutput
	fmt.Fprintf(os.Stdout, "Response from `RolesTftpAPI.TftpGetRoleConfig`: %v\n", resp)
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
	r, err := apiClient.RolesTftpAPI.TftpPutFiles(context.Background()).TftpAPIFilesPutInput(tftpAPIFilesPutInput).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesTftpAPI.TftpPutFiles``: %v\n", err)
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
	r, err := apiClient.RolesTftpAPI.TftpPutRoleConfig(context.Background()).TftpAPIRoleConfigInput(tftpAPIRoleConfigInput).Execute()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error when calling `RolesTftpAPI.TftpPutRoleConfig``: %v\n", err)
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

