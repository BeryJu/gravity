/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.26.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
)

// checks if the DhcpScopeDNS type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &DhcpScopeDNS{}

// DhcpScopeDNS struct for DhcpScopeDNS
type DhcpScopeDNS struct {
	AddZoneInHostname *bool    `json:"addZoneInHostname,omitempty"`
	Search            []string `json:"search,omitempty"`
	Zone              *string  `json:"zone,omitempty"`
}

// NewDhcpScopeDNS instantiates a new DhcpScopeDNS object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDhcpScopeDNS() *DhcpScopeDNS {
	this := DhcpScopeDNS{}
	return &this
}

// NewDhcpScopeDNSWithDefaults instantiates a new DhcpScopeDNS object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDhcpScopeDNSWithDefaults() *DhcpScopeDNS {
	this := DhcpScopeDNS{}
	return &this
}

// GetAddZoneInHostname returns the AddZoneInHostname field value if set, zero value otherwise.
func (o *DhcpScopeDNS) GetAddZoneInHostname() bool {
	if o == nil || IsNil(o.AddZoneInHostname) {
		var ret bool
		return ret
	}
	return *o.AddZoneInHostname
}

// GetAddZoneInHostnameOk returns a tuple with the AddZoneInHostname field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DhcpScopeDNS) GetAddZoneInHostnameOk() (*bool, bool) {
	if o == nil || IsNil(o.AddZoneInHostname) {
		return nil, false
	}
	return o.AddZoneInHostname, true
}

// HasAddZoneInHostname returns a boolean if a field has been set.
func (o *DhcpScopeDNS) HasAddZoneInHostname() bool {
	if o != nil && !IsNil(o.AddZoneInHostname) {
		return true
	}

	return false
}

// SetAddZoneInHostname gets a reference to the given bool and assigns it to the AddZoneInHostname field.
func (o *DhcpScopeDNS) SetAddZoneInHostname(v bool) {
	o.AddZoneInHostname = &v
}

// GetSearch returns the Search field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *DhcpScopeDNS) GetSearch() []string {
	if o == nil {
		var ret []string
		return ret
	}
	return o.Search
}

// GetSearchOk returns a tuple with the Search field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *DhcpScopeDNS) GetSearchOk() ([]string, bool) {
	if o == nil || IsNil(o.Search) {
		return nil, false
	}
	return o.Search, true
}

// HasSearch returns a boolean if a field has been set.
func (o *DhcpScopeDNS) HasSearch() bool {
	if o != nil && IsNil(o.Search) {
		return true
	}

	return false
}

// SetSearch gets a reference to the given []string and assigns it to the Search field.
func (o *DhcpScopeDNS) SetSearch(v []string) {
	o.Search = v
}

// GetZone returns the Zone field value if set, zero value otherwise.
func (o *DhcpScopeDNS) GetZone() string {
	if o == nil || IsNil(o.Zone) {
		var ret string
		return ret
	}
	return *o.Zone
}

// GetZoneOk returns a tuple with the Zone field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DhcpScopeDNS) GetZoneOk() (*string, bool) {
	if o == nil || IsNil(o.Zone) {
		return nil, false
	}
	return o.Zone, true
}

// HasZone returns a boolean if a field has been set.
func (o *DhcpScopeDNS) HasZone() bool {
	if o != nil && !IsNil(o.Zone) {
		return true
	}

	return false
}

// SetZone gets a reference to the given string and assigns it to the Zone field.
func (o *DhcpScopeDNS) SetZone(v string) {
	o.Zone = &v
}

func (o DhcpScopeDNS) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o DhcpScopeDNS) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.AddZoneInHostname) {
		toSerialize["addZoneInHostname"] = o.AddZoneInHostname
	}
	if o.Search != nil {
		toSerialize["search"] = o.Search
	}
	if !IsNil(o.Zone) {
		toSerialize["zone"] = o.Zone
	}
	return toSerialize, nil
}

type NullableDhcpScopeDNS struct {
	value *DhcpScopeDNS
	isSet bool
}

func (v NullableDhcpScopeDNS) Get() *DhcpScopeDNS {
	return v.value
}

func (v *NullableDhcpScopeDNS) Set(val *DhcpScopeDNS) {
	v.value = val
	v.isSet = true
}

func (v NullableDhcpScopeDNS) IsSet() bool {
	return v.isSet
}

func (v *NullableDhcpScopeDNS) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDhcpScopeDNS(val *DhcpScopeDNS) *NullableDhcpScopeDNS {
	return &NullableDhcpScopeDNS{value: val, isSet: true}
}

func (v NullableDhcpScopeDNS) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDhcpScopeDNS) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
