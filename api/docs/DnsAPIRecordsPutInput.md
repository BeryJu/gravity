# DnsAPIRecordsPutInput

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | **string** |  | 
**MxPreference** | Pointer to **int32** |  | [optional] 
**SoaExpire** | Pointer to **int32** |  | [optional] 
**SoaMbox** | Pointer to **string** |  | [optional] 
**SoaRefresh** | Pointer to **int32** |  | [optional] 
**SoaRetry** | Pointer to **int32** |  | [optional] 
**SoaSerial** | Pointer to **int32** |  | [optional] 
**SrvPort** | Pointer to **int32** |  | [optional] 
**SrvPriority** | Pointer to **int32** |  | [optional] 
**SrvWeight** | Pointer to **int32** |  | [optional] 
**Ttl** | **int64** |  | 
**Type** | [**TypesDNSRecordType**](TypesDNSRecordType.md) |  | 

## Methods

### NewDnsAPIRecordsPutInput

`func NewDnsAPIRecordsPutInput(data string, ttl int64, type_ TypesDNSRecordType, ) *DnsAPIRecordsPutInput`

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

### GetSoaExpire

`func (o *DnsAPIRecordsPutInput) GetSoaExpire() int32`

GetSoaExpire returns the SoaExpire field if non-nil, zero value otherwise.

### GetSoaExpireOk

`func (o *DnsAPIRecordsPutInput) GetSoaExpireOk() (*int32, bool)`

GetSoaExpireOk returns a tuple with the SoaExpire field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSoaExpire

`func (o *DnsAPIRecordsPutInput) SetSoaExpire(v int32)`

SetSoaExpire sets SoaExpire field to given value.

### HasSoaExpire

`func (o *DnsAPIRecordsPutInput) HasSoaExpire() bool`

HasSoaExpire returns a boolean if a field has been set.

### GetSoaMbox

`func (o *DnsAPIRecordsPutInput) GetSoaMbox() string`

GetSoaMbox returns the SoaMbox field if non-nil, zero value otherwise.

### GetSoaMboxOk

`func (o *DnsAPIRecordsPutInput) GetSoaMboxOk() (*string, bool)`

GetSoaMboxOk returns a tuple with the SoaMbox field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSoaMbox

`func (o *DnsAPIRecordsPutInput) SetSoaMbox(v string)`

SetSoaMbox sets SoaMbox field to given value.

### HasSoaMbox

`func (o *DnsAPIRecordsPutInput) HasSoaMbox() bool`

HasSoaMbox returns a boolean if a field has been set.

### GetSoaRefresh

`func (o *DnsAPIRecordsPutInput) GetSoaRefresh() int32`

GetSoaRefresh returns the SoaRefresh field if non-nil, zero value otherwise.

### GetSoaRefreshOk

`func (o *DnsAPIRecordsPutInput) GetSoaRefreshOk() (*int32, bool)`

GetSoaRefreshOk returns a tuple with the SoaRefresh field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSoaRefresh

`func (o *DnsAPIRecordsPutInput) SetSoaRefresh(v int32)`

SetSoaRefresh sets SoaRefresh field to given value.

### HasSoaRefresh

`func (o *DnsAPIRecordsPutInput) HasSoaRefresh() bool`

HasSoaRefresh returns a boolean if a field has been set.

### GetSoaRetry

`func (o *DnsAPIRecordsPutInput) GetSoaRetry() int32`

GetSoaRetry returns the SoaRetry field if non-nil, zero value otherwise.

### GetSoaRetryOk

`func (o *DnsAPIRecordsPutInput) GetSoaRetryOk() (*int32, bool)`

GetSoaRetryOk returns a tuple with the SoaRetry field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSoaRetry

`func (o *DnsAPIRecordsPutInput) SetSoaRetry(v int32)`

SetSoaRetry sets SoaRetry field to given value.

### HasSoaRetry

`func (o *DnsAPIRecordsPutInput) HasSoaRetry() bool`

HasSoaRetry returns a boolean if a field has been set.

### GetSoaSerial

`func (o *DnsAPIRecordsPutInput) GetSoaSerial() int32`

GetSoaSerial returns the SoaSerial field if non-nil, zero value otherwise.

### GetSoaSerialOk

`func (o *DnsAPIRecordsPutInput) GetSoaSerialOk() (*int32, bool)`

GetSoaSerialOk returns a tuple with the SoaSerial field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSoaSerial

`func (o *DnsAPIRecordsPutInput) SetSoaSerial(v int32)`

SetSoaSerial sets SoaSerial field to given value.

### HasSoaSerial

`func (o *DnsAPIRecordsPutInput) HasSoaSerial() bool`

HasSoaSerial returns a boolean if a field has been set.

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

### GetTtl

`func (o *DnsAPIRecordsPutInput) GetTtl() int64`

GetTtl returns the Ttl field if non-nil, zero value otherwise.

### GetTtlOk

`func (o *DnsAPIRecordsPutInput) GetTtlOk() (*int64, bool)`

GetTtlOk returns a tuple with the Ttl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTtl

`func (o *DnsAPIRecordsPutInput) SetTtl(v int64)`

SetTtl sets Ttl field to given value.


### GetType

`func (o *DnsAPIRecordsPutInput) GetType() TypesDNSRecordType`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *DnsAPIRecordsPutInput) GetTypeOk() (*TypesDNSRecordType, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *DnsAPIRecordsPutInput) SetType(v TypesDNSRecordType)`

SetType sets Type field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


