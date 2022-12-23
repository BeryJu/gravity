# BackupAPIBackupStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Duration** | Pointer to **int32** |  | [optional] 
**Error** | Pointer to **string** |  | [optional] 
**Filename** | Pointer to **string** |  | [optional] 
**Node** | Pointer to **string** |  | [optional] 
**Size** | Pointer to **int32** |  | [optional] 
**Status** | Pointer to **string** |  | [optional] 
**Time** | Pointer to **time.Time** |  | [optional] 

## Methods

### NewBackupAPIBackupStatus

`func NewBackupAPIBackupStatus() *BackupAPIBackupStatus`

NewBackupAPIBackupStatus instantiates a new BackupAPIBackupStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBackupAPIBackupStatusWithDefaults

`func NewBackupAPIBackupStatusWithDefaults() *BackupAPIBackupStatus`

NewBackupAPIBackupStatusWithDefaults instantiates a new BackupAPIBackupStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDuration

`func (o *BackupAPIBackupStatus) GetDuration() int32`

GetDuration returns the Duration field if non-nil, zero value otherwise.

### GetDurationOk

`func (o *BackupAPIBackupStatus) GetDurationOk() (*int32, bool)`

GetDurationOk returns a tuple with the Duration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDuration

`func (o *BackupAPIBackupStatus) SetDuration(v int32)`

SetDuration sets Duration field to given value.

### HasDuration

`func (o *BackupAPIBackupStatus) HasDuration() bool`

HasDuration returns a boolean if a field has been set.

### GetError

`func (o *BackupAPIBackupStatus) GetError() string`

GetError returns the Error field if non-nil, zero value otherwise.

### GetErrorOk

`func (o *BackupAPIBackupStatus) GetErrorOk() (*string, bool)`

GetErrorOk returns a tuple with the Error field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetError

`func (o *BackupAPIBackupStatus) SetError(v string)`

SetError sets Error field to given value.

### HasError

`func (o *BackupAPIBackupStatus) HasError() bool`

HasError returns a boolean if a field has been set.

### GetFilename

`func (o *BackupAPIBackupStatus) GetFilename() string`

GetFilename returns the Filename field if non-nil, zero value otherwise.

### GetFilenameOk

`func (o *BackupAPIBackupStatus) GetFilenameOk() (*string, bool)`

GetFilenameOk returns a tuple with the Filename field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFilename

`func (o *BackupAPIBackupStatus) SetFilename(v string)`

SetFilename sets Filename field to given value.

### HasFilename

`func (o *BackupAPIBackupStatus) HasFilename() bool`

HasFilename returns a boolean if a field has been set.

### GetNode

`func (o *BackupAPIBackupStatus) GetNode() string`

GetNode returns the Node field if non-nil, zero value otherwise.

### GetNodeOk

`func (o *BackupAPIBackupStatus) GetNodeOk() (*string, bool)`

GetNodeOk returns a tuple with the Node field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNode

`func (o *BackupAPIBackupStatus) SetNode(v string)`

SetNode sets Node field to given value.

### HasNode

`func (o *BackupAPIBackupStatus) HasNode() bool`

HasNode returns a boolean if a field has been set.

### GetSize

`func (o *BackupAPIBackupStatus) GetSize() int32`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *BackupAPIBackupStatus) GetSizeOk() (*int32, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *BackupAPIBackupStatus) SetSize(v int32)`

SetSize sets Size field to given value.

### HasSize

`func (o *BackupAPIBackupStatus) HasSize() bool`

HasSize returns a boolean if a field has been set.

### GetStatus

`func (o *BackupAPIBackupStatus) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *BackupAPIBackupStatus) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *BackupAPIBackupStatus) SetStatus(v string)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *BackupAPIBackupStatus) HasStatus() bool`

HasStatus returns a boolean if a field has been set.

### GetTime

`func (o *BackupAPIBackupStatus) GetTime() time.Time`

GetTime returns the Time field if non-nil, zero value otherwise.

### GetTimeOk

`func (o *BackupAPIBackupStatus) GetTimeOk() (*time.Time, bool)`

GetTimeOk returns a tuple with the Time field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTime

`func (o *BackupAPIBackupStatus) SetTime(v time.Time)`

SetTime sets Time field to given value.

### HasTime

`func (o *BackupAPIBackupStatus) HasTime() bool`

HasTime returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


