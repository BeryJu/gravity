# \RolesBackupApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**BackupGetRoleConfig**](RolesBackupApi.md#BackupGetRoleConfig) | **Get** /api/v1/roles/backup | Backup role config
[**BackupPutRoleConfig**](RolesBackupApi.md#BackupPutRoleConfig) | **Post** /api/v1/roles/backup | Backup role config
[**BackupStart**](RolesBackupApi.md#BackupStart) | **Post** /api/v1/backup/start | Backup start
[**BackupStatus**](RolesBackupApi.md#BackupStatus) | **Get** /api/v1/backup/status | Backup status



## BackupGetRoleConfig

> BackupAPIRoleConfigOutput BackupGetRoleConfig(ctx).Execute()

Backup role config

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
    resp, r, err := apiClient.RolesBackupApi.BackupGetRoleConfig(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesBackupApi.BackupGetRoleConfig``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `BackupGetRoleConfig`: BackupAPIRoleConfigOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesBackupApi.BackupGetRoleConfig`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiBackupGetRoleConfigRequest struct via the builder pattern


### Return type

[**BackupAPIRoleConfigOutput**](BackupAPIRoleConfigOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## BackupPutRoleConfig

> BackupPutRoleConfig(ctx).BackupAPIRoleConfigInput(backupAPIRoleConfigInput).Execute()

Backup role config

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
    backupAPIRoleConfigInput := *openapiclient.NewBackupAPIRoleConfigInput(*openapiclient.NewBackupRoleConfig()) // BackupAPIRoleConfigInput |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.RolesBackupApi.BackupPutRoleConfig(context.Background()).BackupAPIRoleConfigInput(backupAPIRoleConfigInput).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesBackupApi.BackupPutRoleConfig``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiBackupPutRoleConfigRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **backupAPIRoleConfigInput** | [**BackupAPIRoleConfigInput**](BackupAPIRoleConfigInput.md) |  | 

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


## BackupStart

> BackupBackupStatus BackupStart(ctx).Wait(wait).Execute()

Backup start

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
    wait := true // bool | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesBackupApi.BackupStart(context.Background()).Wait(wait).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesBackupApi.BackupStart``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `BackupStart`: BackupBackupStatus
    fmt.Fprintf(os.Stdout, "Response from `RolesBackupApi.BackupStart`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiBackupStartRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **wait** | **bool** |  | 

### Return type

[**BackupBackupStatus**](BackupBackupStatus.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## BackupStatus

> BackupAPIBackupStatusOutput BackupStatus(ctx).Execute()

Backup status

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
    resp, r, err := apiClient.RolesBackupApi.BackupStatus(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesBackupApi.BackupStatus``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `BackupStatus`: BackupAPIBackupStatusOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesBackupApi.BackupStatus`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiBackupStatusRequest struct via the builder pattern


### Return type

[**BackupAPIBackupStatusOutput**](BackupAPIBackupStatusOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)

