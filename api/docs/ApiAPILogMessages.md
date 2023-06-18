# ApiAPILogMessages

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**IsJSON** | Pointer to **bool** |  | [optional] 
**Messages** | Pointer to [**[]ApiAPILogMessage**](ApiAPILogMessage.md) |  | [optional] 

## Methods

### NewApiAPILogMessages

`func NewApiAPILogMessages() *ApiAPILogMessages`

NewApiAPILogMessages instantiates a new ApiAPILogMessages object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewApiAPILogMessagesWithDefaults

`func NewApiAPILogMessagesWithDefaults() *ApiAPILogMessages`

NewApiAPILogMessagesWithDefaults instantiates a new ApiAPILogMessages object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetIsJSON

`func (o *ApiAPILogMessages) GetIsJSON() bool`

GetIsJSON returns the IsJSON field if non-nil, zero value otherwise.

### GetIsJSONOk

`func (o *ApiAPILogMessages) GetIsJSONOk() (*bool, bool)`

GetIsJSONOk returns a tuple with the IsJSON field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsJSON

`func (o *ApiAPILogMessages) SetIsJSON(v bool)`

SetIsJSON sets IsJSON field to given value.

### HasIsJSON

`func (o *ApiAPILogMessages) HasIsJSON() bool`

HasIsJSON returns a boolean if a field has been set.

### GetMessages

`func (o *ApiAPILogMessages) GetMessages() []ApiAPILogMessage`

GetMessages returns the Messages field if non-nil, zero value otherwise.

### GetMessagesOk

`func (o *ApiAPILogMessages) GetMessagesOk() (*[]ApiAPILogMessage, bool)`

GetMessagesOk returns a tuple with the Messages field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessages

`func (o *ApiAPILogMessages) SetMessages(v []ApiAPILogMessage)`

SetMessages sets Messages field to given value.

### HasMessages

`func (o *ApiAPILogMessages) HasMessages() bool`

HasMessages returns a boolean if a field has been set.

### SetMessagesNil

`func (o *ApiAPILogMessages) SetMessagesNil(b bool)`

 SetMessagesNil sets the value for Messages to be an explicit nil

### UnsetMessages
`func (o *ApiAPILogMessages) UnsetMessages()`

UnsetMessages ensures that no value is present for Messages, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


