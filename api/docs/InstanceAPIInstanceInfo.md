# InstanceAPIInstanceInfo

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BuildHash** | **string** |  | 
**CurrentInstanceIP** | **string** |  | 
**CurrentInstanceIdentifier** | **string** |  | 
**Dirs** | [**ExtconfigExtConfigDirs**](ExtconfigExtConfigDirs.md) |  | 
**Version** | **string** |  | 

## Methods

### NewInstanceAPIInstanceInfo

`func NewInstanceAPIInstanceInfo(buildHash string, currentInstanceIP string, currentInstanceIdentifier string, dirs ExtconfigExtConfigDirs, version string, ) *InstanceAPIInstanceInfo`

NewInstanceAPIInstanceInfo instantiates a new InstanceAPIInstanceInfo object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInstanceAPIInstanceInfoWithDefaults

`func NewInstanceAPIInstanceInfoWithDefaults() *InstanceAPIInstanceInfo`

NewInstanceAPIInstanceInfoWithDefaults instantiates a new InstanceAPIInstanceInfo object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBuildHash

`func (o *InstanceAPIInstanceInfo) GetBuildHash() string`

GetBuildHash returns the BuildHash field if non-nil, zero value otherwise.

### GetBuildHashOk

`func (o *InstanceAPIInstanceInfo) GetBuildHashOk() (*string, bool)`

GetBuildHashOk returns a tuple with the BuildHash field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBuildHash

`func (o *InstanceAPIInstanceInfo) SetBuildHash(v string)`

SetBuildHash sets BuildHash field to given value.


### GetCurrentInstanceIP

`func (o *InstanceAPIInstanceInfo) GetCurrentInstanceIP() string`

GetCurrentInstanceIP returns the CurrentInstanceIP field if non-nil, zero value otherwise.

### GetCurrentInstanceIPOk

`func (o *InstanceAPIInstanceInfo) GetCurrentInstanceIPOk() (*string, bool)`

GetCurrentInstanceIPOk returns a tuple with the CurrentInstanceIP field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrentInstanceIP

`func (o *InstanceAPIInstanceInfo) SetCurrentInstanceIP(v string)`

SetCurrentInstanceIP sets CurrentInstanceIP field to given value.


### GetCurrentInstanceIdentifier

`func (o *InstanceAPIInstanceInfo) GetCurrentInstanceIdentifier() string`

GetCurrentInstanceIdentifier returns the CurrentInstanceIdentifier field if non-nil, zero value otherwise.

### GetCurrentInstanceIdentifierOk

`func (o *InstanceAPIInstanceInfo) GetCurrentInstanceIdentifierOk() (*string, bool)`

GetCurrentInstanceIdentifierOk returns a tuple with the CurrentInstanceIdentifier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCurrentInstanceIdentifier

`func (o *InstanceAPIInstanceInfo) SetCurrentInstanceIdentifier(v string)`

SetCurrentInstanceIdentifier sets CurrentInstanceIdentifier field to given value.


### GetDirs

`func (o *InstanceAPIInstanceInfo) GetDirs() ExtconfigExtConfigDirs`

GetDirs returns the Dirs field if non-nil, zero value otherwise.

### GetDirsOk

`func (o *InstanceAPIInstanceInfo) GetDirsOk() (*ExtconfigExtConfigDirs, bool)`

GetDirsOk returns a tuple with the Dirs field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDirs

`func (o *InstanceAPIInstanceInfo) SetDirs(v ExtconfigExtConfigDirs)`

SetDirs sets Dirs field to given value.


### GetVersion

`func (o *InstanceAPIInstanceInfo) GetVersion() string`

GetVersion returns the Version field if non-nil, zero value otherwise.

### GetVersionOk

`func (o *InstanceAPIInstanceInfo) GetVersionOk() (*string, bool)`

GetVersionOk returns a tuple with the Version field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetVersion

`func (o *InstanceAPIInstanceInfo) SetVersion(v string)`

SetVersion sets Version field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


