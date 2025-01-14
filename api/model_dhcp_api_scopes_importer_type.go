/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.24.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
	"fmt"
)

// DhcpAPIScopesImporterType the model 'DhcpAPIScopesImporterType'
type DhcpAPIScopesImporterType string

// List of DhcpAPIScopesImporterType
const (
	DHCPAPISCOPESIMPORTERTYPE_MS_DHCP DhcpAPIScopesImporterType = "ms_dhcp"
)

// All allowed values of DhcpAPIScopesImporterType enum
var AllowedDhcpAPIScopesImporterTypeEnumValues = []DhcpAPIScopesImporterType{
	"ms_dhcp",
}

func (v *DhcpAPIScopesImporterType) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := DhcpAPIScopesImporterType(value)
	for _, existing := range AllowedDhcpAPIScopesImporterTypeEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	return fmt.Errorf("%+v is not a valid DhcpAPIScopesImporterType", value)
}

// NewDhcpAPIScopesImporterTypeFromValue returns a pointer to a valid DhcpAPIScopesImporterType
// for the value passed as argument, or an error if the value passed is not allowed by the enum
func NewDhcpAPIScopesImporterTypeFromValue(v string) (*DhcpAPIScopesImporterType, error) {
	ev := DhcpAPIScopesImporterType(v)
	if ev.IsValid() {
		return &ev, nil
	} else {
		return nil, fmt.Errorf("invalid value '%v' for DhcpAPIScopesImporterType: valid values are %v", v, AllowedDhcpAPIScopesImporterTypeEnumValues)
	}
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v DhcpAPIScopesImporterType) IsValid() bool {
	for _, existing := range AllowedDhcpAPIScopesImporterTypeEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to DhcpAPIScopesImporterType value
func (v DhcpAPIScopesImporterType) Ptr() *DhcpAPIScopesImporterType {
	return &v
}

type NullableDhcpAPIScopesImporterType struct {
	value *DhcpAPIScopesImporterType
	isSet bool
}

func (v NullableDhcpAPIScopesImporterType) Get() *DhcpAPIScopesImporterType {
	return v.value
}

func (v *NullableDhcpAPIScopesImporterType) Set(val *DhcpAPIScopesImporterType) {
	v.value = val
	v.isSet = true
}

func (v NullableDhcpAPIScopesImporterType) IsSet() bool {
	return v.isSet
}

func (v *NullableDhcpAPIScopesImporterType) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDhcpAPIScopesImporterType(val *DhcpAPIScopesImporterType) *NullableDhcpAPIScopesImporterType {
	return &NullableDhcpAPIScopesImporterType{value: val, isSet: true}
}

func (v NullableDhcpAPIScopesImporterType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDhcpAPIScopesImporterType) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
