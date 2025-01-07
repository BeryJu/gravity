/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.21.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
	"time"
)

// checks if the ApiAPILogMessage type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ApiAPILogMessage{}

// ApiAPILogMessage struct for ApiAPILogMessage
type ApiAPILogMessage struct {
	Level   *string    `json:"level,omitempty"`
	Logger  *string    `json:"logger,omitempty"`
	Message *string    `json:"message,omitempty"`
	Node    *string    `json:"node,omitempty"`
	Time    *time.Time `json:"time,omitempty"`
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

// GetLevel returns the Level field value if set, zero value otherwise.
func (o *ApiAPILogMessage) GetLevel() string {
	if o == nil || IsNil(o.Level) {
		var ret string
		return ret
	}
	return *o.Level
}

// GetLevelOk returns a tuple with the Level field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAPILogMessage) GetLevelOk() (*string, bool) {
	if o == nil || IsNil(o.Level) {
		return nil, false
	}
	return o.Level, true
}

// HasLevel returns a boolean if a field has been set.
func (o *ApiAPILogMessage) HasLevel() bool {
	if o != nil && !IsNil(o.Level) {
		return true
	}

	return false
}

// SetLevel gets a reference to the given string and assigns it to the Level field.
func (o *ApiAPILogMessage) SetLevel(v string) {
	o.Level = &v
}

// GetLogger returns the Logger field value if set, zero value otherwise.
func (o *ApiAPILogMessage) GetLogger() string {
	if o == nil || IsNil(o.Logger) {
		var ret string
		return ret
	}
	return *o.Logger
}

// GetLoggerOk returns a tuple with the Logger field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAPILogMessage) GetLoggerOk() (*string, bool) {
	if o == nil || IsNil(o.Logger) {
		return nil, false
	}
	return o.Logger, true
}

// HasLogger returns a boolean if a field has been set.
func (o *ApiAPILogMessage) HasLogger() bool {
	if o != nil && !IsNil(o.Logger) {
		return true
	}

	return false
}

// SetLogger gets a reference to the given string and assigns it to the Logger field.
func (o *ApiAPILogMessage) SetLogger(v string) {
	o.Logger = &v
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

// GetTime returns the Time field value if set, zero value otherwise.
func (o *ApiAPILogMessage) GetTime() time.Time {
	if o == nil || IsNil(o.Time) {
		var ret time.Time
		return ret
	}
	return *o.Time
}

// GetTimeOk returns a tuple with the Time field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiAPILogMessage) GetTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.Time) {
		return nil, false
	}
	return o.Time, true
}

// HasTime returns a boolean if a field has been set.
func (o *ApiAPILogMessage) HasTime() bool {
	if o != nil && !IsNil(o.Time) {
		return true
	}

	return false
}

// SetTime gets a reference to the given time.Time and assigns it to the Time field.
func (o *ApiAPILogMessage) SetTime(v time.Time) {
	o.Time = &v
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
	if !IsNil(o.Level) {
		toSerialize["level"] = o.Level
	}
	if !IsNil(o.Logger) {
		toSerialize["logger"] = o.Logger
	}
	if !IsNil(o.Message) {
		toSerialize["message"] = o.Message
	}
	if !IsNil(o.Node) {
		toSerialize["node"] = o.Node
	}
	if !IsNil(o.Time) {
		toSerialize["time"] = o.Time
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
