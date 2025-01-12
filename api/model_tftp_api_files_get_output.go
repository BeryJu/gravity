/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.23.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
)

// checks if the TftpAPIFilesGetOutput type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TftpAPIFilesGetOutput{}

// TftpAPIFilesGetOutput struct for TftpAPIFilesGetOutput
type TftpAPIFilesGetOutput struct {
	Files []TftpAPIFile `json:"files"`
}

// NewTftpAPIFilesGetOutput instantiates a new TftpAPIFilesGetOutput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTftpAPIFilesGetOutput(files []TftpAPIFile) *TftpAPIFilesGetOutput {
	this := TftpAPIFilesGetOutput{}
	this.Files = files
	return &this
}

// NewTftpAPIFilesGetOutputWithDefaults instantiates a new TftpAPIFilesGetOutput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTftpAPIFilesGetOutputWithDefaults() *TftpAPIFilesGetOutput {
	this := TftpAPIFilesGetOutput{}
	return &this
}

// GetFiles returns the Files field value
// If the value is explicit nil, the zero value for []TftpAPIFile will be returned
func (o *TftpAPIFilesGetOutput) GetFiles() []TftpAPIFile {
	if o == nil {
		var ret []TftpAPIFile
		return ret
	}

	return o.Files
}

// GetFilesOk returns a tuple with the Files field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *TftpAPIFilesGetOutput) GetFilesOk() ([]TftpAPIFile, bool) {
	if o == nil || IsNil(o.Files) {
		return nil, false
	}
	return o.Files, true
}

// SetFiles sets field value
func (o *TftpAPIFilesGetOutput) SetFiles(v []TftpAPIFile) {
	o.Files = v
}

func (o TftpAPIFilesGetOutput) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TftpAPIFilesGetOutput) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if o.Files != nil {
		toSerialize["files"] = o.Files
	}
	return toSerialize, nil
}

type NullableTftpAPIFilesGetOutput struct {
	value *TftpAPIFilesGetOutput
	isSet bool
}

func (v NullableTftpAPIFilesGetOutput) Get() *TftpAPIFilesGetOutput {
	return v.value
}

func (v *NullableTftpAPIFilesGetOutput) Set(val *TftpAPIFilesGetOutput) {
	v.value = val
	v.isSet = true
}

func (v NullableTftpAPIFilesGetOutput) IsSet() bool {
	return v.isSet
}

func (v *NullableTftpAPIFilesGetOutput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTftpAPIFilesGetOutput(val *TftpAPIFilesGetOutput) *NullableTftpAPIFilesGetOutput {
	return &NullableTftpAPIFilesGetOutput{value: val, isSet: true}
}

func (v NullableTftpAPIFilesGetOutput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTftpAPIFilesGetOutput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
