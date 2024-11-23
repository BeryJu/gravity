# ApiAPILogMessage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Level** | Pointer to **string** |  | [optional] 
**Logger** | Pointer to **string** |  | [optional] 
**Message** | Pointer to **string** |  | [optional] 
**Node** | Pointer to **string** |  | [optional] 
**Time** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewApiAPILogMessage

`func NewApiAPILogMessage() *ApiAPILogMessage`

NewApiAPILogMessage instantiates a new ApiAPILogMessage object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewApiAPILogMessageWithDefaults

`func NewApiAPILogMessageWithDefaults() *ApiAPILogMessage`

NewApiAPILogMessageWithDefaults instantiates a new ApiAPILogMessage object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLevel

`func (o *ApiAPILogMessage) GetLevel() string`

GetLevel returns the Level field if non-nil, zero value otherwise.

### GetLevelOk

`func (o *ApiAPILogMessage) GetLevelOk() (*string, bool)`

GetLevelOk returns a tuple with the Level field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLevel

`func (o *ApiAPILogMessage) SetLevel(v string)`

SetLevel sets Level field to given value.

### HasLevel

`func (o *ApiAPILogMessage) HasLevel() bool`

HasLevel returns a boolean if a field has been set.

### GetLogger

`func (o *ApiAPILogMessage) GetLogger() string`

GetLogger returns the Logger field if non-nil, zero value otherwise.

### GetLoggerOk

`func (o *ApiAPILogMessage) GetLoggerOk() (*string, bool)`

GetLoggerOk returns a tuple with the Logger field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLogger

`func (o *ApiAPILogMessage) SetLogger(v string)`

SetLogger sets Logger field to given value.

### HasLogger

`func (o *ApiAPILogMessage) HasLogger() bool`

HasLogger returns a boolean if a field has been set.

### GetMessage

`func (o *ApiAPILogMessage) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *ApiAPILogMessage) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *ApiAPILogMessage) SetMessage(v string)`

SetMessage sets Message field to given value.

### HasMessage

`func (o *ApiAPILogMessage) HasMessage() bool`

HasMessage returns a boolean if a field has been set.

### GetNode

`func (o *ApiAPILogMessage) GetNode() string`

GetNode returns the Node field if non-nil, zero value otherwise.

### GetNodeOk

`func (o *ApiAPILogMessage) GetNodeOk() (*string, bool)`

GetNodeOk returns a tuple with the Node field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNode

`func (o *ApiAPILogMessage) SetNode(v string)`

SetNode sets Node field to given value.

### HasNode

`func (o *ApiAPILogMessage) HasNode() bool`

HasNode returns a boolean if a field has been set.

### GetTime

`func (o *ApiAPILogMessage) GetTime() time.Time`

GetTime returns the Time field if non-nil, zero value otherwise.

### GetTimeOk

`func (o *ApiAPILogMessage) GetTimeOk() (*time.Time, bool)`

GetTimeOk returns a tuple with the Time field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTime

`func (o *ApiAPILogMessage) SetTime(v time.Time)`

SetTime sets Time field to given value.

### HasTime

`func (o *ApiAPILogMessage) HasTime() bool`

HasTime returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


