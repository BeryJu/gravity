/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.18.2
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
)

// checks if the DiscoveryRoleConfig type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &DiscoveryRoleConfig{}

// DiscoveryRoleConfig struct for DiscoveryRoleConfig
type DiscoveryRoleConfig struct {
	Enabled *bool `json:"enabled,omitempty"`
}

// NewDiscoveryRoleConfig instantiates a new DiscoveryRoleConfig object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDiscoveryRoleConfig() *DiscoveryRoleConfig {
	this := DiscoveryRoleConfig{}
	return &this
}

// NewDiscoveryRoleConfigWithDefaults instantiates a new DiscoveryRoleConfig object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDiscoveryRoleConfigWithDefaults() *DiscoveryRoleConfig {
	this := DiscoveryRoleConfig{}
	return &this
}

// GetEnabled returns the Enabled field value if set, zero value otherwise.
func (o *DiscoveryRoleConfig) GetEnabled() bool {
	if o == nil || IsNil(o.Enabled) {
		var ret bool
		return ret
	}
	return *o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DiscoveryRoleConfig) GetEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.Enabled) {
		return nil, false
	}
	return o.Enabled, true
}

// HasEnabled returns a boolean if a field has been set.
func (o *DiscoveryRoleConfig) HasEnabled() bool {
	if o != nil && !IsNil(o.Enabled) {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given bool and assigns it to the Enabled field.
func (o *DiscoveryRoleConfig) SetEnabled(v bool) {
	o.Enabled = &v
}

func (o DiscoveryRoleConfig) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o DiscoveryRoleConfig) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Enabled) {
		toSerialize["enabled"] = o.Enabled
	}
	return toSerialize, nil
}

type NullableDiscoveryRoleConfig struct {
	value *DiscoveryRoleConfig
	isSet bool
}

func (v NullableDiscoveryRoleConfig) Get() *DiscoveryRoleConfig {
	return v.value
}

func (v *NullableDiscoveryRoleConfig) Set(val *DiscoveryRoleConfig) {
	v.value = val
	v.isSet = true
}

func (v NullableDiscoveryRoleConfig) IsSet() bool {
	return v.isSet
}

func (v *NullableDiscoveryRoleConfig) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDiscoveryRoleConfig(val *DiscoveryRoleConfig) *NullableDiscoveryRoleConfig {
	return &NullableDiscoveryRoleConfig{value: val, isSet: true}
}

func (v NullableDiscoveryRoleConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDiscoveryRoleConfig) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
