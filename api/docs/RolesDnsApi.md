# \RolesDnsApi

All URIs are relative to *http://localhost*

Method | HTTP request | Description
------------- | ------------- | -------------
[**DnsDeleteRecords**](RolesDnsApi.md#DnsDeleteRecords) | **Delete** /api/v1/dns/zones/records | DNS Records
[**DnsDeleteZones**](RolesDnsApi.md#DnsDeleteZones) | **Delete** /api/v1/dns/zones | DNS Zones
[**DnsGetRecords**](RolesDnsApi.md#DnsGetRecords) | **Get** /api/v1/dns/zones/records | DNS Records
[**DnsGetRoleConfig**](RolesDnsApi.md#DnsGetRoleConfig) | **Get** /api/v1/roles/dns | DNS role config
[**DnsGetZones**](RolesDnsApi.md#DnsGetZones) | **Get** /api/v1/dns/zones | DNS Zones
[**DnsImportZones**](RolesDnsApi.md#DnsImportZones) | **Post** /api/v1/dns/zones/import | DNS Zones
[**DnsPutRecords**](RolesDnsApi.md#DnsPutRecords) | **Post** /api/v1/dns/zones/records | DNS Records
[**DnsPutRoleConfig**](RolesDnsApi.md#DnsPutRoleConfig) | **Post** /api/v1/roles/dns | DNS role config
[**DnsPutZones**](RolesDnsApi.md#DnsPutZones) | **Post** /api/v1/dns/zones | DNS Zones



## DnsDeleteRecords

> DnsDeleteRecords(ctx).Zone(zone).Hostname(hostname).Uid(uid).Type_(type_).Execute()

DNS Records

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
    zone := "zone_example" // string | 
    hostname := "hostname_example" // string | 
    uid := "uid_example" // string | 
    type_ := "type__example" // string | 

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.RolesDnsApi.DnsDeleteRecords(context.Background()).Zone(zone).Hostname(hostname).Uid(uid).Type_(type_).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDnsApi.DnsDeleteRecords``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDnsDeleteRecordsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **zone** | **string** |  | 
 **hostname** | **string** |  | 
 **uid** | **string** |  | 
 **type_** | **string** |  | 

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


## DnsDeleteZones

> DnsDeleteZones(ctx).Zone(zone).Execute()

DNS Zones

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
    zone := "zone_example" // string |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.RolesDnsApi.DnsDeleteZones(context.Background()).Zone(zone).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDnsApi.DnsDeleteZones``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDnsDeleteZonesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **zone** | **string** |  | 

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


## DnsGetRecords

> DnsAPIRecordsGetOutput DnsGetRecords(ctx).Zone(zone).Hostname(hostname).Type_(type_).Uid(uid).Execute()

DNS Records

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
    zone := "zone_example" // string |  (optional)
    hostname := "hostname_example" // string | Optionally get DNS Records for hostname (optional)
    type_ := "type__example" // string |  (optional)
    uid := "uid_example" // string |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesDnsApi.DnsGetRecords(context.Background()).Zone(zone).Hostname(hostname).Type_(type_).Uid(uid).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDnsApi.DnsGetRecords``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DnsGetRecords`: DnsAPIRecordsGetOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesDnsApi.DnsGetRecords`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDnsGetRecordsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **zone** | **string** |  | 
 **hostname** | **string** | Optionally get DNS Records for hostname | 
 **type_** | **string** |  | 
 **uid** | **string** |  | 

### Return type

[**DnsAPIRecordsGetOutput**](DnsAPIRecordsGetOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DnsGetRoleConfig

> DnsAPIRoleConfigOutput DnsGetRoleConfig(ctx).Execute()

DNS role config

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
    resp, r, err := apiClient.RolesDnsApi.DnsGetRoleConfig(context.Background()).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDnsApi.DnsGetRoleConfig``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DnsGetRoleConfig`: DnsAPIRoleConfigOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesDnsApi.DnsGetRoleConfig`: %v\n", resp)
}
```

### Path Parameters

This endpoint does not need any parameter.

### Other Parameters

Other parameters are passed through a pointer to a apiDnsGetRoleConfigRequest struct via the builder pattern


### Return type

[**DnsAPIRoleConfigOutput**](DnsAPIRoleConfigOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DnsGetZones

> DnsAPIZonesGetOutput DnsGetZones(ctx).Name(name).Execute()

DNS Zones

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
    name := "name_example" // string | Optionally get DNS Zone by name (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesDnsApi.DnsGetZones(context.Background()).Name(name).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDnsApi.DnsGetZones``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DnsGetZones`: DnsAPIZonesGetOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesDnsApi.DnsGetZones`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDnsGetZonesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **name** | **string** | Optionally get DNS Zone by name | 

### Return type

[**DnsAPIZonesGetOutput**](DnsAPIZonesGetOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: Not defined
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DnsImportZones

> DnsAPIZonesImportOutput DnsImportZones(ctx).Zone(zone).DnsAPIZonesImportInput(dnsAPIZonesImportInput).Execute()

DNS Zones

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
    zone := "zone_example" // string |  (optional)
    dnsAPIZonesImportInput := *openapiclient.NewDnsAPIZonesImportInput() // DnsAPIZonesImportInput |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    resp, r, err := apiClient.RolesDnsApi.DnsImportZones(context.Background()).Zone(zone).DnsAPIZonesImportInput(dnsAPIZonesImportInput).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDnsApi.DnsImportZones``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
    // response from `DnsImportZones`: DnsAPIZonesImportOutput
    fmt.Fprintf(os.Stdout, "Response from `RolesDnsApi.DnsImportZones`: %v\n", resp)
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDnsImportZonesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **zone** | **string** |  | 
 **dnsAPIZonesImportInput** | [**DnsAPIZonesImportInput**](DnsAPIZonesImportInput.md) |  | 

### Return type

[**DnsAPIZonesImportOutput**](DnsAPIZonesImportOutput.md)

### Authorization

No authorization required

### HTTP request headers

- **Content-Type**: application/json
- **Accept**: application/json

[[Back to top]](#) [[Back to API list]](../README.md#documentation-for-api-endpoints)
[[Back to Model list]](../README.md#documentation-for-models)
[[Back to README]](../README.md)


## DnsPutRecords

> DnsPutRecords(ctx).Zone(zone).Hostname(hostname).Uid(uid).DnsAPIRecordsPutInput(dnsAPIRecordsPutInput).Execute()

DNS Records

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
    zone := "zone_example" // string | 
    hostname := "hostname_example" // string | 
    uid := "uid_example" // string |  (optional)
    dnsAPIRecordsPutInput := *openapiclient.NewDnsAPIRecordsPutInput("Data_example", "Type_example") // DnsAPIRecordsPutInput |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.RolesDnsApi.DnsPutRecords(context.Background()).Zone(zone).Hostname(hostname).Uid(uid).DnsAPIRecordsPutInput(dnsAPIRecordsPutInput).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDnsApi.DnsPutRecords``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDnsPutRecordsRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **zone** | **string** |  | 
 **hostname** | **string** |  | 
 **uid** | **string** |  | 
 **dnsAPIRecordsPutInput** | [**DnsAPIRecordsPutInput**](DnsAPIRecordsPutInput.md) |  | 

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


## DnsPutRoleConfig

> DnsPutRoleConfig(ctx).DnsAPIRoleConfigInput(dnsAPIRoleConfigInput).Execute()

DNS role config

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
    dnsAPIRoleConfigInput := *openapiclient.NewDnsAPIRoleConfigInput(*openapiclient.NewDnsRoleConfig()) // DnsAPIRoleConfigInput |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.RolesDnsApi.DnsPutRoleConfig(context.Background()).DnsAPIRoleConfigInput(dnsAPIRoleConfigInput).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDnsApi.DnsPutRoleConfig``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDnsPutRoleConfigRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **dnsAPIRoleConfigInput** | [**DnsAPIRoleConfigInput**](DnsAPIRoleConfigInput.md) |  | 

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


## DnsPutZones

> DnsPutZones(ctx).Zone(zone).DnsAPIZonesPutInput(dnsAPIZonesPutInput).Execute()

DNS Zones

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
    zone := "zone_example" // string | 
    dnsAPIZonesPutInput := *openapiclient.NewDnsAPIZonesPutInput(false, int32(123), []map[string]interface{}{map[string]interface{}{"key": interface{}(123)}}, "Hook_example") // DnsAPIZonesPutInput |  (optional)

    configuration := openapiclient.NewConfiguration()
    apiClient := openapiclient.NewAPIClient(configuration)
    r, err := apiClient.RolesDnsApi.DnsPutZones(context.Background()).Zone(zone).DnsAPIZonesPutInput(dnsAPIZonesPutInput).Execute()
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error when calling `RolesDnsApi.DnsPutZones``: %v\n", err)
        fmt.Fprintf(os.Stderr, "Full HTTP response: %v\n", r)
    }
}
```

### Path Parameters



### Other Parameters

Other parameters are passed through a pointer to a apiDnsPutZonesRequest struct via the builder pattern


Name | Type | Description  | Notes
------------- | ------------- | ------------- | -------------
 **zone** | **string** |  | 
 **dnsAPIZonesPutInput** | [**DnsAPIZonesPutInput**](DnsAPIZonesPutInput.md) |  | 

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

