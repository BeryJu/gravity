/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.27.2
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
	"time"
)

// checks if the TypesAPIMetricsRecord type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TypesAPIMetricsRecord{}

// TypesAPIMetricsRecord struct for TypesAPIMetricsRecord
type TypesAPIMetricsRecord struct {
	Keys  []string  `json:"keys"`
	Node  string    `json:"node"`
	Time  time.Time `json:"time"`
	Value int32     `json:"value"`
}

// NewTypesAPIMetricsRecord instantiates a new TypesAPIMetricsRecord object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTypesAPIMetricsRecord(keys []string, node string, time time.Time, value int32) *TypesAPIMetricsRecord {
	this := TypesAPIMetricsRecord{}
	this.Keys = keys
	this.Node = node
	this.Time = time
	this.Value = value
	return &this
}

// NewTypesAPIMetricsRecordWithDefaults instantiates a new TypesAPIMetricsRecord object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTypesAPIMetricsRecordWithDefaults() *TypesAPIMetricsRecord {
	this := TypesAPIMetricsRecord{}
	return &this
}

// GetKeys returns the Keys field value
// If the value is explicit nil, the zero value for []string will be returned
func (o *TypesAPIMetricsRecord) GetKeys() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.Keys
}

// GetKeysOk returns a tuple with the Keys field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TypesAPIMetricsRecord) GetKeysOk() ([]string, bool) {
	if o == nil || IsNil(o.Keys) {
		return nil, false
	}
	return o.Keys, true
}

// SetKeys sets field value
func (o *TypesAPIMetricsRecord) SetKeys(v []string) {
	o.Keys = v
}

// GetNode returns the Node field value
func (o *TypesAPIMetricsRecord) GetNode() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Node
}

// GetNodeOk returns a tuple with the Node field value
// and a boolean to check if the value has been set.
func (o *TypesAPIMetricsRecord) GetNodeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Node, true
}

// SetNode sets field value
func (o *TypesAPIMetricsRecord) SetNode(v string) {
	o.Node = v
}

// GetTime returns the Time field value
func (o *TypesAPIMetricsRecord) GetTime() time.Time {
	if o == nil {
		var ret time.Time
		return ret
	}

	return o.Time
}

// GetTimeOk returns a tuple with the Time field value
// and a boolean to check if the value has been set.
func (o *TypesAPIMetricsRecord) GetTimeOk() (*time.Time, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Time, true
}

// SetTime sets field value
func (o *TypesAPIMetricsRecord) SetTime(v time.Time) {
	o.Time = v
}

// GetValue returns the Value field value
func (o *TypesAPIMetricsRecord) GetValue() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Value
}

// GetValueOk returns a tuple with the Value field value
// and a boolean to check if the value has been set.
func (o *TypesAPIMetricsRecord) GetValueOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Value, true
}

// SetValue sets field value
func (o *TypesAPIMetricsRecord) SetValue(v int32) {
	o.Value = v
}

func (o TypesAPIMetricsRecord) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TypesAPIMetricsRecord) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if o.Keys != nil {
		toSerialize["keys"] = o.Keys
	}
	toSerialize["node"] = o.Node
	toSerialize["time"] = o.Time
	toSerialize["value"] = o.Value
	return toSerialize, nil
}

type NullableTypesAPIMetricsRecord struct {
	value *TypesAPIMetricsRecord
	isSet bool
}

func (v NullableTypesAPIMetricsRecord) Get() *TypesAPIMetricsRecord {
	return v.value
}

func (v *NullableTypesAPIMetricsRecord) Set(val *TypesAPIMetricsRecord) {
	v.value = val
	v.isSet = true
}

func (v NullableTypesAPIMetricsRecord) IsSet() bool {
	return v.isSet
}

func (v *NullableTypesAPIMetricsRecord) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTypesAPIMetricsRecord(val *TypesAPIMetricsRecord) *NullableTypesAPIMetricsRecord {
	return &NullableTypesAPIMetricsRecord{value: val, isSet: true}
}

func (v NullableTypesAPIMetricsRecord) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTypesAPIMetricsRecord) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
