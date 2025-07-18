/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.27.2
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
)

// checks if the DnsAPIZonesImportInput type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &DnsAPIZonesImportInput{}

// DnsAPIZonesImportInput struct for DnsAPIZonesImportInput
type DnsAPIZonesImportInput struct {
	Payload *string                  `json:"payload,omitempty"`
	Type    *DnsAPIZonesImporterType `json:"type,omitempty"`
}

// NewDnsAPIZonesImportInput instantiates a new DnsAPIZonesImportInput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDnsAPIZonesImportInput() *DnsAPIZonesImportInput {
	this := DnsAPIZonesImportInput{}
	return &this
}

// NewDnsAPIZonesImportInputWithDefaults instantiates a new DnsAPIZonesImportInput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDnsAPIZonesImportInputWithDefaults() *DnsAPIZonesImportInput {
	this := DnsAPIZonesImportInput{}
	return &this
}

// GetPayload returns the Payload field value if set, zero value otherwise.
func (o *DnsAPIZonesImportInput) GetPayload() string {
	if o == nil || IsNil(o.Payload) {
		var ret string
		return ret
	}
	return *o.Payload
}

// GetPayloadOk returns a tuple with the Payload field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DnsAPIZonesImportInput) GetPayloadOk() (*string, bool) {
	if o == nil || IsNil(o.Payload) {
		return nil, false
	}
	return o.Payload, true
}

// HasPayload returns a boolean if a field has been set.
func (o *DnsAPIZonesImportInput) HasPayload() bool {
	if o != nil && !IsNil(o.Payload) {
		return true
	}

	return false
}

// SetPayload gets a reference to the given string and assigns it to the Payload field.
func (o *DnsAPIZonesImportInput) SetPayload(v string) {
	o.Payload = &v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *DnsAPIZonesImportInput) GetType() DnsAPIZonesImporterType {
	if o == nil || IsNil(o.Type) {
		var ret DnsAPIZonesImporterType
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DnsAPIZonesImportInput) GetTypeOk() (*DnsAPIZonesImporterType, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *DnsAPIZonesImportInput) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given DnsAPIZonesImporterType and assigns it to the Type field.
func (o *DnsAPIZonesImportInput) SetType(v DnsAPIZonesImporterType) {
	o.Type = &v
}

func (o DnsAPIZonesImportInput) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o DnsAPIZonesImportInput) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Payload) {
		toSerialize["payload"] = o.Payload
	}
	if !IsNil(o.Type) {
		toSerialize["type"] = o.Type
	}
	return toSerialize, nil
}

type NullableDnsAPIZonesImportInput struct {
	value *DnsAPIZonesImportInput
	isSet bool
}

func (v NullableDnsAPIZonesImportInput) Get() *DnsAPIZonesImportInput {
	return v.value
}

func (v *NullableDnsAPIZonesImportInput) Set(val *DnsAPIZonesImportInput) {
	v.value = val
	v.isSet = true
}

func (v NullableDnsAPIZonesImportInput) IsSet() bool {
	return v.isSet
}

func (v *NullableDnsAPIZonesImportInput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDnsAPIZonesImportInput(val *DnsAPIZonesImportInput) *NullableDnsAPIZonesImportInput {
	return &NullableDnsAPIZonesImportInput{value: val, isSet: true}
}

func (v NullableDnsAPIZonesImportInput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDnsAPIZonesImportInput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
