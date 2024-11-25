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

// checks if the ApiAPILogMessages type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ApiAPILogMessages{}

// ApiAPILogMessages struct for ApiAPILogMessages
type ApiAPILogMessages struct {
	Messages []ApiAPILogMessage `json:"messages,omitempty"`
}

// NewApiAPILogMessages instantiates a new ApiAPILogMessages object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiAPILogMessages() *ApiAPILogMessages {
	this := ApiAPILogMessages{}
	return &this
}

// NewApiAPILogMessagesWithDefaults instantiates a new ApiAPILogMessages object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiAPILogMessagesWithDefaults() *ApiAPILogMessages {
	this := ApiAPILogMessages{}
	return &this
}

// GetMessages returns the Messages field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *ApiAPILogMessages) GetMessages() []ApiAPILogMessage {
	if o == nil {
		var ret []ApiAPILogMessage
		return ret
	}
	return o.Messages
}

// GetMessagesOk returns a tuple with the Messages field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *ApiAPILogMessages) GetMessagesOk() ([]ApiAPILogMessage, bool) {
	if o == nil || IsNil(o.Messages) {
		return nil, false
	}
	return o.Messages, true
}

// HasMessages returns a boolean if a field has been set.
func (o *ApiAPILogMessages) HasMessages() bool {
	if o != nil && IsNil(o.Messages) {
		return true
	}

	return false
}

// SetMessages gets a reference to the given []ApiAPILogMessage and assigns it to the Messages field.
func (o *ApiAPILogMessages) SetMessages(v []ApiAPILogMessage) {
	o.Messages = v
}

func (o ApiAPILogMessages) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ApiAPILogMessages) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if o.Messages != nil {
		toSerialize["messages"] = o.Messages
	}
	return toSerialize, nil
}

type NullableApiAPILogMessages struct {
	value *ApiAPILogMessages
	isSet bool
}

func (v NullableApiAPILogMessages) Get() *ApiAPILogMessages {
	return v.value
}

func (v *NullableApiAPILogMessages) Set(val *ApiAPILogMessages) {
	v.value = val
	v.isSet = true
}

func (v NullableApiAPILogMessages) IsSet() bool {
	return v.isSet
}

func (v *NullableApiAPILogMessages) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableApiAPILogMessages(val *ApiAPILogMessages) *NullableApiAPILogMessages {
	return &NullableApiAPILogMessages{value: val, isSet: true}
}

func (v NullableApiAPILogMessages) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableApiAPILogMessages) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
