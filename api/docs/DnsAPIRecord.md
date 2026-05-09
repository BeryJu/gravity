# DnsAPIRecord

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Data** | **string** |  | 
**Fqdn** | **string** |  | 
**Hostname** | **string** |  | 
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
**Uid** | **string** |  | 

## Methods

### NewDnsAPIRecord

`func NewDnsAPIRecord(data string, fqdn string, hostname string, ttl int64, type_ TypesDNSRecordType, uid string, ) *DnsAPIRecord`

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

### GetSoaExpire

`func (o *DnsAPIRecord) GetSoaExpire() int32`

GetSoaExpire returns the SoaExpire field if non-nil, zero value otherwise.

### GetSoaExpireOk

`func (o *DnsAPIRecord) GetSoaExpireOk() (*int32, bool)`

GetSoaExpireOk returns a tuple with the SoaExpire field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSoaExpire

`func (o *DnsAPIRecord) SetSoaExpire(v int32)`

SetSoaExpire sets SoaExpire field to given value.

### HasSoaExpire

`func (o *DnsAPIRecord) HasSoaExpire() bool`

HasSoaExpire returns a boolean if a field has been set.

### GetSoaMbox

`func (o *DnsAPIRecord) GetSoaMbox() string`

GetSoaMbox returns the SoaMbox field if non-nil, zero value otherwise.

### GetSoaMboxOk

`func (o *DnsAPIRecord) GetSoaMboxOk() (*string, bool)`

GetSoaMboxOk returns a tuple with the SoaMbox field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSoaMbox

`func (o *DnsAPIRecord) SetSoaMbox(v string)`

SetSoaMbox sets SoaMbox field to given value.

### HasSoaMbox

`func (o *DnsAPIRecord) HasSoaMbox() bool`

HasSoaMbox returns a boolean if a field has been set.

### GetSoaRefresh

`func (o *DnsAPIRecord) GetSoaRefresh() int32`

GetSoaRefresh returns the SoaRefresh field if non-nil, zero value otherwise.

### GetSoaRefreshOk

`func (o *DnsAPIRecord) GetSoaRefreshOk() (*int32, bool)`

GetSoaRefreshOk returns a tuple with the SoaRefresh field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSoaRefresh

`func (o *DnsAPIRecord) SetSoaRefresh(v int32)`

SetSoaRefresh sets SoaRefresh field to given value.

### HasSoaRefresh

`func (o *DnsAPIRecord) HasSoaRefresh() bool`

HasSoaRefresh returns a boolean if a field has been set.

### GetSoaRetry

`func (o *DnsAPIRecord) GetSoaRetry() int32`

GetSoaRetry returns the SoaRetry field if non-nil, zero value otherwise.

### GetSoaRetryOk

`func (o *DnsAPIRecord) GetSoaRetryOk() (*int32, bool)`

GetSoaRetryOk returns a tuple with the SoaRetry field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSoaRetry

`func (o *DnsAPIRecord) SetSoaRetry(v int32)`

SetSoaRetry sets SoaRetry field to given value.

### HasSoaRetry

`func (o *DnsAPIRecord) HasSoaRetry() bool`

HasSoaRetry returns a boolean if a field has been set.

### GetSoaSerial

`func (o *DnsAPIRecord) GetSoaSerial() int32`

GetSoaSerial returns the SoaSerial field if non-nil, zero value otherwise.

### GetSoaSerialOk

`func (o *DnsAPIRecord) GetSoaSerialOk() (*int32, bool)`

GetSoaSerialOk returns a tuple with the SoaSerial field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSoaSerial

`func (o *DnsAPIRecord) SetSoaSerial(v int32)`

SetSoaSerial sets SoaSerial field to given value.

### HasSoaSerial

`func (o *DnsAPIRecord) HasSoaSerial() bool`

HasSoaSerial returns a boolean if a field has been set.

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

### GetTtl

`func (o *DnsAPIRecord) GetTtl() int64`

GetTtl returns the Ttl field if non-nil, zero value otherwise.

### GetTtlOk

`func (o *DnsAPIRecord) GetTtlOk() (*int64, bool)`

GetTtlOk returns a tuple with the Ttl field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTtl

`func (o *DnsAPIRecord) SetTtl(v int64)`

SetTtl sets Ttl field to given value.


### GetType

`func (o *DnsAPIRecord) GetType() TypesDNSRecordType`

GetType returns the Type field if non-nil, zero value otherwise.

### GetTypeOk

`func (o *DnsAPIRecord) GetTypeOk() (*TypesDNSRecordType, bool)`

GetTypeOk returns a tuple with the Type field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetType

`func (o *DnsAPIRecord) SetType(v TypesDNSRecordType)`

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


