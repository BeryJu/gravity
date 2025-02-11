/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.26.2
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
	"time"
)

// checks if the BackupAPIBackupStatus type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &BackupAPIBackupStatus{}

// BackupAPIBackupStatus struct for BackupAPIBackupStatus
type BackupAPIBackupStatus struct {
	Duration int32     `json:"duration"`
	Error    string    `json:"error"`
	Filename string    `json:"filename"`
	Node     *string   `json:"node,omitempty"`
	Size     int32     `json:"size"`
	Status   string    `json:"status"`
	Time     time.Time `json:"time"`
}

// NewBackupAPIBackupStatus instantiates a new BackupAPIBackupStatus object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBackupAPIBackupStatus(duration int32, error_ string, filename string, size int32, status string, time time.Time) *BackupAPIBackupStatus {
	this := BackupAPIBackupStatus{}
	this.Duration = duration
	this.Error = error_
	this.Filename = filename
	this.Size = size
	this.Status = status
	this.Time = time
	return &this
}

// NewBackupAPIBackupStatusWithDefaults instantiates a new BackupAPIBackupStatus object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBackupAPIBackupStatusWithDefaults() *BackupAPIBackupStatus {
	this := BackupAPIBackupStatus{}
	return &this
}

// GetDuration returns the Duration field value
func (o *BackupAPIBackupStatus) GetDuration() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Duration
}

// GetDurationOk returns a tuple with the Duration field value
// and a boolean to check if the value has been set.
func (o *BackupAPIBackupStatus) GetDurationOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Duration, true
}

// SetDuration sets field value
func (o *BackupAPIBackupStatus) SetDuration(v int32) {
	o.Duration = v
}

// GetError returns the Error field value
func (o *BackupAPIBackupStatus) GetError() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Error
}

// GetErrorOk returns a tuple with the Error field value
// and a boolean to check if the value has been set.
func (o *BackupAPIBackupStatus) GetErrorOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Error, true
}

// SetError sets field value
func (o *BackupAPIBackupStatus) SetError(v string) {
	o.Error = v
}

// GetFilename returns the Filename field value
func (o *BackupAPIBackupStatus) GetFilename() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Filename
}

// GetFilenameOk returns a tuple with the Filename field value
// and a boolean to check if the value has been set.
func (o *BackupAPIBackupStatus) GetFilenameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Filename, true
}

// SetFilename sets field value
func (o *BackupAPIBackupStatus) SetFilename(v string) {
	o.Filename = v
}

// GetNode returns the Node field value if set, zero value otherwise.
func (o *BackupAPIBackupStatus) GetNode() string {
	if o == nil || IsNil(o.Node) {
		var ret string
		return ret
	}
	return *o.Node
}

// GetNodeOk returns a tuple with the Node field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BackupAPIBackupStatus) GetNodeOk() (*string, bool) {
	if o == nil || IsNil(o.Node) {
		return nil, false
	}
	return o.Node, true
}

// HasNode returns a boolean if a field has been set.
func (o *BackupAPIBackupStatus) HasNode() bool {
	if o != nil && !IsNil(o.Node) {
		return true
	}

	return false
}

// SetNode gets a reference to the given string and assigns it to the Node field.
func (o *BackupAPIBackupStatus) SetNode(v string) {
	o.Node = &v
}

// GetSize returns the Size field value
func (o *BackupAPIBackupStatus) GetSize() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Size
}

// GetSizeOk returns a tuple with the Size field value
// and a boolean to check if the value has been set.
func (o *BackupAPIBackupStatus) GetSizeOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Size, true
}

// SetSize sets field value
func (o *BackupAPIBackupStatus) SetSize(v int32) {
	o.Size = v
}

// GetStatus returns the Status field value
func (o *BackupAPIBackupStatus) GetStatus() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Status
}

// GetStatusOk returns a tuple with the Status field value
// and a boolean to check if the value has been set.
func (o *BackupAPIBackupStatus) GetStatusOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Status, true
}

// SetStatus sets field value
func (o *BackupAPIBackupStatus) SetStatus(v string) {
	o.Status = v
}

// GetTime returns the Time field value
func (o *BackupAPIBackupStatus) GetTime() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.Time
}

// GetTimeOk returns a tuple with the Time field value
// and a boolean to check if the value has been set.
func (o *BackupAPIBackupStatus) GetTimeOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Time, true
}

// SetTime sets field value
func (o *BackupAPIBackupStatus) SetTime(v time.Time) {
	o.Time = v
}

func (o BackupAPIBackupStatus) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o BackupAPIBackupStatus) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["duration"] = o.Duration
	toSerialize["error"] = o.Error
	toSerialize["filename"] = o.Filename
	if !IsNil(o.Node) {
		toSerialize["node"] = o.Node
	}
	toSerialize["size"] = o.Size
	toSerialize["status"] = o.Status
	toSerialize["time"] = o.Time
	return toSerialize, nil
}

type NullableBackupAPIBackupStatus struct {
	value *BackupAPIBackupStatus
	isSet bool
}

func (v NullableBackupAPIBackupStatus) Get() *BackupAPIBackupStatus {
	return v.value
}

func (v *NullableBackupAPIBackupStatus) Set(val *BackupAPIBackupStatus) {
	v.value = val
	v.isSet = true
}

func (v NullableBackupAPIBackupStatus) IsSet() bool {
	return v.isSet
}

func (v *NullableBackupAPIBackupStatus) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBackupAPIBackupStatus(val *BackupAPIBackupStatus) *NullableBackupAPIBackupStatus {
	return &NullableBackupAPIBackupStatus{value: val, isSet: true}
}

func (v NullableBackupAPIBackupStatus) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBackupAPIBackupStatus) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
