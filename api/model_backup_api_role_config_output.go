/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.6.11
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
)

// checks if the BackupAPIRoleConfigOutput type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &BackupAPIRoleConfigOutput{}

// BackupAPIRoleConfigOutput struct for BackupAPIRoleConfigOutput
type BackupAPIRoleConfigOutput struct {
	Config BackupRoleConfig `json:"config"`
}

// NewBackupAPIRoleConfigOutput instantiates a new BackupAPIRoleConfigOutput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBackupAPIRoleConfigOutput(config BackupRoleConfig) *BackupAPIRoleConfigOutput {
	this := BackupAPIRoleConfigOutput{}
	this.Config = config
	return &this
}

// NewBackupAPIRoleConfigOutputWithDefaults instantiates a new BackupAPIRoleConfigOutput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBackupAPIRoleConfigOutputWithDefaults() *BackupAPIRoleConfigOutput {
	this := BackupAPIRoleConfigOutput{}
	return &this
}

// GetConfig returns the Config field value
func (o *BackupAPIRoleConfigOutput) GetConfig() BackupRoleConfig {
	if o == nil {
		var ret BackupRoleConfig
		return ret
	}

	return o.Config
}

// GetConfigOk returns a tuple with the Config field value
// and a boolean to check if the value has been set.
func (o *BackupAPIRoleConfigOutput) GetConfigOk() (*BackupRoleConfig, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Config, true
}

// SetConfig sets field value
func (o *BackupAPIRoleConfigOutput) SetConfig(v BackupRoleConfig) {
	o.Config = v
}

func (o BackupAPIRoleConfigOutput) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o BackupAPIRoleConfigOutput) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["config"] = o.Config
	return toSerialize, nil
}

type NullableBackupAPIRoleConfigOutput struct {
	value *BackupAPIRoleConfigOutput
	isSet bool
}

func (v NullableBackupAPIRoleConfigOutput) Get() *BackupAPIRoleConfigOutput {
	return v.value
}

func (v *NullableBackupAPIRoleConfigOutput) Set(val *BackupAPIRoleConfigOutput) {
	v.value = val
	v.isSet = true
}

func (v NullableBackupAPIRoleConfigOutput) IsSet() bool {
	return v.isSet
}

func (v *NullableBackupAPIRoleConfigOutput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBackupAPIRoleConfigOutput(val *BackupAPIRoleConfigOutput) *NullableBackupAPIRoleConfigOutput {
	return &NullableBackupAPIRoleConfigOutput{value: val, isSet: true}
}

func (v NullableBackupAPIRoleConfigOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBackupAPIRoleConfigOutput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
