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

// ApiAPIToolPortmapOutputPort struct for ApiAPIToolPortmapOutputPort
type ApiAPIToolPortmapOutputPort struct {
	Name     *string `json:"name,omitempty"`
	Port     *int32  `json:"port,omitempty"`
	Protocol *string `json:"protocol,omitempty"`
	Reason   *string `json:"reason,omitempty"`
}

// NewApiAPIToolPortmapOutputPort instantiates a new ApiAPIToolPortmapOutputPort object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiAPIToolPortmapOutputPort() *ApiAPIToolPortmapOutputPort {
	this := ApiAPIToolPortmapOutputPort{}
	return &this
}

// NewApiAPIToolPortmapOutputPortWithDefaults instantiates a new ApiAPIToolPortmapOutputPort object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiAPIToolPortmapOutputPortWithDefaults() *ApiAPIToolPortmapOutputPort {
	this := ApiAPIToolPortmapOutputPort{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *ApiAPIToolPortmapOutputPort) GetName() string {
	if o == nil || o.Name == nil {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAPIToolPortmapOutputPort) GetNameOk() (*string, bool) {
	if o == nil || o.Name == nil {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *ApiAPIToolPortmapOutputPort) HasName() bool {
	if o != nil && o.Name != nil {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *ApiAPIToolPortmapOutputPort) SetName(v string) {
	o.Name = &v
}

// GetPort returns the Port field value if set, zero value otherwise.
func (o *ApiAPIToolPortmapOutputPort) GetPort() int32 {
	if o == nil || o.Port == nil {
		var ret int32
		return ret
	}
	return *o.Port
}

// GetPortOk returns a tuple with the Port field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAPIToolPortmapOutputPort) GetPortOk() (*int32, bool) {
	if o == nil || o.Port == nil {
		return nil, false
	}
	return o.Port, true
}

// HasPort returns a boolean if a field has been set.
func (o *ApiAPIToolPortmapOutputPort) HasPort() bool {
	if o != nil && o.Port != nil {
		return true
	}

	return false
}

// SetPort gets a reference to the given int32 and assigns it to the Port field.
func (o *ApiAPIToolPortmapOutputPort) SetPort(v int32) {
	o.Port = &v
}

// GetProtocol returns the Protocol field value if set, zero value otherwise.
func (o *ApiAPIToolPortmapOutputPort) GetProtocol() string {
	if o == nil || o.Protocol == nil {
		var ret string
		return ret
	}
	return *o.Protocol
}

// GetProtocolOk returns a tuple with the Protocol field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAPIToolPortmapOutputPort) GetProtocolOk() (*string, bool) {
	if o == nil || o.Protocol == nil {
		return nil, false
	}
	return o.Protocol, true
}

// HasProtocol returns a boolean if a field has been set.
func (o *ApiAPIToolPortmapOutputPort) HasProtocol() bool {
	if o != nil && o.Protocol != nil {
		return true
	}

	return false
}

// SetProtocol gets a reference to the given string and assigns it to the Protocol field.
func (o *ApiAPIToolPortmapOutputPort) SetProtocol(v string) {
	o.Protocol = &v
}

// GetReason returns the Reason field value if set, zero value otherwise.
func (o *ApiAPIToolPortmapOutputPort) GetReason() string {
	if o == nil || o.Reason == nil {
		var ret string
		return ret
	}
	return *o.Reason
}

// GetReasonOk returns a tuple with the Reason field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAPIToolPortmapOutputPort) GetReasonOk() (*string, bool) {
	if o == nil || o.Reason == nil {
		return nil, false
	}
	return o.Reason, true
}

// HasReason returns a boolean if a field has been set.
func (o *ApiAPIToolPortmapOutputPort) HasReason() bool {
	if o != nil && o.Reason != nil {
		return true
	}

	return false
}

// SetReason gets a reference to the given string and assigns it to the Reason field.
func (o *ApiAPIToolPortmapOutputPort) SetReason(v string) {
	o.Reason = &v
}

func (o ApiAPIToolPortmapOutputPort) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Name != nil {
		toSerialize["name"] = o.Name
	}
	if o.Port != nil {
		toSerialize["port"] = o.Port
	}
	if o.Protocol != nil {
		toSerialize["protocol"] = o.Protocol
	}
	if o.Reason != nil {
		toSerialize["reason"] = o.Reason
	}
	return json.Marshal(toSerialize)
}

type NullableApiAPIToolPortmapOutputPort struct {
	value *ApiAPIToolPortmapOutputPort
	isSet bool
}

func (v NullableApiAPIToolPortmapOutputPort) Get() *ApiAPIToolPortmapOutputPort {
	return v.value
}

func (v *NullableApiAPIToolPortmapOutputPort) Set(val *ApiAPIToolPortmapOutputPort) {
	v.value = val
	v.isSet = true
}

func (v NullableApiAPIToolPortmapOutputPort) IsSet() bool {
	return v.isSet
}

func (v *NullableApiAPIToolPortmapOutputPort) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableApiAPIToolPortmapOutputPort(val *ApiAPIToolPortmapOutputPort) *NullableApiAPIToolPortmapOutputPort {
	return &NullableApiAPIToolPortmapOutputPort{value: val, isSet: true}
}

func (v NullableApiAPIToolPortmapOutputPort) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableApiAPIToolPortmapOutputPort) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
