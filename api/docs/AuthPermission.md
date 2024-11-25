# AuthPermission

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Methods** | Pointer to **[]string** |  | [optional] 
**Path** | Pointer to **string** |  | [optional] 

## Methods

### NewAuthPermission

`func NewAuthPermission() *AuthPermission`

NewAuthPermission instantiates a new AuthPermission object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewAuthPermissionWithDefaults

`func NewAuthPermissionWithDefaults() *AuthPermission`

NewAuthPermissionWithDefaults instantiates a new AuthPermission object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMethods

`func (o *AuthPermission) GetMethods() []string`

GetMethods returns the Methods field if non-nil, zero value otherwise.

### GetMethodsOk

`func (o *AuthPermission) GetMethodsOk() (*[]string, bool)`

GetMethodsOk returns a tuple with the Methods field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMethods

`func (o *AuthPermission) SetMethods(v []string)`

SetMethods sets Methods field to given value.

### HasMethods

`func (o *AuthPermission) HasMethods() bool`

HasMethods returns a boolean if a field has been set.

### SetMethodsNil

`func (o *AuthPermission) SetMethodsNil(b bool)`

 SetMethodsNil sets the value for Methods to be an explicit nil

### UnsetMethods
`func (o *AuthPermission) UnsetMethods()`

UnsetMethods ensures that no value is present for Methods, not even an explicit nil
### GetPath

`func (o *AuthPermission) GetPath() string`

GetPath returns the Path field if non-nil, zero value otherwise.

### GetPathOk

`func (o *AuthPermission) GetPathOk() (*string, bool)`

GetPathOk returns a tuple with the Path field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPath

`func (o *AuthPermission) SetPath(v string)`

SetPath sets Path field to given value.

### HasPath

`func (o *AuthPermission) HasPath() bool`

HasPath returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


