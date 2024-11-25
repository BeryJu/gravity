# AuthAPIUser

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Permissions** | [**[]AuthPermission**](AuthPermission.md) |  | 
**Username** | **string** |  | 

## Methods

### NewAuthAPIUser

`func NewAuthAPIUser(permissions []AuthPermission, username string, ) *AuthAPIUser`

NewAuthAPIUser instantiates a new AuthAPIUser object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAuthAPIUserWithDefaults

`func NewAuthAPIUserWithDefaults() *AuthAPIUser`

NewAuthAPIUserWithDefaults instantiates a new AuthAPIUser object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPermissions

`func (o *AuthAPIUser) GetPermissions() []AuthPermission`

GetPermissions returns the Permissions field if non-nil, zero value otherwise.

### GetPermissionsOk

`func (o *AuthAPIUser) GetPermissionsOk() (*[]AuthPermission, bool)`

GetPermissionsOk returns a tuple with the Permissions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPermissions

`func (o *AuthAPIUser) SetPermissions(v []AuthPermission)`

SetPermissions sets Permissions field to given value.


### SetPermissionsNil

`func (o *AuthAPIUser) SetPermissionsNil(b bool)`

 SetPermissionsNil sets the value for Permissions to be an explicit nil

### UnsetPermissions
`func (o *AuthAPIUser) UnsetPermissions()`

UnsetPermissions ensures that no value is present for Permissions, not even an explicit nil
### GetUsername

`func (o *AuthAPIUser) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *AuthAPIUser) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *AuthAPIUser) SetUsername(v string)`

SetUsername sets Username field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


