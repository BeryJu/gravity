# DnsAPIRecordsPutInput

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | **string** |  | 
**MxPreference** | Pointer to **int32** |  | [optional] 
**SrvPort** | Pointer to **int32** |  | [optional] 
**SrvPriority** | Pointer to **int32** |  | [optional] 
**SrvWeight** | Pointer to **int32** |  | [optional] 
**Type** | **string** |  | 

## Methods

### NewDnsAPIRecordsPutInput

`func NewDnsAPIRecordsPutInput(data string, type_ string, ) *DnsAPIRecordsPutInput`

NewDnsAPIRecordsPutInput instantiates a new DnsAPIRecordsPutInput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDnsAPIRecordsPutInputWithDefaults

`func NewDnsAPIRecordsPutInputWithDefaults() *DnsAPIRecordsPutInput`

NewDnsAPIRecordsPutInputWithDefaults instantiates a new DnsAPIRecordsPutInput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetData

`func (o *DnsAPIRecordsPutInput) GetData() string`

GetData returns the Data field if non-nil, zero value otherwise.

### GetDataOk

`func (o *DnsAPIRecordsPutInput) GetDataOk() (*string, bool)`

GetDataOk returns a tuple with the Data field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetData

`func (o *DnsAPIRecordsPutInput) SetData(v string)`

SetData sets Data field to given value.


### GetMxPreference

`func (o *DnsAPIRecordsPutInput) GetMxPreference() int32`

GetMxPreference returns the MxPreference field if non-nil, zero value otherwise.

### GetMxPreferenceOk

`func (o *DnsAPIRecordsPutInput) GetMxPreferenceOk() (*int32, bool)`

GetMxPreferenceOk returns a tuple with the MxPreference field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMxPreference

`func (o *DnsAPIRecordsPutInput) SetMxPreference(v int32)`

SetMxPreference sets MxPreference field to given value.

### HasMxPreference

`func (o *DnsAPIRecordsPutInput) HasMxPreference() bool`

HasMxPreference returns a boolean if a field has been set.

### GetSrvPort

`func (o *DnsAPIRecordsPutInput) GetSrvPort() int32`

GetSrvPort returns the SrvPort field if non-nil, zero value otherwise.

### GetSrvPortOk

`func (o *DnsAPIRecordsPutInput) GetSrvPortOk() (*int32, bool)`

GetSrvPortOk returns a tuple with the SrvPort field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSrvPort

`func (o *DnsAPIRecordsPutInput) SetSrvPort(v int32)`

SetSrvPort sets SrvPort field to given value.

### HasSrvPort

`func (o *DnsAPIRecordsPutInput) HasSrvPort() bool`

HasSrvPort returns a boolean if a field has been set.

### GetSrvPriority

`func (o *DnsAPIRecordsPutInput) GetSrvPriority() int32`

GetSrvPriority returns the SrvPriority field if non-nil, zero value otherwise.

### GetSrvPriorityOk

`func (o *DnsAPIRecordsPutInput) GetSrvPriorityOk() (*int32, bool)`

GetSrvPriorityOk returns a tuple with the SrvPriority field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSrvPriority

`func (o *DnsAPIRecordsPutInput) SetSrvPriority(v int32)`

SetSrvPriority sets SrvPriority field to given value.

### HasSrvPriority

`func (o *DnsAPIRecordsPutInput) HasSrvPriority() bool`

HasSrvPriority returns a boolean if a field has been set.

### GetSrvWeight

`func (o *DnsAPIRecordsPutInput) GetSrvWeight() int32`

GetSrvWeight returns the SrvWeight field if non-nil, zero value otherwise.

### GetSrvWeightOk

`func (o *DnsAPIRecordsPutInput) GetSrvWeightOk() (*int32, bool)`

GetSrvWeightOk returns a tuple with the SrvWeight field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSrvWeight

`func (o *DnsAPIRecordsPutInput) SetSrvWeight(v int32)`

SetSrvWeight sets SrvWeight field to given value.

### HasSrvWeight

`func (o *DnsAPIRecordsPutInput) HasSrvWeight() bool`

HasSrvWeight returns a boolean if a field has been set.

### GetType

`func (o *DnsAPIRecordsPutInput) GetType() string`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *DnsAPIRecordsPutInput) GetTypeOk() (*string, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *DnsAPIRecordsPutInput) SetType(v string)`

SetType sets Type field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


