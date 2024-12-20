# AuthAPIMeOutput

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Authenticated** | **bool** |  | 
**Permissions** | [**[]AuthPermission**](AuthPermission.md) |  | 
**Username** | **string** |  | 

## Methods

### NewAuthAPIMeOutput

`func NewAuthAPIMeOutput(authenticated bool, permissions []AuthPermission, username string, ) *AuthAPIMeOutput`

NewAuthAPIMeOutput instantiates a new AuthAPIMeOutput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAuthAPIMeOutputWithDefaults

`func NewAuthAPIMeOutputWithDefaults() *AuthAPIMeOutput`

NewAuthAPIMeOutputWithDefaults instantiates a new AuthAPIMeOutput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAuthenticated

`func (o *AuthAPIMeOutput) GetAuthenticated() bool`

GetAuthenticated returns the Authenticated field if non-nil, zero value otherwise.

### GetAuthenticatedOk

`func (o *AuthAPIMeOutput) GetAuthenticatedOk() (*bool, bool)`

GetAuthenticatedOk returns a tuple with the Authenticated field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthenticated

`func (o *AuthAPIMeOutput) SetAuthenticated(v bool)`

SetAuthenticated sets Authenticated field to given value.


### GetPermissions

`func (o *AuthAPIMeOutput) GetPermissions() []AuthPermission`

GetPermissions returns the Permissions field if non-nil, zero value otherwise.

### GetPermissionsOk

`func (o *AuthAPIMeOutput) GetPermissionsOk() (*[]AuthPermission, bool)`

GetPermissionsOk returns a tuple with the Permissions field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPermissions

`func (o *AuthAPIMeOutput) SetPermissions(v []AuthPermission)`

SetPermissions sets Permissions field to given value.


### SetPermissionsNil

`func (o *AuthAPIMeOutput) SetPermissionsNil(b bool)`

 SetPermissionsNil sets the value for Permissions to be an explicit nil

### UnsetPermissions
`func (o *AuthAPIMeOutput) UnsetPermissions()`

UnsetPermissions ensures that no value is present for Permissions, not even an explicit nil
### GetUsername

`func (o *AuthAPIMeOutput) GetUsername() string`

GetUsername returns the Username field if non-nil, zero value otherwise.

### GetUsernameOk

`func (o *AuthAPIMeOutput) GetUsernameOk() (*string, bool)`

GetUsernameOk returns a tuple with the Username field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUsername

`func (o *AuthAPIMeOutput) SetUsername(v string)`

SetUsername sets Username field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


