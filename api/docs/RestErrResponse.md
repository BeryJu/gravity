# RestErrResponse

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Code** | Pointer to **int32** | Application-specific error code. | [optional] 
**Context** | Pointer to **map[string]interface{}** | Application context. | [optional] 
**Error** | Pointer to **string** | Error message. | [optional] 
**Status** | Pointer to **string** | Status text. | [optional] 

## Methods

### NewRestErrResponse

`func NewRestErrResponse() *RestErrResponse`

NewRestErrResponse instantiates a new RestErrResponse object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewRestErrResponseWithDefaults

`func NewRestErrResponseWithDefaults() *RestErrResponse`

NewRestErrResponseWithDefaults instantiates a new RestErrResponse object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCode

`func (o *RestErrResponse) GetCode() int32`

GetCode returns the Code field if non-nil, zero value otherwise.

### GetCodeOk

`func (o *RestErrResponse) GetCodeOk() (*int32, bool)`

GetCodeOk returns a tuple with the Code field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCode

`func (o *RestErrResponse) SetCode(v int32)`

SetCode sets Code field to given value.

### HasCode

`func (o *RestErrResponse) HasCode() bool`

HasCode returns a boolean if a field has been set.

### GetContext

`func (o *RestErrResponse) GetContext() map[string]interface{}`

GetContext returns the Context field if non-nil, zero value otherwise.

### GetContextOk

`func (o *RestErrResponse) GetContextOk() (*map[string]interface{}, bool)`

GetContextOk returns a tuple with the Context field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetContext

`func (o *RestErrResponse) SetContext(v map[string]interface{})`

SetContext sets Context field to given value.

### HasContext

`func (o *RestErrResponse) HasContext() bool`

HasContext returns a boolean if a field has been set.

### GetError

`func (o *RestErrResponse) GetError() string`

GetError returns the Error field if non-nil, zero value otherwise.

### GetErrorOk

`func (o *RestErrResponse) GetErrorOk() (*string, bool)`

GetErrorOk returns a tuple with the Error field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetError

`func (o *RestErrResponse) SetError(v string)`

SetError sets Error field to given value.

### HasError

`func (o *RestErrResponse) HasError() bool`

HasError returns a boolean if a field has been set.

### GetStatus

`func (o *RestErrResponse) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *RestErrResponse) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *RestErrResponse) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *RestErrResponse) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


