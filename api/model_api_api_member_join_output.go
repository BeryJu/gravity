/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.4.7
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
)

// ApiAPIMemberJoinOutput struct for ApiAPIMemberJoinOutput
type ApiAPIMemberJoinOutput struct {
	EtcdInitialCluster *string `json:"etcdInitialCluster,omitempty"`
}

// NewApiAPIMemberJoinOutput instantiates a new ApiAPIMemberJoinOutput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiAPIMemberJoinOutput() *ApiAPIMemberJoinOutput {
	this := ApiAPIMemberJoinOutput{}
	return &this
}

// NewApiAPIMemberJoinOutputWithDefaults instantiates a new ApiAPIMemberJoinOutput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiAPIMemberJoinOutputWithDefaults() *ApiAPIMemberJoinOutput {
	this := ApiAPIMemberJoinOutput{}
	return &this
}

// GetEtcdInitialCluster returns the EtcdInitialCluster field value if set, zero value otherwise.
func (o *ApiAPIMemberJoinOutput) GetEtcdInitialCluster() string {
	if o == nil || o.EtcdInitialCluster == nil {
		var ret string
		return ret
	}
	return *o.EtcdInitialCluster
}

// GetEtcdInitialClusterOk returns a tuple with the EtcdInitialCluster field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAPIMemberJoinOutput) GetEtcdInitialClusterOk() (*string, bool) {
	if o == nil || o.EtcdInitialCluster == nil {
		return nil, false
	}
	return o.EtcdInitialCluster, true
}

// HasEtcdInitialCluster returns a boolean if a field has been set.
func (o *ApiAPIMemberJoinOutput) HasEtcdInitialCluster() bool {
	if o != nil && o.EtcdInitialCluster != nil {
		return true
	}

	return false
}

// SetEtcdInitialCluster gets a reference to the given string and assigns it to the EtcdInitialCluster field.
func (o *ApiAPIMemberJoinOutput) SetEtcdInitialCluster(v string) {
	o.EtcdInitialCluster = &v
}

func (o ApiAPIMemberJoinOutput) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.EtcdInitialCluster != nil {
		toSerialize["etcdInitialCluster"] = o.EtcdInitialCluster
	}
	return json.Marshal(toSerialize)
}

type NullableApiAPIMemberJoinOutput struct {
	value *ApiAPIMemberJoinOutput
	isSet bool
}

func (v NullableApiAPIMemberJoinOutput) Get() *ApiAPIMemberJoinOutput {
	return v.value
}

func (v *NullableApiAPIMemberJoinOutput) Set(val *ApiAPIMemberJoinOutput) {
	v.value = val
	v.isSet = true
}

func (v NullableApiAPIMemberJoinOutput) IsSet() bool {
	return v.isSet
}

func (v *NullableApiAPIMemberJoinOutput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableApiAPIMemberJoinOutput(val *ApiAPIMemberJoinOutput) *NullableApiAPIMemberJoinOutput {
	return &NullableApiAPIMemberJoinOutput{value: val, isSet: true}
}

func (v NullableApiAPIMemberJoinOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableApiAPIMemberJoinOutput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
