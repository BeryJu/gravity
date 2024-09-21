/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.9.1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
)

// checks if the ApiAPILogMessage type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ApiAPILogMessage{}

// ApiAPILogMessage struct for ApiAPILogMessage
type ApiAPILogMessage struct {
	Message *string `json:"message,omitempty"`
	Node    *string `json:"node,omitempty"`
}

// NewApiAPILogMessage instantiates a new ApiAPILogMessage object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiAPILogMessage() *ApiAPILogMessage {
	this := ApiAPILogMessage{}
	return &this
}

// NewApiAPILogMessageWithDefaults instantiates a new ApiAPILogMessage object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiAPILogMessageWithDefaults() *ApiAPILogMessage {
	this := ApiAPILogMessage{}
	return &this
}

// GetMessage returns the Message field value if set, zero value otherwise.
func (o *ApiAPILogMessage) GetMessage() string {
	if o == nil || IsNil(o.Message) {
		var ret string
		return ret
	}
	return *o.Message
}

// GetMessageOk returns a tuple with the Message field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAPILogMessage) GetMessageOk() (*string, bool) {
	if o == nil || IsNil(o.Message) {
		return nil, false
	}
	return o.Message, true
}

// HasMessage returns a boolean if a field has been set.
func (o *ApiAPILogMessage) HasMessage() bool {
	if o != nil && !IsNil(o.Message) {
		return true
	}

	return false
}

// SetMessage gets a reference to the given string and assigns it to the Message field.
func (o *ApiAPILogMessage) SetMessage(v string) {
	o.Message = &v
}

// GetNode returns the Node field value if set, zero value otherwise.
func (o *ApiAPILogMessage) GetNode() string {
	if o == nil || IsNil(o.Node) {
		var ret string
		return ret
	}
	return *o.Node
}

// GetNodeOk returns a tuple with the Node field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAPILogMessage) GetNodeOk() (*string, bool) {
	if o == nil || IsNil(o.Node) {
		return nil, false
	}
	return o.Node, true
}

// HasNode returns a boolean if a field has been set.
func (o *ApiAPILogMessage) HasNode() bool {
	if o != nil && !IsNil(o.Node) {
		return true
	}

	return false
}

// SetNode gets a reference to the given string and assigns it to the Node field.
func (o *ApiAPILogMessage) SetNode(v string) {
	o.Node = &v
}

func (o ApiAPILogMessage) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ApiAPILogMessage) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Message) {
		toSerialize["message"] = o.Message
	}
	if !IsNil(o.Node) {
		toSerialize["node"] = o.Node
	}
	return toSerialize, nil
}

type NullableApiAPILogMessage struct {
	value *ApiAPILogMessage
	isSet bool
}

func (v NullableApiAPILogMessage) Get() *ApiAPILogMessage {
	return v.value
}

func (v *NullableApiAPILogMessage) Set(val *ApiAPILogMessage) {
	v.value = val
	v.isSet = true
}

func (v NullableApiAPILogMessage) IsSet() bool {
	return v.isSet
}

func (v *NullableApiAPILogMessage) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableApiAPILogMessage(val *ApiAPILogMessage) *NullableApiAPILogMessage {
	return &NullableApiAPILogMessage{value: val, isSet: true}
}

func (v NullableApiAPILogMessage) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableApiAPILogMessage) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
