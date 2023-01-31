/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.4.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
)

// InstanceAPIRoleRestartInput struct for InstanceAPIRoleRestartInput
type InstanceAPIRoleRestartInput struct {
	RoleId *string `json:"roleId,omitempty"`
}

// NewInstanceAPIRoleRestartInput instantiates a new InstanceAPIRoleRestartInput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewInstanceAPIRoleRestartInput() *InstanceAPIRoleRestartInput {
	this := InstanceAPIRoleRestartInput{}
	return &this
}

// NewInstanceAPIRoleRestartInputWithDefaults instantiates a new InstanceAPIRoleRestartInput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewInstanceAPIRoleRestartInputWithDefaults() *InstanceAPIRoleRestartInput {
	this := InstanceAPIRoleRestartInput{}
	return &this
}

// GetRoleId returns the RoleId field value if set, zero value otherwise.
func (o *InstanceAPIRoleRestartInput) GetRoleId() string {
	if o == nil || o.RoleId == nil {
		var ret string
		return ret
	}
	return *o.RoleId
}

// GetRoleIdOk returns a tuple with the RoleId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *InstanceAPIRoleRestartInput) GetRoleIdOk() (*string, bool) {
	if o == nil || o.RoleId == nil {
		return nil, false
	}
	return o.RoleId, true
}

// HasRoleId returns a boolean if a field has been set.
func (o *InstanceAPIRoleRestartInput) HasRoleId() bool {
	if o != nil && o.RoleId != nil {
		return true
	}

	return false
}

// SetRoleId gets a reference to the given string and assigns it to the RoleId field.
func (o *InstanceAPIRoleRestartInput) SetRoleId(v string) {
	o.RoleId = &v
}

func (o InstanceAPIRoleRestartInput) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.RoleId != nil {
		toSerialize["roleId"] = o.RoleId
	}
	return json.Marshal(toSerialize)
}

type NullableInstanceAPIRoleRestartInput struct {
	value *InstanceAPIRoleRestartInput
	isSet bool
}

func (v NullableInstanceAPIRoleRestartInput) Get() *InstanceAPIRoleRestartInput {
	return v.value
}

func (v *NullableInstanceAPIRoleRestartInput) Set(val *InstanceAPIRoleRestartInput) {
	v.value = val
	v.isSet = true
}

func (v NullableInstanceAPIRoleRestartInput) IsSet() bool {
	return v.isSet
}

func (v *NullableInstanceAPIRoleRestartInput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableInstanceAPIRoleRestartInput(val *InstanceAPIRoleRestartInput) *NullableInstanceAPIRoleRestartInput {
	return &NullableInstanceAPIRoleRestartInput{value: val, isSet: true}
}

func (v NullableInstanceAPIRoleRestartInput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableInstanceAPIRoleRestartInput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
