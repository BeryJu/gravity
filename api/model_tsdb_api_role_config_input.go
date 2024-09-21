/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.9.1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
)

// checks if the TsdbAPIRoleConfigInput type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TsdbAPIRoleConfigInput{}

// TsdbAPIRoleConfigInput struct for TsdbAPIRoleConfigInput
type TsdbAPIRoleConfigInput struct {
	Config TsdbRoleConfig `json:"config"`
}

// NewTsdbAPIRoleConfigInput instantiates a new TsdbAPIRoleConfigInput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTsdbAPIRoleConfigInput(config TsdbRoleConfig) *TsdbAPIRoleConfigInput {
	this := TsdbAPIRoleConfigInput{}
	this.Config = config
	return &this
}

// NewTsdbAPIRoleConfigInputWithDefaults instantiates a new TsdbAPIRoleConfigInput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTsdbAPIRoleConfigInputWithDefaults() *TsdbAPIRoleConfigInput {
	this := TsdbAPIRoleConfigInput{}
	return &this
}

// GetConfig returns the Config field value
func (o *TsdbAPIRoleConfigInput) GetConfig() TsdbRoleConfig {
	if o == nil {
		var ret TsdbRoleConfig
		return ret
	}

	return o.Config
}

// GetConfigOk returns a tuple with the Config field value
// and a boolean to check if the value has been set.
func (o *TsdbAPIRoleConfigInput) GetConfigOk() (*TsdbRoleConfig, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Config, true
}

// SetConfig sets field value
func (o *TsdbAPIRoleConfigInput) SetConfig(v TsdbRoleConfig) {
	o.Config = v
}

func (o TsdbAPIRoleConfigInput) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TsdbAPIRoleConfigInput) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["config"] = o.Config
	return toSerialize, nil
}

type NullableTsdbAPIRoleConfigInput struct {
	value *TsdbAPIRoleConfigInput
	isSet bool
}

func (v NullableTsdbAPIRoleConfigInput) Get() *TsdbAPIRoleConfigInput {
	return v.value
}

func (v *NullableTsdbAPIRoleConfigInput) Set(val *TsdbAPIRoleConfigInput) {
	v.value = val
	v.isSet = true
}

func (v NullableTsdbAPIRoleConfigInput) IsSet() bool {
	return v.isSet
}

func (v *NullableTsdbAPIRoleConfigInput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTsdbAPIRoleConfigInput(val *TsdbAPIRoleConfigInput) *NullableTsdbAPIRoleConfigInput {
	return &NullableTsdbAPIRoleConfigInput{value: val, isSet: true}
}

func (v NullableTsdbAPIRoleConfigInput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTsdbAPIRoleConfigInput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
