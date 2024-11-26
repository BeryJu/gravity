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

// checks if the MonitoringAPIRoleConfigInput type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &MonitoringAPIRoleConfigInput{}

// MonitoringAPIRoleConfigInput struct for MonitoringAPIRoleConfigInput
type MonitoringAPIRoleConfigInput struct {
	Config MonitoringRoleConfig `json:"config"`
}

// NewMonitoringAPIRoleConfigInput instantiates a new MonitoringAPIRoleConfigInput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewMonitoringAPIRoleConfigInput(config MonitoringRoleConfig) *MonitoringAPIRoleConfigInput {
	this := MonitoringAPIRoleConfigInput{}
	this.Config = config
	return &this
}

// NewMonitoringAPIRoleConfigInputWithDefaults instantiates a new MonitoringAPIRoleConfigInput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewMonitoringAPIRoleConfigInputWithDefaults() *MonitoringAPIRoleConfigInput {
	this := MonitoringAPIRoleConfigInput{}
	return &this
}

// GetConfig returns the Config field value
func (o *MonitoringAPIRoleConfigInput) GetConfig() MonitoringRoleConfig {
	if o == nil {
		var ret MonitoringRoleConfig
		return ret
	}

	return o.Config
}

// GetConfigOk returns a tuple with the Config field value
// and a boolean to check if the value has been set.
func (o *MonitoringAPIRoleConfigInput) GetConfigOk() (*MonitoringRoleConfig, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Config, true
}

// SetConfig sets field value
func (o *MonitoringAPIRoleConfigInput) SetConfig(v MonitoringRoleConfig) {
	o.Config = v
}

func (o MonitoringAPIRoleConfigInput) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o MonitoringAPIRoleConfigInput) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["config"] = o.Config
	return toSerialize, nil
}

type NullableMonitoringAPIRoleConfigInput struct {
	value *MonitoringAPIRoleConfigInput
	isSet bool
}

func (v NullableMonitoringAPIRoleConfigInput) Get() *MonitoringAPIRoleConfigInput {
	return v.value
}

func (v *NullableMonitoringAPIRoleConfigInput) Set(val *MonitoringAPIRoleConfigInput) {
	v.value = val
	v.isSet = true
}

func (v NullableMonitoringAPIRoleConfigInput) IsSet() bool {
	return v.isSet
}

func (v *NullableMonitoringAPIRoleConfigInput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableMonitoringAPIRoleConfigInput(val *MonitoringAPIRoleConfigInput) *NullableMonitoringAPIRoleConfigInput {
	return &NullableMonitoringAPIRoleConfigInput{value: val, isSet: true}
}

func (v NullableMonitoringAPIRoleConfigInput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableMonitoringAPIRoleConfigInput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
