/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.4.9
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
)

// DhcpAPILeaseInfo struct for DhcpAPILeaseInfo
type DhcpAPILeaseInfo struct {
	Vendor *string `json:"vendor,omitempty"`
}

// NewDhcpAPILeaseInfo instantiates a new DhcpAPILeaseInfo object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDhcpAPILeaseInfo() *DhcpAPILeaseInfo {
	this := DhcpAPILeaseInfo{}
	return &this
}

// NewDhcpAPILeaseInfoWithDefaults instantiates a new DhcpAPILeaseInfo object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDhcpAPILeaseInfoWithDefaults() *DhcpAPILeaseInfo {
	this := DhcpAPILeaseInfo{}
	return &this
}

// GetVendor returns the Vendor field value if set, zero value otherwise.
func (o *DhcpAPILeaseInfo) GetVendor() string {
	if o == nil || o.Vendor == nil {
		var ret string
		return ret
	}
	return *o.Vendor
}

// GetVendorOk returns a tuple with the Vendor field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DhcpAPILeaseInfo) GetVendorOk() (*string, bool) {
	if o == nil || o.Vendor == nil {
		return nil, false
	}
	return o.Vendor, true
}

// HasVendor returns a boolean if a field has been set.
func (o *DhcpAPILeaseInfo) HasVendor() bool {
	if o != nil && o.Vendor != nil {
		return true
	}

	return false
}

// SetVendor gets a reference to the given string and assigns it to the Vendor field.
func (o *DhcpAPILeaseInfo) SetVendor(v string) {
	o.Vendor = &v
}

func (o DhcpAPILeaseInfo) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Vendor != nil {
		toSerialize["vendor"] = o.Vendor
	}
	return json.Marshal(toSerialize)
}

type NullableDhcpAPILeaseInfo struct {
	value *DhcpAPILeaseInfo
	isSet bool
}

func (v NullableDhcpAPILeaseInfo) Get() *DhcpAPILeaseInfo {
	return v.value
}

func (v *NullableDhcpAPILeaseInfo) Set(val *DhcpAPILeaseInfo) {
	v.value = val
	v.isSet = true
}

func (v NullableDhcpAPILeaseInfo) IsSet() bool {
	return v.isSet
}

func (v *NullableDhcpAPILeaseInfo) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDhcpAPILeaseInfo(val *DhcpAPILeaseInfo) *NullableDhcpAPILeaseInfo {
	return &NullableDhcpAPILeaseInfo{value: val, isSet: true}
}

func (v NullableDhcpAPILeaseInfo) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDhcpAPILeaseInfo) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
