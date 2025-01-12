/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.23.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
)

// checks if the ApiAPIRoleConfigOutput type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ApiAPIRoleConfigOutput{}

// ApiAPIRoleConfigOutput struct for ApiAPIRoleConfigOutput
type ApiAPIRoleConfigOutput struct {
	Config ApiRoleConfig `json:"config"`
}

// NewApiAPIRoleConfigOutput instantiates a new ApiAPIRoleConfigOutput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiAPIRoleConfigOutput(config ApiRoleConfig) *ApiAPIRoleConfigOutput {
	this := ApiAPIRoleConfigOutput{}
	this.Config = config
	return &this
}

// NewApiAPIRoleConfigOutputWithDefaults instantiates a new ApiAPIRoleConfigOutput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiAPIRoleConfigOutputWithDefaults() *ApiAPIRoleConfigOutput {
	this := ApiAPIRoleConfigOutput{}
	return &this
}

// GetConfig returns the Config field value
func (o *ApiAPIRoleConfigOutput) GetConfig() ApiRoleConfig {
	if o == nil {
		var ret ApiRoleConfig
		return ret
	}

	return o.Config
}

// GetConfigOk returns a tuple with the Config field value
// and a boolean to check if the value has been set.
func (o *ApiAPIRoleConfigOutput) GetConfigOk() (*ApiRoleConfig, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Config, true
}

// SetConfig sets field value
func (o *ApiAPIRoleConfigOutput) SetConfig(v ApiRoleConfig) {
	o.Config = v
}

func (o ApiAPIRoleConfigOutput) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ApiAPIRoleConfigOutput) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["config"] = o.Config
	return toSerialize, nil
}

type NullableApiAPIRoleConfigOutput struct {
	value *ApiAPIRoleConfigOutput
	isSet bool
}

func (v NullableApiAPIRoleConfigOutput) Get() *ApiAPIRoleConfigOutput {
	return v.value
}

func (v *NullableApiAPIRoleConfigOutput) Set(val *ApiAPIRoleConfigOutput) {
	v.value = val
	v.isSet = true
}

func (v NullableApiAPIRoleConfigOutput) IsSet() bool {
	return v.isSet
}

func (v *NullableApiAPIRoleConfigOutput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableApiAPIRoleConfigOutput(val *ApiAPIRoleConfigOutput) *NullableApiAPIRoleConfigOutput {
	return &NullableApiAPIRoleConfigOutput{value: val, isSet: true}
}

func (v NullableApiAPIRoleConfigOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableApiAPIRoleConfigOutput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
