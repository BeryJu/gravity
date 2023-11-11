/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.6.18
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
)

// checks if the ExtconfigExtConfigDirs type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ExtconfigExtConfigDirs{}

// ExtconfigExtConfigDirs struct for ExtconfigExtConfigDirs
type ExtconfigExtConfigDirs struct {
	BackupDir *string `json:"backupDir,omitempty"`
	CertDir   *string `json:"certDir,omitempty"`
	EtcdDir   *string `json:"etcdDir,omitempty"`
}

// NewExtconfigExtConfigDirs instantiates a new ExtconfigExtConfigDirs object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewExtconfigExtConfigDirs() *ExtconfigExtConfigDirs {
	this := ExtconfigExtConfigDirs{}
	return &this
}

// NewExtconfigExtConfigDirsWithDefaults instantiates a new ExtconfigExtConfigDirs object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewExtconfigExtConfigDirsWithDefaults() *ExtconfigExtConfigDirs {
	this := ExtconfigExtConfigDirs{}
	return &this
}

// GetBackupDir returns the BackupDir field value if set, zero value otherwise.
func (o *ExtconfigExtConfigDirs) GetBackupDir() string {
	if o == nil || IsNil(o.BackupDir) {
		var ret string
		return ret
	}
	return *o.BackupDir
}

// GetBackupDirOk returns a tuple with the BackupDir field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ExtconfigExtConfigDirs) GetBackupDirOk() (*string, bool) {
	if o == nil || IsNil(o.BackupDir) {
		return nil, false
	}
	return o.BackupDir, true
}

// HasBackupDir returns a boolean if a field has been set.
func (o *ExtconfigExtConfigDirs) HasBackupDir() bool {
	if o != nil && !IsNil(o.BackupDir) {
		return true
	}

	return false
}

// SetBackupDir gets a reference to the given string and assigns it to the BackupDir field.
func (o *ExtconfigExtConfigDirs) SetBackupDir(v string) {
	o.BackupDir = &v
}

// GetCertDir returns the CertDir field value if set, zero value otherwise.
func (o *ExtconfigExtConfigDirs) GetCertDir() string {
	if o == nil || IsNil(o.CertDir) {
		var ret string
		return ret
	}
	return *o.CertDir
}

// GetCertDirOk returns a tuple with the CertDir field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ExtconfigExtConfigDirs) GetCertDirOk() (*string, bool) {
	if o == nil || IsNil(o.CertDir) {
		return nil, false
	}
	return o.CertDir, true
}

// HasCertDir returns a boolean if a field has been set.
func (o *ExtconfigExtConfigDirs) HasCertDir() bool {
	if o != nil && !IsNil(o.CertDir) {
		return true
	}

	return false
}

// SetCertDir gets a reference to the given string and assigns it to the CertDir field.
func (o *ExtconfigExtConfigDirs) SetCertDir(v string) {
	o.CertDir = &v
}

// GetEtcdDir returns the EtcdDir field value if set, zero value otherwise.
func (o *ExtconfigExtConfigDirs) GetEtcdDir() string {
	if o == nil || IsNil(o.EtcdDir) {
		var ret string
		return ret
	}
	return *o.EtcdDir
}

// GetEtcdDirOk returns a tuple with the EtcdDir field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ExtconfigExtConfigDirs) GetEtcdDirOk() (*string, bool) {
	if o == nil || IsNil(o.EtcdDir) {
		return nil, false
	}
	return o.EtcdDir, true
}

// HasEtcdDir returns a boolean if a field has been set.
func (o *ExtconfigExtConfigDirs) HasEtcdDir() bool {
	if o != nil && !IsNil(o.EtcdDir) {
		return true
	}

	return false
}

// SetEtcdDir gets a reference to the given string and assigns it to the EtcdDir field.
func (o *ExtconfigExtConfigDirs) SetEtcdDir(v string) {
	o.EtcdDir = &v
}

func (o ExtconfigExtConfigDirs) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ExtconfigExtConfigDirs) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.BackupDir) {
		toSerialize["backupDir"] = o.BackupDir
	}
	if !IsNil(o.CertDir) {
		toSerialize["certDir"] = o.CertDir
	}
	if !IsNil(o.EtcdDir) {
		toSerialize["etcdDir"] = o.EtcdDir
	}
	return toSerialize, nil
}

type NullableExtconfigExtConfigDirs struct {
	value *ExtconfigExtConfigDirs
	isSet bool
}

func (v NullableExtconfigExtConfigDirs) Get() *ExtconfigExtConfigDirs {
	return v.value
}

func (v *NullableExtconfigExtConfigDirs) Set(val *ExtconfigExtConfigDirs) {
	v.value = val
	v.isSet = true
}

func (v NullableExtconfigExtConfigDirs) IsSet() bool {
	return v.isSet
}

func (v *NullableExtconfigExtConfigDirs) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableExtconfigExtConfigDirs(val *ExtconfigExtConfigDirs) *NullableExtconfigExtConfigDirs {
	return &NullableExtconfigExtConfigDirs{value: val, isSet: true}
}

func (v NullableExtconfigExtConfigDirs) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableExtconfigExtConfigDirs) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
