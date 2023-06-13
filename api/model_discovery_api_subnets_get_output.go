/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.6.4
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
)

// DiscoveryAPISubnetsGetOutput struct for DiscoveryAPISubnetsGetOutput
type DiscoveryAPISubnetsGetOutput struct {
	Subnets []DiscoveryAPISubnet `json:"subnets,omitempty"`
}

// NewDiscoveryAPISubnetsGetOutput instantiates a new DiscoveryAPISubnetsGetOutput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDiscoveryAPISubnetsGetOutput() *DiscoveryAPISubnetsGetOutput {
	this := DiscoveryAPISubnetsGetOutput{}
	return &this
}

// NewDiscoveryAPISubnetsGetOutputWithDefaults instantiates a new DiscoveryAPISubnetsGetOutput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDiscoveryAPISubnetsGetOutputWithDefaults() *DiscoveryAPISubnetsGetOutput {
	this := DiscoveryAPISubnetsGetOutput{}
	return &this
}

// GetSubnets returns the Subnets field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *DiscoveryAPISubnetsGetOutput) GetSubnets() []DiscoveryAPISubnet {
	if o == nil {
		var ret []DiscoveryAPISubnet
		return ret
	}
	return o.Subnets
}

// GetSubnetsOk returns a tuple with the Subnets field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *DiscoveryAPISubnetsGetOutput) GetSubnetsOk() ([]DiscoveryAPISubnet, bool) {
	if o == nil || o.Subnets == nil {
		return nil, false
	}
	return o.Subnets, true
}

// HasSubnets returns a boolean if a field has been set.
func (o *DiscoveryAPISubnetsGetOutput) HasSubnets() bool {
	if o != nil && o.Subnets != nil {
		return true
	}

	return false
}

// SetSubnets gets a reference to the given []DiscoveryAPISubnet and assigns it to the Subnets field.
func (o *DiscoveryAPISubnetsGetOutput) SetSubnets(v []DiscoveryAPISubnet) {
	o.Subnets = v
}

func (o DiscoveryAPISubnetsGetOutput) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Subnets != nil {
		toSerialize["subnets"] = o.Subnets
	}
	return json.Marshal(toSerialize)
}

type NullableDiscoveryAPISubnetsGetOutput struct {
	value *DiscoveryAPISubnetsGetOutput
	isSet bool
}

func (v NullableDiscoveryAPISubnetsGetOutput) Get() *DiscoveryAPISubnetsGetOutput {
	return v.value
}

func (v *NullableDiscoveryAPISubnetsGetOutput) Set(val *DiscoveryAPISubnetsGetOutput) {
	v.value = val
	v.isSet = true
}

func (v NullableDiscoveryAPISubnetsGetOutput) IsSet() bool {
	return v.isSet
}

func (v *NullableDiscoveryAPISubnetsGetOutput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDiscoveryAPISubnetsGetOutput(val *DiscoveryAPISubnetsGetOutput) *NullableDiscoveryAPISubnetsGetOutput {
	return &NullableDiscoveryAPISubnetsGetOutput{value: val, isSet: true}
}

func (v NullableDiscoveryAPISubnetsGetOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDiscoveryAPISubnetsGetOutput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
