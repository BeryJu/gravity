# AuthAPIUsersPutInput

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Password** | **string** |  | 
**Permissions** | [**[]AuthPermission**](AuthPermission.md) |  | 

## Methods

### NewAuthAPIUsersPutInput

`func NewAuthAPIUsersPutInput(password string, permissions []AuthPermission, ) *AuthAPIUsersPutInput`

NewAuthAPIUsersPutInput instantiates a new AuthAPIUsersPutInput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAuthAPIUsersPutInputWithDefaults

`func NewAuthAPIUsersPutInputWithDefaults() *AuthAPIUsersPutInput`

NewAuthAPIUsersPutInputWithDefaults instantiates a new AuthAPIUsersPutInput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPassword

`func (o *AuthAPIUsersPutInput) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *AuthAPIUsersPutInput) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *AuthAPIUsersPutInput) SetPassword(v string)`

SetPassword sets Password field to given value.


### GetPermissions

`func (o *AuthAPIUsersPutInput) GetPermissions() []AuthPermission`

GetPermissions returns the Permissions field if non-nil, zero value otherwise.

### GetPermissionsOk

`func (o *AuthAPIUsersPutInput) GetPermissionsOk() (*[]AuthPermission, bool)`

GetPermissionsOk returns a tuple with the Permissions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPermissions

`func (o *AuthAPIUsersPutInput) SetPermissions(v []AuthPermission)`

SetPermissions sets Permissions field to given value.


### SetPermissionsNil

`func (o *AuthAPIUsersPutInput) SetPermissionsNil(b bool)`

 SetPermissionsNil sets the value for Permissions to be an explicit nil

### UnsetPermissions
`func (o *AuthAPIUsersPutInput) UnsetPermissions()`

UnsetPermissions ensures that no value is present for Permissions, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


