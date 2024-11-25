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

// checks if the TypesAPIMetricsGetOutput type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TypesAPIMetricsGetOutput{}

// TypesAPIMetricsGetOutput struct for TypesAPIMetricsGetOutput
type TypesAPIMetricsGetOutput struct {
	Records []TypesAPIMetricsRecord `json:"records"`
}

// NewTypesAPIMetricsGetOutput instantiates a new TypesAPIMetricsGetOutput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTypesAPIMetricsGetOutput(records []TypesAPIMetricsRecord) *TypesAPIMetricsGetOutput {
	this := TypesAPIMetricsGetOutput{}
	this.Records = records
	return &this
}

// NewTypesAPIMetricsGetOutputWithDefaults instantiates a new TypesAPIMetricsGetOutput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTypesAPIMetricsGetOutputWithDefaults() *TypesAPIMetricsGetOutput {
	this := TypesAPIMetricsGetOutput{}
	return &this
}

// GetRecords returns the Records field value
// If the value is explicit nil, the zero value for []TypesAPIMetricsRecord will be returned
func (o *TypesAPIMetricsGetOutput) GetRecords() []TypesAPIMetricsRecord {
	if o == nil {
		var ret []TypesAPIMetricsRecord
		return ret
	}

	return o.Records
}

// GetRecordsOk returns a tuple with the Records field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TypesAPIMetricsGetOutput) GetRecordsOk() ([]TypesAPIMetricsRecord, bool) {
	if o == nil || IsNil(o.Records) {
		return nil, false
	}
	return o.Records, true
}

// SetRecords sets field value
func (o *TypesAPIMetricsGetOutput) SetRecords(v []TypesAPIMetricsRecord) {
	o.Records = v
}

func (o TypesAPIMetricsGetOutput) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TypesAPIMetricsGetOutput) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if o.Records != nil {
		toSerialize["records"] = o.Records
	}
	return toSerialize, nil
}

type NullableTypesAPIMetricsGetOutput struct {
	value *TypesAPIMetricsGetOutput
	isSet bool
}

func (v NullableTypesAPIMetricsGetOutput) Get() *TypesAPIMetricsGetOutput {
	return v.value
}

func (v *NullableTypesAPIMetricsGetOutput) Set(val *TypesAPIMetricsGetOutput) {
	v.value = val
	v.isSet = true
}

func (v NullableTypesAPIMetricsGetOutput) IsSet() bool {
	return v.isSet
}

func (v *NullableTypesAPIMetricsGetOutput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTypesAPIMetricsGetOutput(val *TypesAPIMetricsGetOutput) *NullableTypesAPIMetricsGetOutput {
	return &NullableTypesAPIMetricsGetOutput{value: val, isSet: true}
}

func (v NullableTypesAPIMetricsGetOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTypesAPIMetricsGetOutput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
