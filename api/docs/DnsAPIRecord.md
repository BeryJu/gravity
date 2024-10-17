# DnsAPIRecord

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | **string** |  | 
**Fqdn** | **string** |  | 
**Hostname** | **string** |  | 
**MxPreference** | Pointer to **int32** |  | [optional] 
**SrvPort** | Pointer to **int32** |  | [optional] 
**SrvPriority** | Pointer to **int32** |  | [optional] 
**SrvWeight** | Pointer to **int32** |  | [optional] 
**Type** | **string** |  | 
**Uid** | **string** |  | 

## Methods

### NewDnsAPIRecord

`func NewDnsAPIRecord(data string, fqdn string, hostname string, type_ string, uid string, ) *DnsAPIRecord`

NewDnsAPIRecord instantiates a new DnsAPIRecord object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDnsAPIRecordWithDefaults

`func NewDnsAPIRecordWithDefaults() *DnsAPIRecord`

NewDnsAPIRecordWithDefaults instantiates a new DnsAPIRecord object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *DnsAPIRecord) GetData() string`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *DnsAPIRecord) GetDataOk() (*string, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *DnsAPIRecord) SetData(v string)`

SetData sets Data field to given value.


### GetFqdn

`func (o *DnsAPIRecord) GetFqdn() string`

GetFqdn returns the Fqdn field if non-nil, zero value otherwise.

### GetFqdnOk

`func (o *DnsAPIRecord) GetFqdnOk() (*string, bool)`

GetFqdnOk returns a tuple with the Fqdn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFqdn

`func (o *DnsAPIRecord) SetFqdn(v string)`

SetFqdn sets Fqdn field to given value.


### GetHostname

`func (o *DnsAPIRecord) GetHostname() string`

GetHostname returns the Hostname field if non-nil, zero value otherwise.

### GetHostnameOk

`func (o *DnsAPIRecord) GetHostnameOk() (*string, bool)`

GetHostnameOk returns a tuple with the Hostname field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHostname

`func (o *DnsAPIRecord) SetHostname(v string)`

SetHostname sets Hostname field to given value.


### GetMxPreference

`func (o *DnsAPIRecord) GetMxPreference() int32`

GetMxPreference returns the MxPreference field if non-nil, zero value otherwise.

### GetMxPreferenceOk

`func (o *DnsAPIRecord) GetMxPreferenceOk() (*int32, bool)`

GetMxPreferenceOk returns a tuple with the MxPreference field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMxPreference

`func (o *DnsAPIRecord) SetMxPreference(v int32)`

SetMxPreference sets MxPreference field to given value.

### HasMxPreference

`func (o *DnsAPIRecord) HasMxPreference() bool`

HasMxPreference returns a boolean if a field has been set.

### GetSrvPort

`func (o *DnsAPIRecord) GetSrvPort() int32`

GetSrvPort returns the SrvPort field if non-nil, zero value otherwise.

### GetSrvPortOk

`func (o *DnsAPIRecord) GetSrvPortOk() (*int32, bool)`

GetSrvPortOk returns a tuple with the SrvPort field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSrvPort

`func (o *DnsAPIRecord) SetSrvPort(v int32)`

SetSrvPort sets SrvPort field to given value.

### HasSrvPort

`func (o *DnsAPIRecord) HasSrvPort() bool`

HasSrvPort returns a boolean if a field has been set.

### GetSrvPriority

`func (o *DnsAPIRecord) GetSrvPriority() int32`

GetSrvPriority returns the SrvPriority field if non-nil, zero value otherwise.

### GetSrvPriorityOk

`func (o *DnsAPIRecord) GetSrvPriorityOk() (*int32, bool)`

GetSrvPriorityOk returns a tuple with the SrvPriority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSrvPriority

`func (o *DnsAPIRecord) SetSrvPriority(v int32)`

SetSrvPriority sets SrvPriority field to given value.

### HasSrvPriority

`func (o *DnsAPIRecord) HasSrvPriority() bool`

HasSrvPriority returns a boolean if a field has been set.

### GetSrvWeight

`func (o *DnsAPIRecord) GetSrvWeight() int32`

GetSrvWeight returns the SrvWeight field if non-nil, zero value otherwise.

### GetSrvWeightOk

`func (o *DnsAPIRecord) GetSrvWeightOk() (*int32, bool)`

GetSrvWeightOk returns a tuple with the SrvWeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSrvWeight

`func (o *DnsAPIRecord) SetSrvWeight(v int32)`

SetSrvWeight sets SrvWeight field to given value.

### HasSrvWeight

`func (o *DnsAPIRecord) HasSrvWeight() bool`

HasSrvWeight returns a boolean if a field has been set.

### GetType

`func (o *DnsAPIRecord) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *DnsAPIRecord) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *DnsAPIRecord) SetType(v string)`

SetType sets Type field to given value.


### GetUid

`func (o *DnsAPIRecord) GetUid() string`

GetUid returns the Uid field if non-nil, zero value otherwise.

### GetUidOk

`func (o *DnsAPIRecord) GetUidOk() (*string, bool)`

GetUidOk returns a tuple with the Uid field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUid

`func (o *DnsAPIRecord) SetUid(v string)`

SetUid sets Uid field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


