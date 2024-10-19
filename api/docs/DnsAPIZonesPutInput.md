# DnsAPIZonesPutInput

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Authoritative** | **bool** |  | 
**DefaultTTL** | **int32** |  | 
**HandlerConfigs** | **[]map[string]string** |  | 
**Hook** | **string** |  | 

## Methods

### NewDnsAPIZonesPutInput

`func NewDnsAPIZonesPutInput(authoritative bool, defaultTTL int32, handlerConfigs []map[string]string, hook string, ) *DnsAPIZonesPutInput`

NewDnsAPIZonesPutInput instantiates a new DnsAPIZonesPutInput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDnsAPIZonesPutInputWithDefaults

`func NewDnsAPIZonesPutInputWithDefaults() *DnsAPIZonesPutInput`

NewDnsAPIZonesPutInputWithDefaults instantiates a new DnsAPIZonesPutInput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAuthoritative

`func (o *DnsAPIZonesPutInput) GetAuthoritative() bool`

GetAuthoritative returns the Authoritative field if non-nil, zero value otherwise.

### GetAuthoritativeOk

`func (o *DnsAPIZonesPutInput) GetAuthoritativeOk() (*bool, bool)`

GetAuthoritativeOk returns a tuple with the Authoritative field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthoritative

`func (o *DnsAPIZonesPutInput) SetAuthoritative(v bool)`

SetAuthoritative sets Authoritative field to given value.


### GetDefaultTTL

`func (o *DnsAPIZonesPutInput) GetDefaultTTL() int32`

GetDefaultTTL returns the DefaultTTL field if non-nil, zero value otherwise.

### GetDefaultTTLOk

`func (o *DnsAPIZonesPutInput) GetDefaultTTLOk() (*int32, bool)`

GetDefaultTTLOk returns a tuple with the DefaultTTL field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefaultTTL

`func (o *DnsAPIZonesPutInput) SetDefaultTTL(v int32)`

SetDefaultTTL sets DefaultTTL field to given value.


### GetHandlerConfigs

`func (o *DnsAPIZonesPutInput) GetHandlerConfigs() []map[string]string`

GetHandlerConfigs returns the HandlerConfigs field if non-nil, zero value otherwise.

### GetHandlerConfigsOk

`func (o *DnsAPIZonesPutInput) GetHandlerConfigsOk() (*[]map[string]string, bool)`

GetHandlerConfigsOk returns a tuple with the HandlerConfigs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHandlerConfigs

`func (o *DnsAPIZonesPutInput) SetHandlerConfigs(v []map[string]string)`

SetHandlerConfigs sets HandlerConfigs field to given value.


### SetHandlerConfigsNil

`func (o *DnsAPIZonesPutInput) SetHandlerConfigsNil(b bool)`

 SetHandlerConfigsNil sets the value for HandlerConfigs to be an explicit nil

### UnsetHandlerConfigs
`func (o *DnsAPIZonesPutInput) UnsetHandlerConfigs()`

UnsetHandlerConfigs ensures that no value is present for HandlerConfigs, not even an explicit nil
### GetHook

`func (o *DnsAPIZonesPutInput) GetHook() string`

GetHook returns the Hook field if non-nil, zero value otherwise.

### GetHookOk

`func (o *DnsAPIZonesPutInput) GetHookOk() (*string, bool)`

GetHookOk returns a tuple with the Hook field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHook

`func (o *DnsAPIZonesPutInput) SetHook(v string)`

SetHook sets Hook field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


