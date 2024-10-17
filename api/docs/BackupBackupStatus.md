# BackupBackupStatus

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Duration** | **int32** |  | 
**Error** | **string** |  | 
**Filename** | **string** |  | 
**Size** | **int32** |  | 
**Status** | **string** |  | 
**Time** | **time.Time** |  | 

## Methods

### NewBackupBackupStatus

`func NewBackupBackupStatus(duration int32, error_ string, filename string, size int32, status string, time time.Time, ) *BackupBackupStatus`

NewBackupBackupStatus instantiates a new BackupBackupStatus object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBackupBackupStatusWithDefaults

`func NewBackupBackupStatusWithDefaults() *BackupBackupStatus`

NewBackupBackupStatusWithDefaults instantiates a new BackupBackupStatus object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDuration

`func (o *BackupBackupStatus) GetDuration() int32`

GetDuration returns the Duration field if non-nil, zero value otherwise.

### GetDurationOk

`func (o *BackupBackupStatus) GetDurationOk() (*int32, bool)`

GetDurationOk returns a tuple with the Duration field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDuration

`func (o *BackupBackupStatus) SetDuration(v int32)`

SetDuration sets Duration field to given value.


### GetError

`func (o *BackupBackupStatus) GetError() string`

GetError returns the Error field if non-nil, zero value otherwise.

### GetErrorOk

`func (o *BackupBackupStatus) GetErrorOk() (*string, bool)`

GetErrorOk returns a tuple with the Error field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetError

`func (o *BackupBackupStatus) SetError(v string)`

SetError sets Error field to given value.


### GetFilename

`func (o *BackupBackupStatus) GetFilename() string`

GetFilename returns the Filename field if non-nil, zero value otherwise.

### GetFilenameOk

`func (o *BackupBackupStatus) GetFilenameOk() (*string, bool)`

GetFilenameOk returns a tuple with the Filename field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetFilename

`func (o *BackupBackupStatus) SetFilename(v string)`

SetFilename sets Filename field to given value.


### GetSize

`func (o *BackupBackupStatus) GetSize() int32`

GetSize returns the Size field if non-nil, zero value otherwise.

### GetSizeOk

`func (o *BackupBackupStatus) GetSizeOk() (*int32, bool)`

GetSizeOk returns a tuple with the Size field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSize

`func (o *BackupBackupStatus) SetSize(v int32)`

SetSize sets Size field to given value.


### GetStatus

`func (o *BackupBackupStatus) GetStatus() string`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *BackupBackupStatus) GetStatusOk() (*string, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *BackupBackupStatus) SetStatus(v string)`

SetStatus sets Status field to given value.


### GetTime

`func (o *BackupBackupStatus) GetTime() time.Time`

GetTime returns the Time field if non-nil, zero value otherwise.

### GetTimeOk

`func (o *BackupBackupStatus) GetTimeOk() (*time.Time, bool)`

GetTimeOk returns a tuple with the Time field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetTime

`func (o *BackupBackupStatus) SetTime(v time.Time)`

SetTime sets Time field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


