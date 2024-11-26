/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.17.2
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
)

// checks if the DiscoveryAPIRoleConfigInput type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &DiscoveryAPIRoleConfigInput{}

// DiscoveryAPIRoleConfigInput struct for DiscoveryAPIRoleConfigInput
type DiscoveryAPIRoleConfigInput struct {
	Config DiscoveryRoleConfig `json:"config"`
}

// NewDiscoveryAPIRoleConfigInput instantiates a new DiscoveryAPIRoleConfigInput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDiscoveryAPIRoleConfigInput(config DiscoveryRoleConfig) *DiscoveryAPIRoleConfigInput {
	this := DiscoveryAPIRoleConfigInput{}
	this.Config = config
	return &this
}

// NewDiscoveryAPIRoleConfigInputWithDefaults instantiates a new DiscoveryAPIRoleConfigInput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDiscoveryAPIRoleConfigInputWithDefaults() *DiscoveryAPIRoleConfigInput {
	this := DiscoveryAPIRoleConfigInput{}
	return &this
}

// GetConfig returns the Config field value
func (o *DiscoveryAPIRoleConfigInput) GetConfig() DiscoveryRoleConfig {
	if o == nil {
		var ret DiscoveryRoleConfig
		return ret
	}

	return o.Config
}

// GetConfigOk returns a tuple with the Config field value
// and a boolean to check if the value has been set.
func (o *DiscoveryAPIRoleConfigInput) GetConfigOk() (*DiscoveryRoleConfig, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Config, true
}

// SetConfig sets field value
func (o *DiscoveryAPIRoleConfigInput) SetConfig(v DiscoveryRoleConfig) {
	o.Config = v
}

func (o DiscoveryAPIRoleConfigInput) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o DiscoveryAPIRoleConfigInput) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["config"] = o.Config
	return toSerialize, nil
}

type NullableDiscoveryAPIRoleConfigInput struct {
	value *DiscoveryAPIRoleConfigInput
	isSet bool
}

func (v NullableDiscoveryAPIRoleConfigInput) Get() *DiscoveryAPIRoleConfigInput {
	return v.value
}

func (v *NullableDiscoveryAPIRoleConfigInput) Set(val *DiscoveryAPIRoleConfigInput) {
	v.value = val
	v.isSet = true
}

func (v NullableDiscoveryAPIRoleConfigInput) IsSet() bool {
	return v.isSet
}

func (v *NullableDiscoveryAPIRoleConfigInput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDiscoveryAPIRoleConfigInput(val *DiscoveryAPIRoleConfigInput) *NullableDiscoveryAPIRoleConfigInput {
	return &NullableDiscoveryAPIRoleConfigInput{value: val, isSet: true}
}

func (v NullableDiscoveryAPIRoleConfigInput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDiscoveryAPIRoleConfigInput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
