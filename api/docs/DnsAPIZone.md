# DnsAPIZone

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Authoritative** | **bool** |  | 
**DefaultTTL** | **int32** |  | 
**HandlerConfigs** | **[]map[string]interface{}** |  | 
**Hook** | **string** |  | 
**Name** | **string** |  | 

## Methods

### NewDnsAPIZone

`func NewDnsAPIZone(authoritative bool, defaultTTL int32, handlerConfigs []map[string]interface{}, hook string, name string, ) *DnsAPIZone`

NewDnsAPIZone instantiates a new DnsAPIZone object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDnsAPIZoneWithDefaults

`func NewDnsAPIZoneWithDefaults() *DnsAPIZone`

NewDnsAPIZoneWithDefaults instantiates a new DnsAPIZone object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAuthoritative

`func (o *DnsAPIZone) GetAuthoritative() bool`

GetAuthoritative returns the Authoritative field if non-nil, zero value otherwise.

### GetAuthoritativeOk

`func (o *DnsAPIZone) GetAuthoritativeOk() (*bool, bool)`

GetAuthoritativeOk returns a tuple with the Authoritative field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthoritative

`func (o *DnsAPIZone) SetAuthoritative(v bool)`

SetAuthoritative sets Authoritative field to given value.


### GetDefaultTTL

`func (o *DnsAPIZone) GetDefaultTTL() int32`

GetDefaultTTL returns the DefaultTTL field if non-nil, zero value otherwise.

### GetDefaultTTLOk

`func (o *DnsAPIZone) GetDefaultTTLOk() (*int32, bool)`

GetDefaultTTLOk returns a tuple with the DefaultTTL field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultTTL

`func (o *DnsAPIZone) SetDefaultTTL(v int32)`

SetDefaultTTL sets DefaultTTL field to given value.


### GetHandlerConfigs

`func (o *DnsAPIZone) GetHandlerConfigs() []map[string]interface{}`

GetHandlerConfigs returns the HandlerConfigs field if non-nil, zero value otherwise.

### GetHandlerConfigsOk

`func (o *DnsAPIZone) GetHandlerConfigsOk() (*[]map[string]interface{}, bool)`

GetHandlerConfigsOk returns a tuple with the HandlerConfigs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHandlerConfigs

`func (o *DnsAPIZone) SetHandlerConfigs(v []map[string]interface{})`

SetHandlerConfigs sets HandlerConfigs field to given value.


### SetHandlerConfigsNil

`func (o *DnsAPIZone) SetHandlerConfigsNil(b bool)`

 SetHandlerConfigsNil sets the value for HandlerConfigs to be an explicit nil

### UnsetHandlerConfigs
`func (o *DnsAPIZone) UnsetHandlerConfigs()`

UnsetHandlerConfigs ensures that no value is present for HandlerConfigs, not even an explicit nil
### GetHook

`func (o *DnsAPIZone) GetHook() string`

GetHook returns the Hook field if non-nil, zero value otherwise.

### GetHookOk

`func (o *DnsAPIZone) GetHookOk() (*string, bool)`

GetHookOk returns a tuple with the Hook field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHook

`func (o *DnsAPIZone) SetHook(v string)`

SetHook sets Hook field to given value.


### GetName

`func (o *DnsAPIZone) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *DnsAPIZone) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *DnsAPIZone) SetName(v string)`

SetName sets Name field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


