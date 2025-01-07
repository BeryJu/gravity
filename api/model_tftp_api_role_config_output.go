/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.21.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
)

// checks if the TftpAPIRoleConfigOutput type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TftpAPIRoleConfigOutput{}

// TftpAPIRoleConfigOutput struct for TftpAPIRoleConfigOutput
type TftpAPIRoleConfigOutput struct {
	Config TftpRoleConfig `json:"config"`
}

// NewTftpAPIRoleConfigOutput instantiates a new TftpAPIRoleConfigOutput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTftpAPIRoleConfigOutput(config TftpRoleConfig) *TftpAPIRoleConfigOutput {
	this := TftpAPIRoleConfigOutput{}
	this.Config = config
	return &this
}

// NewTftpAPIRoleConfigOutputWithDefaults instantiates a new TftpAPIRoleConfigOutput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTftpAPIRoleConfigOutputWithDefaults() *TftpAPIRoleConfigOutput {
	this := TftpAPIRoleConfigOutput{}
	return &this
}

// GetConfig returns the Config field value
func (o *TftpAPIRoleConfigOutput) GetConfig() TftpRoleConfig {
	if o == nil {
		var ret TftpRoleConfig
		return ret
	}

	return o.Config
}

// GetConfigOk returns a tuple with the Config field value
// and a boolean to check if the value has been set.
func (o *TftpAPIRoleConfigOutput) GetConfigOk() (*TftpRoleConfig, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Config, true
}

// SetConfig sets field value
func (o *TftpAPIRoleConfigOutput) SetConfig(v TftpRoleConfig) {
	o.Config = v
}

func (o TftpAPIRoleConfigOutput) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TftpAPIRoleConfigOutput) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["config"] = o.Config
	return toSerialize, nil
}

type NullableTftpAPIRoleConfigOutput struct {
	value *TftpAPIRoleConfigOutput
	isSet bool
}

func (v NullableTftpAPIRoleConfigOutput) Get() *TftpAPIRoleConfigOutput {
	return v.value
}

func (v *NullableTftpAPIRoleConfigOutput) Set(val *TftpAPIRoleConfigOutput) {
	v.value = val
	v.isSet = true
}

func (v NullableTftpAPIRoleConfigOutput) IsSet() bool {
	return v.isSet
}

func (v *NullableTftpAPIRoleConfigOutput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTftpAPIRoleConfigOutput(val *TftpAPIRoleConfigOutput) *NullableTftpAPIRoleConfigOutput {
	return &NullableTftpAPIRoleConfigOutput{value: val, isSet: true}
}

func (v NullableTftpAPIRoleConfigOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTftpAPIRoleConfigOutput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
