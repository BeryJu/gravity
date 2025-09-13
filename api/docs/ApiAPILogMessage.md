# ApiAPILogMessage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Fields** | **map[string]interface{}** |  | 
**Level** | **string** |  | 
**Logger** | **string** |  | 
**Message** | **string** |  | 
**Node** | **string** |  | 
**Time** | **time.Time** |  | 

## Methods

### NewApiAPILogMessage

`func NewApiAPILogMessage(fields map[string]interface{}, level string, logger string, message string, node string, time time.Time, ) *ApiAPILogMessage`

NewApiAPILogMessage instantiates a new ApiAPILogMessage object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewApiAPILogMessageWithDefaults

`func NewApiAPILogMessageWithDefaults() *ApiAPILogMessage`

NewApiAPILogMessageWithDefaults instantiates a new ApiAPILogMessage object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetFields

`func (o *ApiAPILogMessage) GetFields() map[string]interface{}`

GetFields returns the Fields field if non-nil, zero value otherwise.

### GetFieldsOk

`func (o *ApiAPILogMessage) GetFieldsOk() (*map[string]interface{}, bool)`

GetFieldsOk returns a tuple with the Fields field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFields

`func (o *ApiAPILogMessage) SetFields(v map[string]interface{})`

SetFields sets Fields field to given value.


### SetFieldsNil

`func (o *ApiAPILogMessage) SetFieldsNil(b bool)`

 SetFieldsNil sets the value for Fields to be an explicit nil

### UnsetFields
`func (o *ApiAPILogMessage) UnsetFields()`

UnsetFields ensures that no value is present for Fields, not even an explicit nil
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



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


