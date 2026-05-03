# BackupAPIBackupStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Duration** | **int64** |  | 
**Error** | **string** |  | 
**Filename** | **string** |  | 
**Node** | Pointer to **string** |  | [optional] 
**Size** | **int64** |  | 
**Status** | **string** |  | 
**Time** | **time.Time** |  | 

## Methods

### NewBackupAPIBackupStatus

`func NewBackupAPIBackupStatus(duration int64, error_ string, filename string, size int64, status string, time time.Time, ) *BackupAPIBackupStatus`

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

`func (o *BackupAPIBackupStatus) GetDuration() int64`

GetDuration returns the Duration field if non-nil, zero value otherwise.

### GetDurationOk

`func (o *BackupAPIBackupStatus) GetDurationOk() (*int64, bool)`

GetDurationOk returns a tuple with the Duration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDuration

`func (o *BackupAPIBackupStatus) SetDuration(v int64)`

SetDuration sets Duration field to given value.


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

`func (o *BackupAPIBackupStatus) GetSize() int64`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *BackupAPIBackupStatus) GetSizeOk() (*int64, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *BackupAPIBackupStatus) SetSize(v int64)`

SetSize sets Size field to given value.


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



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


