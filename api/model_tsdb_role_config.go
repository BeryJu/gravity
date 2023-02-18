/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.4.9
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
)

// TsdbRoleConfig struct for TsdbRoleConfig
type TsdbRoleConfig struct {
	Enabled *bool  `json:"enabled,omitempty"`
	Expire  *int32 `json:"expire,omitempty"`
	Scrape  *int32 `json:"scrape,omitempty"`
}

// NewTsdbRoleConfig instantiates a new TsdbRoleConfig object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTsdbRoleConfig() *TsdbRoleConfig {
	this := TsdbRoleConfig{}
	return &this
}

// NewTsdbRoleConfigWithDefaults instantiates a new TsdbRoleConfig object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTsdbRoleConfigWithDefaults() *TsdbRoleConfig {
	this := TsdbRoleConfig{}
	return &this
}

// GetEnabled returns the Enabled field value if set, zero value otherwise.
func (o *TsdbRoleConfig) GetEnabled() bool {
	if o == nil || o.Enabled == nil {
		var ret bool
		return ret
	}
	return *o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TsdbRoleConfig) GetEnabledOk() (*bool, bool) {
	if o == nil || o.Enabled == nil {
		return nil, false
	}
	return o.Enabled, true
}

// HasEnabled returns a boolean if a field has been set.
func (o *TsdbRoleConfig) HasEnabled() bool {
	if o != nil && o.Enabled != nil {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given bool and assigns it to the Enabled field.
func (o *TsdbRoleConfig) SetEnabled(v bool) {
	o.Enabled = &v
}

// GetExpire returns the Expire field value if set, zero value otherwise.
func (o *TsdbRoleConfig) GetExpire() int32 {
	if o == nil || o.Expire == nil {
		var ret int32
		return ret
	}
	return *o.Expire
}

// GetExpireOk returns a tuple with the Expire field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TsdbRoleConfig) GetExpireOk() (*int32, bool) {
	if o == nil || o.Expire == nil {
		return nil, false
	}
	return o.Expire, true
}

// HasExpire returns a boolean if a field has been set.
func (o *TsdbRoleConfig) HasExpire() bool {
	if o != nil && o.Expire != nil {
		return true
	}

	return false
}

// SetExpire gets a reference to the given int32 and assigns it to the Expire field.
func (o *TsdbRoleConfig) SetExpire(v int32) {
	o.Expire = &v
}

// GetScrape returns the Scrape field value if set, zero value otherwise.
func (o *TsdbRoleConfig) GetScrape() int32 {
	if o == nil || o.Scrape == nil {
		var ret int32
		return ret
	}
	return *o.Scrape
}

// GetScrapeOk returns a tuple with the Scrape field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TsdbRoleConfig) GetScrapeOk() (*int32, bool) {
	if o == nil || o.Scrape == nil {
		return nil, false
	}
	return o.Scrape, true
}

// HasScrape returns a boolean if a field has been set.
func (o *TsdbRoleConfig) HasScrape() bool {
	if o != nil && o.Scrape != nil {
		return true
	}

	return false
}

// SetScrape gets a reference to the given int32 and assigns it to the Scrape field.
func (o *TsdbRoleConfig) SetScrape(v int32) {
	o.Scrape = &v
}

func (o TsdbRoleConfig) MarshalJSON() ([]byte, error) {
	toSerialize := map[string]interface{}{}
	if o.Enabled != nil {
		toSerialize["enabled"] = o.Enabled
	}
	if o.Expire != nil {
		toSerialize["expire"] = o.Expire
	}
	if o.Scrape != nil {
		toSerialize["scrape"] = o.Scrape
	}
	return json.Marshal(toSerialize)
}

type NullableTsdbRoleConfig struct {
	value *TsdbRoleConfig
	isSet bool
}

func (v NullableTsdbRoleConfig) Get() *TsdbRoleConfig {
	return v.value
}

func (v *NullableTsdbRoleConfig) Set(val *TsdbRoleConfig) {
	v.value = val
	v.isSet = true
}

func (v NullableTsdbRoleConfig) IsSet() bool {
	return v.isSet
}

func (v *NullableTsdbRoleConfig) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTsdbRoleConfig(val *TsdbRoleConfig) *NullableTsdbRoleConfig {
	return &NullableTsdbRoleConfig{value: val, isSet: true}
}

func (v NullableTsdbRoleConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTsdbRoleConfig) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
