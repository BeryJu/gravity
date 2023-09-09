/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.6.12
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
)

// checks if the DhcpAPIScopesGetOutput type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &DhcpAPIScopesGetOutput{}

// DhcpAPIScopesGetOutput struct for DhcpAPIScopesGetOutput
type DhcpAPIScopesGetOutput struct {
	Scopes []DhcpAPIScope `json:"scopes"`
}

// NewDhcpAPIScopesGetOutput instantiates a new DhcpAPIScopesGetOutput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDhcpAPIScopesGetOutput(scopes []DhcpAPIScope) *DhcpAPIScopesGetOutput {
	this := DhcpAPIScopesGetOutput{}
	this.Scopes = scopes
	return &this
}

// NewDhcpAPIScopesGetOutputWithDefaults instantiates a new DhcpAPIScopesGetOutput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDhcpAPIScopesGetOutputWithDefaults() *DhcpAPIScopesGetOutput {
	this := DhcpAPIScopesGetOutput{}
	return &this
}

// GetScopes returns the Scopes field value
// If the value is explicit nil, the zero value for []DhcpAPIScope will be returned
func (o *DhcpAPIScopesGetOutput) GetScopes() []DhcpAPIScope {
	if o == nil {
		var ret []DhcpAPIScope
		return ret
	}

	return o.Scopes
}

// GetScopesOk returns a tuple with the Scopes field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *DhcpAPIScopesGetOutput) GetScopesOk() ([]DhcpAPIScope, bool) {
	if o == nil || IsNil(o.Scopes) {
		return nil, false
	}
	return o.Scopes, true
}

// SetScopes sets field value
func (o *DhcpAPIScopesGetOutput) SetScopes(v []DhcpAPIScope) {
	o.Scopes = v
}

func (o DhcpAPIScopesGetOutput) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o DhcpAPIScopesGetOutput) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if o.Scopes != nil {
		toSerialize["scopes"] = o.Scopes
	}
	return toSerialize, nil
}

type NullableDhcpAPIScopesGetOutput struct {
	value *DhcpAPIScopesGetOutput
	isSet bool
}

func (v NullableDhcpAPIScopesGetOutput) Get() *DhcpAPIScopesGetOutput {
	return v.value
}

func (v *NullableDhcpAPIScopesGetOutput) Set(val *DhcpAPIScopesGetOutput) {
	v.value = val
	v.isSet = true
}

func (v NullableDhcpAPIScopesGetOutput) IsSet() bool {
	return v.isSet
}

func (v *NullableDhcpAPIScopesGetOutput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDhcpAPIScopesGetOutput(val *DhcpAPIScopesGetOutput) *NullableDhcpAPIScopesGetOutput {
	return &NullableDhcpAPIScopesGetOutput{value: val, isSet: true}
}

func (v NullableDhcpAPIScopesGetOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDhcpAPIScopesGetOutput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
