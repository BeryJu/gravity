# DhcpAPIScopesPutInput

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Default** | **bool** |  | 
**Dns** | Pointer to [**DhcpScopeDNS**](DhcpScopeDNS.md) |  | [optional] 
**Hook** | **string** |  | 
**Ipam** | Pointer to **map[string]string** |  | [optional] 
**Options** | [**[]TypesOption**](TypesOption.md) |  | 
**SubnetCidr** | **string** |  | 
**Ttl** | **int32** |  | 

## Methods

### NewDhcpAPIScopesPutInput

`func NewDhcpAPIScopesPutInput(default_ bool, hook string, options []TypesOption, subnetCidr string, ttl int32, ) *DhcpAPIScopesPutInput`

NewDhcpAPIScopesPutInput instantiates a new DhcpAPIScopesPutInput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDhcpAPIScopesPutInputWithDefaults

`func NewDhcpAPIScopesPutInputWithDefaults() *DhcpAPIScopesPutInput`

NewDhcpAPIScopesPutInputWithDefaults instantiates a new DhcpAPIScopesPutInput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDefault

`func (o *DhcpAPIScopesPutInput) GetDefault() bool`

GetDefault returns the Default field if non-nil, zero value otherwise.

### GetDefaultOk

`func (o *DhcpAPIScopesPutInput) GetDefaultOk() (*bool, bool)`

GetDefaultOk returns a tuple with the Default field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDefault

`func (o *DhcpAPIScopesPutInput) SetDefault(v bool)`

SetDefault sets Default field to given value.


### GetDns

`func (o *DhcpAPIScopesPutInput) GetDns() DhcpScopeDNS`

GetDns returns the Dns field if non-nil, zero value otherwise.

### GetDnsOk

`func (o *DhcpAPIScopesPutInput) GetDnsOk() (*DhcpScopeDNS, bool)`

GetDnsOk returns a tuple with the Dns field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDns

`func (o *DhcpAPIScopesPutInput) SetDns(v DhcpScopeDNS)`

SetDns sets Dns field to given value.

### HasDns

`func (o *DhcpAPIScopesPutInput) HasDns() bool`

HasDns returns a boolean if a field has been set.

### GetHook

`func (o *DhcpAPIScopesPutInput) GetHook() string`

GetHook returns the Hook field if non-nil, zero value otherwise.

### GetHookOk

`func (o *DhcpAPIScopesPutInput) GetHookOk() (*string, bool)`

GetHookOk returns a tuple with the Hook field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHook

`func (o *DhcpAPIScopesPutInput) SetHook(v string)`

SetHook sets Hook field to given value.


### GetIpam

`func (o *DhcpAPIScopesPutInput) GetIpam() map[string]string`

GetIpam returns the Ipam field if non-nil, zero value otherwise.

### GetIpamOk

`func (o *DhcpAPIScopesPutInput) GetIpamOk() (*map[string]string, bool)`

GetIpamOk returns a tuple with the Ipam field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIpam

`func (o *DhcpAPIScopesPutInput) SetIpam(v map[string]string)`

SetIpam sets Ipam field to given value.

### HasIpam

`func (o *DhcpAPIScopesPutInput) HasIpam() bool`

HasIpam returns a boolean if a field has been set.

### SetIpamNil

`func (o *DhcpAPIScopesPutInput) SetIpamNil(b bool)`

 SetIpamNil sets the value for Ipam to be an explicit nil

### UnsetIpam
`func (o *DhcpAPIScopesPutInput) UnsetIpam()`

UnsetIpam ensures that no value is present for Ipam, not even an explicit nil
### GetOptions

`func (o *DhcpAPIScopesPutInput) GetOptions() []TypesOption`

GetOptions returns the Options field if non-nil, zero value otherwise.

### GetOptionsOk

`func (o *DhcpAPIScopesPutInput) GetOptionsOk() (*[]TypesOption, bool)`

GetOptionsOk returns a tuple with the Options field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOptions

`func (o *DhcpAPIScopesPutInput) SetOptions(v []TypesOption)`

SetOptions sets Options field to given value.


### SetOptionsNil

`func (o *DhcpAPIScopesPutInput) SetOptionsNil(b bool)`

 SetOptionsNil sets the value for Options to be an explicit nil

### UnsetOptions
`func (o *DhcpAPIScopesPutInput) UnsetOptions()`

UnsetOptions ensures that no value is present for Options, not even an explicit nil
### GetSubnetCidr

`func (o *DhcpAPIScopesPutInput) GetSubnetCidr() string`

GetSubnetCidr returns the SubnetCidr field if non-nil, zero value otherwise.

### GetSubnetCidrOk

`func (o *DhcpAPIScopesPutInput) GetSubnetCidrOk() (*string, bool)`

GetSubnetCidrOk returns a tuple with the SubnetCidr field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubnetCidr

`func (o *DhcpAPIScopesPutInput) SetSubnetCidr(v string)`

SetSubnetCidr sets SubnetCidr field to given value.


### GetTtl

`func (o *DhcpAPIScopesPutInput) GetTtl() int32`

GetTtl returns the Ttl field if non-nil, zero value otherwise.

### GetTtlOk

`func (o *DhcpAPIScopesPutInput) GetTtlOk() (*int32, bool)`

GetTtlOk returns a tuple with the Ttl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTtl

`func (o *DhcpAPIScopesPutInput) SetTtl(v int32)`

SetTtl sets Ttl field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


