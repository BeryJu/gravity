/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.8.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
)

// checks if the AuthAPILoginOutput type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AuthAPILoginOutput{}

// AuthAPILoginOutput struct for AuthAPILoginOutput
type AuthAPILoginOutput struct {
	Successful *bool `json:"successful,omitempty"`
}

// NewAuthAPILoginOutput instantiates a new AuthAPILoginOutput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAuthAPILoginOutput() *AuthAPILoginOutput {
	this := AuthAPILoginOutput{}
	return &this
}

// NewAuthAPILoginOutputWithDefaults instantiates a new AuthAPILoginOutput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAuthAPILoginOutputWithDefaults() *AuthAPILoginOutput {
	this := AuthAPILoginOutput{}
	return &this
}

// GetSuccessful returns the Successful field value if set, zero value otherwise.
func (o *AuthAPILoginOutput) GetSuccessful() bool {
	if o == nil || IsNil(o.Successful) {
		var ret bool
		return ret
	}
	return *o.Successful
}

// GetSuccessfulOk returns a tuple with the Successful field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AuthAPILoginOutput) GetSuccessfulOk() (*bool, bool) {
	if o == nil || IsNil(o.Successful) {
		return nil, false
	}
	return o.Successful, true
}

// HasSuccessful returns a boolean if a field has been set.
func (o *AuthAPILoginOutput) HasSuccessful() bool {
	if o != nil && !IsNil(o.Successful) {
		return true
	}

	return false
}

// SetSuccessful gets a reference to the given bool and assigns it to the Successful field.
func (o *AuthAPILoginOutput) SetSuccessful(v bool) {
	o.Successful = &v
}

func (o AuthAPILoginOutput) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o AuthAPILoginOutput) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Successful) {
		toSerialize["successful"] = o.Successful
	}
	return toSerialize, nil
}

type NullableAuthAPILoginOutput struct {
	value *AuthAPILoginOutput
	isSet bool
}

func (v NullableAuthAPILoginOutput) Get() *AuthAPILoginOutput {
	return v.value
}

func (v *NullableAuthAPILoginOutput) Set(val *AuthAPILoginOutput) {
	v.value = val
	v.isSet = true
}

func (v NullableAuthAPILoginOutput) IsSet() bool {
	return v.isSet
}

func (v *NullableAuthAPILoginOutput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAuthAPILoginOutput(val *AuthAPILoginOutput) *NullableAuthAPILoginOutput {
	return &NullableAuthAPILoginOutput{value: val, isSet: true}
}

func (v NullableAuthAPILoginOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAuthAPILoginOutput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
