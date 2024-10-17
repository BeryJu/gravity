# TypesAPIMetricsRecord

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Keys** | **[]string** |  | 
**Node** | **string** |  | 
**Time** | **time.Time** |  | 
**Value** | **int32** |  | 

## Methods

### NewTypesAPIMetricsRecord

`func NewTypesAPIMetricsRecord(keys []string, node string, time time.Time, value int32, ) *TypesAPIMetricsRecord`

NewTypesAPIMetricsRecord instantiates a new TypesAPIMetricsRecord object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTypesAPIMetricsRecordWithDefaults

`func NewTypesAPIMetricsRecordWithDefaults() *TypesAPIMetricsRecord`

NewTypesAPIMetricsRecordWithDefaults instantiates a new TypesAPIMetricsRecord object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetKeys

`func (o *TypesAPIMetricsRecord) GetKeys() []string`

GetKeys returns the Keys field if non-nil, zero value otherwise.

### GetKeysOk

`func (o *TypesAPIMetricsRecord) GetKeysOk() (*[]string, bool)`

GetKeysOk returns a tuple with the Keys field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetKeys

`func (o *TypesAPIMetricsRecord) SetKeys(v []string)`

SetKeys sets Keys field to given value.


### SetKeysNil

`func (o *TypesAPIMetricsRecord) SetKeysNil(b bool)`

 SetKeysNil sets the value for Keys to be an explicit nil

### UnsetKeys
`func (o *TypesAPIMetricsRecord) UnsetKeys()`

UnsetKeys ensures that no value is present for Keys, not even an explicit nil
### GetNode

`func (o *TypesAPIMetricsRecord) GetNode() string`

GetNode returns the Node field if non-nil, zero value otherwise.

### GetNodeOk

`func (o *TypesAPIMetricsRecord) GetNodeOk() (*string, bool)`

GetNodeOk returns a tuple with the Node field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNode

`func (o *TypesAPIMetricsRecord) SetNode(v string)`

SetNode sets Node field to given value.


### GetTime

`func (o *TypesAPIMetricsRecord) GetTime() time.Time`

GetTime returns the Time field if non-nil, zero value otherwise.

### GetTimeOk

`func (o *TypesAPIMetricsRecord) GetTimeOk() (*time.Time, bool)`

GetTimeOk returns a tuple with the Time field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTime

`func (o *TypesAPIMetricsRecord) SetTime(v time.Time)`

SetTime sets Time field to given value.


### GetValue

`func (o *TypesAPIMetricsRecord) GetValue() int32`

GetValue returns the Value field if non-nil, zero value otherwise.

### GetValueOk

`func (o *TypesAPIMetricsRecord) GetValueOk() (*int32, bool)`

GetValueOk returns a tuple with the Value field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetValue

`func (o *TypesAPIMetricsRecord) SetValue(v int32)`

SetValue sets Value field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


