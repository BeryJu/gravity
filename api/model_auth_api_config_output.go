/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.3.17
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
)

// AuthAPIConfigOutput struct for AuthAPIConfigOutput
type AuthAPIConfigOutput struct {
	Bool *bool `json:"bool,omitempty"`
	Oidc *bool `json:"oidc,omitempty"`
}

// NewAuthAPIConfigOutput instantiates a new AuthAPIConfigOutput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAuthAPIConfigOutput() *AuthAPIConfigOutput {
	this := AuthAPIConfigOutput{}
	return &this
}

// NewAuthAPIConfigOutputWithDefaults instantiates a new AuthAPIConfigOutput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAuthAPIConfigOutputWithDefaults() *AuthAPIConfigOutput {
	this := AuthAPIConfigOutput{}
	return &this
}

// GetBool returns the Bool field value if set, zero value otherwise.
func (o *AuthAPIConfigOutput) GetBool() bool {
	if o == nil || o.Bool == nil {
		var ret bool
		return ret
	}
	return *o.Bool
}

// GetBoolOk returns a tuple with the Bool field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AuthAPIConfigOutput) GetBoolOk() (*bool, bool) {
	if o == nil || o.Bool == nil {
		return nil, false
	}
	return o.Bool, true
}

// HasBool returns a boolean if a field has been set.
func (o *AuthAPIConfigOutput) HasBool() bool {
	if o != nil && o.Bool != nil {
		return true
	}

	return false
}

// SetBool gets a reference to the given bool and assigns it to the Bool field.
func (o *AuthAPIConfigOutput) SetBool(v bool) {
	o.Bool = &v
}

// GetOidc returns the Oidc field value if set, zero value otherwise.
func (o *AuthAPIConfigOutput) GetOidc() bool {
	if o == nil || o.Oidc == nil {
		var ret bool
		return ret
	}
	return *o.Oidc
}

// GetOidcOk returns a tuple with the Oidc field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *AuthAPIConfigOutput) GetOidcOk() (*bool, bool) {
	if o == nil || o.Oidc == nil {
		return nil, false
	}
	return o.Oidc, true
}

// HasOidc returns a boolean if a field has been set.
func (o *AuthAPIConfigOutput) HasOidc() bool {
	if o != nil && o.Oidc != nil {
		return true
	}

	return false
}

// SetOidc gets a reference to the given bool and assigns it to the Oidc field.
func (o *AuthAPIConfigOutput) SetOidc(v bool) {
	o.Oidc = &v
}

func (o AuthAPIConfigOutput) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Bool != nil {
		toSerialize["bool"] = o.Bool
	}
	if o.Oidc != nil {
		toSerialize["oidc"] = o.Oidc
	}
	return json.Marshal(toSerialize)
}

type NullableAuthAPIConfigOutput struct {
	value *AuthAPIConfigOutput
	isSet bool
}

func (v NullableAuthAPIConfigOutput) Get() *AuthAPIConfigOutput {
	return v.value
}

func (v *NullableAuthAPIConfigOutput) Set(val *AuthAPIConfigOutput) {
	v.value = val
	v.isSet = true
}

func (v NullableAuthAPIConfigOutput) IsSet() bool {
	return v.isSet
}

func (v *NullableAuthAPIConfigOutput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAuthAPIConfigOutput(val *AuthAPIConfigOutput) *NullableAuthAPIConfigOutput {
	return &NullableAuthAPIConfigOutput{value: val, isSet: true}
}

func (v NullableAuthAPIConfigOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAuthAPIConfigOutput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
