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

// checks if the DhcpAPIScope type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &DhcpAPIScope{}

// DhcpAPIScope struct for DhcpAPIScope
type DhcpAPIScope struct {
	Default    bool                   `json:"default"`
	Dns        *DhcpScopeDNS          `json:"dns,omitempty"`
	Hook       string                 `json:"hook"`
	Ipam       map[string]string      `json:"ipam"`
	Options    []TypesDHCPOption      `json:"options"`
	Scope      string                 `json:"scope"`
	Statistics DhcpAPIScopeStatistics `json:"statistics"`
	SubnetCidr string                 `json:"subnetCidr"`
	Ttl        int32                  `json:"ttl"`
}

// NewDhcpAPIScope instantiates a new DhcpAPIScope object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDhcpAPIScope(default_ bool, hook string, ipam map[string]string, options []TypesDHCPOption, scope string, statistics DhcpAPIScopeStatistics, subnetCidr string, ttl int32) *DhcpAPIScope {
	this := DhcpAPIScope{}
	this.Default = default_
	this.Hook = hook
	this.Ipam = ipam
	this.Options = options
	this.Scope = scope
	this.Statistics = statistics
	this.SubnetCidr = subnetCidr
	this.Ttl = ttl
	return &this
}

// NewDhcpAPIScopeWithDefaults instantiates a new DhcpAPIScope object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDhcpAPIScopeWithDefaults() *DhcpAPIScope {
	this := DhcpAPIScope{}
	return &this
}

// GetDefault returns the Default field value
func (o *DhcpAPIScope) GetDefault() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.Default
}

// GetDefaultOk returns a tuple with the Default field value
// and a boolean to check if the value has been set.
func (o *DhcpAPIScope) GetDefaultOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Default, true
}

// SetDefault sets field value
func (o *DhcpAPIScope) SetDefault(v bool) {
	o.Default = v
}

// GetDns returns the Dns field value if set, zero value otherwise.
func (o *DhcpAPIScope) GetDns() DhcpScopeDNS {
	if o == nil || IsNil(o.Dns) {
		var ret DhcpScopeDNS
		return ret
	}
	return *o.Dns
}

// GetDnsOk returns a tuple with the Dns field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DhcpAPIScope) GetDnsOk() (*DhcpScopeDNS, bool) {
	if o == nil || IsNil(o.Dns) {
		return nil, false
	}
	return o.Dns, true
}

// HasDns returns a boolean if a field has been set.
func (o *DhcpAPIScope) HasDns() bool {
	if o != nil && !IsNil(o.Dns) {
		return true
	}

	return false
}

// SetDns gets a reference to the given DhcpScopeDNS and assigns it to the Dns field.
func (o *DhcpAPIScope) SetDns(v DhcpScopeDNS) {
	o.Dns = &v
}

// GetHook returns the Hook field value
func (o *DhcpAPIScope) GetHook() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Hook
}

// GetHookOk returns a tuple with the Hook field value
// and a boolean to check if the value has been set.
func (o *DhcpAPIScope) GetHookOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Hook, true
}

// SetHook sets field value
func (o *DhcpAPIScope) SetHook(v string) {
	o.Hook = v
}

// GetIpam returns the Ipam field value
// If the value is explicit nil, the zero value for map[string]string will be returned
func (o *DhcpAPIScope) GetIpam() map[string]string {
	if o == nil {
		var ret map[string]string
		return ret
	}

	return o.Ipam
}

// GetIpamOk returns a tuple with the Ipam field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *DhcpAPIScope) GetIpamOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.Ipam) {
		return nil, false
	}
	return &o.Ipam, true
}

// SetIpam sets field value
func (o *DhcpAPIScope) SetIpam(v map[string]string) {
	o.Ipam = v
}

// GetOptions returns the Options field value
// If the value is explicit nil, the zero value for []TypesDHCPOption will be returned
func (o *DhcpAPIScope) GetOptions() []TypesDHCPOption {
	if o == nil {
		var ret []TypesDHCPOption
		return ret
	}

	return o.Options
}

// GetOptionsOk returns a tuple with the Options field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *DhcpAPIScope) GetOptionsOk() ([]TypesDHCPOption, bool) {
	if o == nil || IsNil(o.Options) {
		return nil, false
	}
	return o.Options, true
}

// SetOptions sets field value
func (o *DhcpAPIScope) SetOptions(v []TypesDHCPOption) {
	o.Options = v
}

// GetScope returns the Scope field value
func (o *DhcpAPIScope) GetScope() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Scope
}

// GetScopeOk returns a tuple with the Scope field value
// and a boolean to check if the value has been set.
func (o *DhcpAPIScope) GetScopeOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Scope, true
}

// SetScope sets field value
func (o *DhcpAPIScope) SetScope(v string) {
	o.Scope = v
}

// GetStatistics returns the Statistics field value
func (o *DhcpAPIScope) GetStatistics() DhcpAPIScopeStatistics {
	if o == nil {
		var ret DhcpAPIScopeStatistics
		return ret
	}

	return o.Statistics
}

// GetStatisticsOk returns a tuple with the Statistics field value
// and a boolean to check if the value has been set.
func (o *DhcpAPIScope) GetStatisticsOk() (*DhcpAPIScopeStatistics, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Statistics, true
}

// SetStatistics sets field value
func (o *DhcpAPIScope) SetStatistics(v DhcpAPIScopeStatistics) {
	o.Statistics = v
}

// GetSubnetCidr returns the SubnetCidr field value
func (o *DhcpAPIScope) GetSubnetCidr() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.SubnetCidr
}

// GetSubnetCidrOk returns a tuple with the SubnetCidr field value
// and a boolean to check if the value has been set.
func (o *DhcpAPIScope) GetSubnetCidrOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.SubnetCidr, true
}

// SetSubnetCidr sets field value
func (o *DhcpAPIScope) SetSubnetCidr(v string) {
	o.SubnetCidr = v
}

// GetTtl returns the Ttl field value
func (o *DhcpAPIScope) GetTtl() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Ttl
}

// GetTtlOk returns a tuple with the Ttl field value
// and a boolean to check if the value has been set.
func (o *DhcpAPIScope) GetTtlOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Ttl, true
}

// SetTtl sets field value
func (o *DhcpAPIScope) SetTtl(v int32) {
	o.Ttl = v
}

func (o DhcpAPIScope) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o DhcpAPIScope) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["default"] = o.Default
	if !IsNil(o.Dns) {
		toSerialize["dns"] = o.Dns
	}
	toSerialize["hook"] = o.Hook
	if o.Ipam != nil {
		toSerialize["ipam"] = o.Ipam
	}
	if o.Options != nil {
		toSerialize["options"] = o.Options
	}
	toSerialize["scope"] = o.Scope
	toSerialize["statistics"] = o.Statistics
	toSerialize["subnetCidr"] = o.SubnetCidr
	toSerialize["ttl"] = o.Ttl
	return toSerialize, nil
}

type NullableDhcpAPIScope struct {
	value *DhcpAPIScope
	isSet bool
}

func (v NullableDhcpAPIScope) Get() *DhcpAPIScope {
	return v.value
}

func (v *NullableDhcpAPIScope) Set(val *DhcpAPIScope) {
	v.value = val
	v.isSet = true
}

func (v NullableDhcpAPIScope) IsSet() bool {
	return v.isSet
}

func (v *NullableDhcpAPIScope) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDhcpAPIScope(val *DhcpAPIScope) *NullableDhcpAPIScope {
	return &NullableDhcpAPIScope{value: val, isSet: true}
}

func (v NullableDhcpAPIScope) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDhcpAPIScope) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
