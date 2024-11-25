/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.17.1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
)

// checks if the DnsAPIZonesPutInput type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &DnsAPIZonesPutInput{}

// DnsAPIZonesPutInput struct for DnsAPIZonesPutInput
type DnsAPIZonesPutInput struct {
	Authoritative  bool                     `json:"authoritative"`
	DefaultTTL     int32                    `json:"defaultTTL"`
	HandlerConfigs []map[string]interface{} `json:"handlerConfigs"`
	Hook           string                   `json:"hook"`
}

// NewDnsAPIZonesPutInput instantiates a new DnsAPIZonesPutInput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDnsAPIZonesPutInput(authoritative bool, defaultTTL int32, handlerConfigs []map[string]interface{}, hook string) *DnsAPIZonesPutInput {
	this := DnsAPIZonesPutInput{}
	this.Authoritative = authoritative
	this.DefaultTTL = defaultTTL
	this.HandlerConfigs = handlerConfigs
	this.Hook = hook
	return &this
}

// NewDnsAPIZonesPutInputWithDefaults instantiates a new DnsAPIZonesPutInput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDnsAPIZonesPutInputWithDefaults() *DnsAPIZonesPutInput {
	this := DnsAPIZonesPutInput{}
	return &this
}

// GetAuthoritative returns the Authoritative field value
func (o *DnsAPIZonesPutInput) GetAuthoritative() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.Authoritative
}

// GetAuthoritativeOk returns a tuple with the Authoritative field value
// and a boolean to check if the value has been set.
func (o *DnsAPIZonesPutInput) GetAuthoritativeOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Authoritative, true
}

// SetAuthoritative sets field value
func (o *DnsAPIZonesPutInput) SetAuthoritative(v bool) {
	o.Authoritative = v
}

// GetDefaultTTL returns the DefaultTTL field value
func (o *DnsAPIZonesPutInput) GetDefaultTTL() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.DefaultTTL
}

// GetDefaultTTLOk returns a tuple with the DefaultTTL field value
// and a boolean to check if the value has been set.
func (o *DnsAPIZonesPutInput) GetDefaultTTLOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DefaultTTL, true
}

// SetDefaultTTL sets field value
func (o *DnsAPIZonesPutInput) SetDefaultTTL(v int32) {
	o.DefaultTTL = v
}

// GetHandlerConfigs returns the HandlerConfigs field value
// If the value is explicit nil, the zero value for []map[string]interface{} will be returned
func (o *DnsAPIZonesPutInput) GetHandlerConfigs() []map[string]interface{} {
	if o == nil {
		var ret []map[string]interface{}
		return ret
	}

	return o.HandlerConfigs
}

// GetHandlerConfigsOk returns a tuple with the HandlerConfigs field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *DnsAPIZonesPutInput) GetHandlerConfigsOk() ([]map[string]interface{}, bool) {
	if o == nil || IsNil(o.HandlerConfigs) {
		return nil, false
	}
	return o.HandlerConfigs, true
}

// SetHandlerConfigs sets field value
func (o *DnsAPIZonesPutInput) SetHandlerConfigs(v []map[string]interface{}) {
	o.HandlerConfigs = v
}

// GetHook returns the Hook field value
func (o *DnsAPIZonesPutInput) GetHook() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Hook
}

// GetHookOk returns a tuple with the Hook field value
// and a boolean to check if the value has been set.
func (o *DnsAPIZonesPutInput) GetHookOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Hook, true
}

// SetHook sets field value
func (o *DnsAPIZonesPutInput) SetHook(v string) {
	o.Hook = v
}

func (o DnsAPIZonesPutInput) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o DnsAPIZonesPutInput) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["authoritative"] = o.Authoritative
	toSerialize["defaultTTL"] = o.DefaultTTL
	if o.HandlerConfigs != nil {
		toSerialize["handlerConfigs"] = o.HandlerConfigs
	}
	toSerialize["hook"] = o.Hook
	return toSerialize, nil
}

type NullableDnsAPIZonesPutInput struct {
	value *DnsAPIZonesPutInput
	isSet bool
}

func (v NullableDnsAPIZonesPutInput) Get() *DnsAPIZonesPutInput {
	return v.value
}

func (v *NullableDnsAPIZonesPutInput) Set(val *DnsAPIZonesPutInput) {
	v.value = val
	v.isSet = true
}

func (v NullableDnsAPIZonesPutInput) IsSet() bool {
	return v.isSet
}

func (v *NullableDnsAPIZonesPutInput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDnsAPIZonesPutInput(val *DnsAPIZonesPutInput) *NullableDnsAPIZonesPutInput {
	return &NullableDnsAPIZonesPutInput{value: val, isSet: true}
}

func (v NullableDnsAPIZonesPutInput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDnsAPIZonesPutInput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
