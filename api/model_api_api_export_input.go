/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.17.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
)

// checks if the ApiAPIExportInput type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ApiAPIExportInput{}

// ApiAPIExportInput struct for ApiAPIExportInput
type ApiAPIExportInput struct {
	Safe *bool `json:"safe,omitempty"`
}

// NewApiAPIExportInput instantiates a new ApiAPIExportInput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiAPIExportInput() *ApiAPIExportInput {
	this := ApiAPIExportInput{}
	return &this
}

// NewApiAPIExportInputWithDefaults instantiates a new ApiAPIExportInput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiAPIExportInputWithDefaults() *ApiAPIExportInput {
	this := ApiAPIExportInput{}
	return &this
}

// GetSafe returns the Safe field value if set, zero value otherwise.
func (o *ApiAPIExportInput) GetSafe() bool {
	if o == nil || IsNil(o.Safe) {
		var ret bool
		return ret
	}
	return *o.Safe
}

// GetSafeOk returns a tuple with the Safe field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAPIExportInput) GetSafeOk() (*bool, bool) {
	if o == nil || IsNil(o.Safe) {
		return nil, false
	}
	return o.Safe, true
}

// HasSafe returns a boolean if a field has been set.
func (o *ApiAPIExportInput) HasSafe() bool {
	if o != nil && !IsNil(o.Safe) {
		return true
	}

	return false
}

// SetSafe gets a reference to the given bool and assigns it to the Safe field.
func (o *ApiAPIExportInput) SetSafe(v bool) {
	o.Safe = &v
}

func (o ApiAPIExportInput) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ApiAPIExportInput) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Safe) {
		toSerialize["safe"] = o.Safe
	}
	return toSerialize, nil
}

type NullableApiAPIExportInput struct {
	value *ApiAPIExportInput
	isSet bool
}

func (v NullableApiAPIExportInput) Get() *ApiAPIExportInput {
	return v.value
}

func (v *NullableApiAPIExportInput) Set(val *ApiAPIExportInput) {
	v.value = val
	v.isSet = true
}

func (v NullableApiAPIExportInput) IsSet() bool {
	return v.isSet
}

func (v *NullableApiAPIExportInput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableApiAPIExportInput(val *ApiAPIExportInput) *NullableApiAPIExportInput {
	return &NullableApiAPIExportInput{value: val, isSet: true}
}

func (v NullableApiAPIExportInput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableApiAPIExportInput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
