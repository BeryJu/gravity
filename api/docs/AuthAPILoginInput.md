# AuthAPILoginInput

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Password** | Pointer to **string** |  | [optional] 
**Username** | Pointer to **string** |  | [optional] 

## Methods

### NewAuthAPILoginInput

`func NewAuthAPILoginInput() *AuthAPILoginInput`

NewAuthAPILoginInput instantiates a new AuthAPILoginInput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAuthAPILoginInputWithDefaults

`func NewAuthAPILoginInputWithDefaults() *AuthAPILoginInput`

NewAuthAPILoginInputWithDefaults instantiates a new AuthAPILoginInput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPassword

`func (o *AuthAPILoginInput) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *AuthAPILoginInput) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *AuthAPILoginInput) SetPassword(v string)`

SetPassword sets Password field to given value.

### HasPassword

`func (o *AuthAPILoginInput) HasPassword() bool`

HasPassword returns a boolean if a field has been set.

### GetUsername

`func (o *AuthAPILoginInput) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *AuthAPILoginInput) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *AuthAPILoginInput) SetUsername(v string)`

SetUsername sets Username field to given value.

### HasUsername

`func (o *AuthAPILoginInput) HasUsername() bool`

HasUsername returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


