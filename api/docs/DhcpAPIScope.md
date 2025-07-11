# DhcpAPIScope

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Default** | **bool** |  | 
**Dns** | Pointer to [**DhcpScopeDNS**](DhcpScopeDNS.md) |  | [optional] 
**Hook** | **string** |  | 
**Ipam** | **map[string]string** |  | 
**Options** | [**[]TypesOption**](TypesOption.md) |  | 
**Scope** | **string** |  | 
**Statistics** | [**DhcpAPIScopeStatistics**](DhcpAPIScopeStatistics.md) |  | 
**SubnetCidr** | **string** |  | 
**Ttl** | **int32** |  | 

## Methods

### NewDhcpAPIScope

`func NewDhcpAPIScope(default_ bool, hook string, ipam map[string]string, options []TypesOption, scope string, statistics DhcpAPIScopeStatistics, subnetCidr string, ttl int32, ) *DhcpAPIScope`

NewDhcpAPIScope instantiates a new DhcpAPIScope object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDhcpAPIScopeWithDefaults

`func NewDhcpAPIScopeWithDefaults() *DhcpAPIScope`

NewDhcpAPIScopeWithDefaults instantiates a new DhcpAPIScope object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDefault

`func (o *DhcpAPIScope) GetDefault() bool`

GetDefault returns the Default field if non-nil, zero value otherwise.

### GetDefaultOk

`func (o *DhcpAPIScope) GetDefaultOk() (*bool, bool)`

GetDefaultOk returns a tuple with the Default field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefault

`func (o *DhcpAPIScope) SetDefault(v bool)`

SetDefault sets Default field to given value.


### GetDns

`func (o *DhcpAPIScope) GetDns() DhcpScopeDNS`

GetDns returns the Dns field if non-nil, zero value otherwise.

### GetDnsOk

`func (o *DhcpAPIScope) GetDnsOk() (*DhcpScopeDNS, bool)`

GetDnsOk returns a tuple with the Dns field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDns

`func (o *DhcpAPIScope) SetDns(v DhcpScopeDNS)`

SetDns sets Dns field to given value.

### HasDns

`func (o *DhcpAPIScope) HasDns() bool`

HasDns returns a boolean if a field has been set.

### GetHook

`func (o *DhcpAPIScope) GetHook() string`

GetHook returns the Hook field if non-nil, zero value otherwise.

### GetHookOk

`func (o *DhcpAPIScope) GetHookOk() (*string, bool)`

GetHookOk returns a tuple with the Hook field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHook

`func (o *DhcpAPIScope) SetHook(v string)`

SetHook sets Hook field to given value.


### GetIpam

`func (o *DhcpAPIScope) GetIpam() map[string]string`

GetIpam returns the Ipam field if non-nil, zero value otherwise.

### GetIpamOk

`func (o *DhcpAPIScope) GetIpamOk() (*map[string]string, bool)`

GetIpamOk returns a tuple with the Ipam field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpam

`func (o *DhcpAPIScope) SetIpam(v map[string]string)`

SetIpam sets Ipam field to given value.


### SetIpamNil

`func (o *DhcpAPIScope) SetIpamNil(b bool)`

 SetIpamNil sets the value for Ipam to be an explicit nil

### UnsetIpam
`func (o *DhcpAPIScope) UnsetIpam()`

UnsetIpam ensures that no value is present for Ipam, not even an explicit nil
### GetOptions

`func (o *DhcpAPIScope) GetOptions() []TypesOption`

GetOptions returns the Options field if non-nil, zero value otherwise.

### GetOptionsOk

`func (o *DhcpAPIScope) GetOptionsOk() (*[]TypesOption, bool)`

GetOptionsOk returns a tuple with the Options field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOptions

`func (o *DhcpAPIScope) SetOptions(v []TypesOption)`

SetOptions sets Options field to given value.


### SetOptionsNil

`func (o *DhcpAPIScope) SetOptionsNil(b bool)`

 SetOptionsNil sets the value for Options to be an explicit nil

### UnsetOptions
`func (o *DhcpAPIScope) UnsetOptions()`

UnsetOptions ensures that no value is present for Options, not even an explicit nil
### GetScope

`func (o *DhcpAPIScope) GetScope() string`

GetScope returns the Scope field if non-nil, zero value otherwise.

### GetScopeOk

`func (o *DhcpAPIScope) GetScopeOk() (*string, bool)`

GetScopeOk returns a tuple with the Scope field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScope

`func (o *DhcpAPIScope) SetScope(v string)`

SetScope sets Scope field to given value.


### GetStatistics

`func (o *DhcpAPIScope) GetStatistics() DhcpAPIScopeStatistics`

GetStatistics returns the Statistics field if non-nil, zero value otherwise.

### GetStatisticsOk

`func (o *DhcpAPIScope) GetStatisticsOk() (*DhcpAPIScopeStatistics, bool)`

GetStatisticsOk returns a tuple with the Statistics field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatistics

`func (o *DhcpAPIScope) SetStatistics(v DhcpAPIScopeStatistics)`

SetStatistics sets Statistics field to given value.


### GetSubnetCidr

`func (o *DhcpAPIScope) GetSubnetCidr() string`

GetSubnetCidr returns the SubnetCidr field if non-nil, zero value otherwise.

### GetSubnetCidrOk

`func (o *DhcpAPIScope) GetSubnetCidrOk() (*string, bool)`

GetSubnetCidrOk returns a tuple with the SubnetCidr field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubnetCidr

`func (o *DhcpAPIScope) SetSubnetCidr(v string)`

SetSubnetCidr sets SubnetCidr field to given value.


### GetTtl

`func (o *DhcpAPIScope) GetTtl() int32`

GetTtl returns the Ttl field if non-nil, zero value otherwise.

### GetTtlOk

`func (o *DhcpAPIScope) GetTtlOk() (*int32, bool)`

GetTtlOk returns a tuple with the Ttl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTtl

`func (o *DhcpAPIScope) SetTtl(v int32)`

SetTtl sets Ttl field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


