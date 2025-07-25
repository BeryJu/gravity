/*
gravity

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 0.27.2
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package api

import (
	"encoding/json"
)

// checks if the DhcpAPIScopesPutInput type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &DhcpAPIScopesPutInput{}

// DhcpAPIScopesPutInput struct for DhcpAPIScopesPutInput
type DhcpAPIScopesPutInput struct {
	Default    bool              `json:"default"`
	Dns        *DhcpScopeDNS     `json:"dns,omitempty"`
	Hook       string            `json:"hook"`
	Ipam       map[string]string `json:"ipam,omitempty"`
	Options    []TypesDHCPOption `json:"options"`
	SubnetCidr string            `json:"subnetCidr"`
	Ttl        int32             `json:"ttl"`
}

// NewDhcpAPIScopesPutInput instantiates a new DhcpAPIScopesPutInput object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDhcpAPIScopesPutInput(default_ bool, hook string, options []TypesDHCPOption, subnetCidr string, ttl int32) *DhcpAPIScopesPutInput {
	this := DhcpAPIScopesPutInput{}
	this.Default = default_
	this.Hook = hook
	this.Options = options
	this.SubnetCidr = subnetCidr
	this.Ttl = ttl
	return &this
}

// NewDhcpAPIScopesPutInputWithDefaults instantiates a new DhcpAPIScopesPutInput object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDhcpAPIScopesPutInputWithDefaults() *DhcpAPIScopesPutInput {
	this := DhcpAPIScopesPutInput{}
	return &this
}

// GetDefault returns the Default field value
func (o *DhcpAPIScopesPutInput) GetDefault() bool {
	if o == nil {
		var ret bool
		return ret
	}

	return o.Default
}

// GetDefaultOk returns a tuple with the Default field value
// and a boolean to check if the value has been set.
func (o *DhcpAPIScopesPutInput) GetDefaultOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Default, true
}

// SetDefault sets field value
func (o *DhcpAPIScopesPutInput) SetDefault(v bool) {
	o.Default = v
}

// GetDns returns the Dns field value if set, zero value otherwise.
func (o *DhcpAPIScopesPutInput) GetDns() DhcpScopeDNS {
	if o == nil || IsNil(o.Dns) {
		var ret DhcpScopeDNS
		return ret
	}
	return *o.Dns
}

// GetDnsOk returns a tuple with the Dns field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DhcpAPIScopesPutInput) GetDnsOk() (*DhcpScopeDNS, bool) {
	if o == nil || IsNil(o.Dns) {
		return nil, false
	}
	return o.Dns, true
}

// HasDns returns a boolean if a field has been set.
func (o *DhcpAPIScopesPutInput) HasDns() bool {
	if o != nil && !IsNil(o.Dns) {
		return true
	}

	return false
}

// SetDns gets a reference to the given DhcpScopeDNS and assigns it to the Dns field.
func (o *DhcpAPIScopesPutInput) SetDns(v DhcpScopeDNS) {
	o.Dns = &v
}

// GetHook returns the Hook field value
func (o *DhcpAPIScopesPutInput) GetHook() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.Hook
}

// GetHookOk returns a tuple with the Hook field value
// and a boolean to check if the value has been set.
func (o *DhcpAPIScopesPutInput) GetHookOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Hook, true
}

// SetHook sets field value
func (o *DhcpAPIScopesPutInput) SetHook(v string) {
	o.Hook = v
}

// GetIpam returns the Ipam field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *DhcpAPIScopesPutInput) GetIpam() map[string]string {
	if o == nil {
		var ret map[string]string
		return ret
	}
	return o.Ipam
}

// GetIpamOk returns a tuple with the Ipam field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *DhcpAPIScopesPutInput) GetIpamOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.Ipam) {
		return nil, false
	}
	return &o.Ipam, true
}

// HasIpam returns a boolean if a field has been set.
func (o *DhcpAPIScopesPutInput) HasIpam() bool {
	if o != nil && IsNil(o.Ipam) {
		return true
	}

	return false
}

// SetIpam gets a reference to the given map[string]string and assigns it to the Ipam field.
func (o *DhcpAPIScopesPutInput) SetIpam(v map[string]string) {
	o.Ipam = v
}

// GetOptions returns the Options field value
// If the value is explicit nil, the zero value for []TypesDHCPOption will be returned
func (o *DhcpAPIScopesPutInput) GetOptions() []TypesDHCPOption {
	if o == nil {
		var ret []TypesDHCPOption
		return ret
	}

	return o.Options
}

// GetOptionsOk returns a tuple with the Options field value
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *DhcpAPIScopesPutInput) GetOptionsOk() ([]TypesDHCPOption, bool) {
	if o == nil || IsNil(o.Options) {
		return nil, false
	}
	return o.Options, true
}

// SetOptions sets field value
func (o *DhcpAPIScopesPutInput) SetOptions(v []TypesDHCPOption) {
	o.Options = v
}

// GetSubnetCidr returns the SubnetCidr field value
func (o *DhcpAPIScopesPutInput) GetSubnetCidr() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.SubnetCidr
}

// GetSubnetCidrOk returns a tuple with the SubnetCidr field value
// and a boolean to check if the value has been set.
func (o *DhcpAPIScopesPutInput) GetSubnetCidrOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.SubnetCidr, true
}

// SetSubnetCidr sets field value
func (o *DhcpAPIScopesPutInput) SetSubnetCidr(v string) {
	o.SubnetCidr = v
}

// GetTtl returns the Ttl field value
func (o *DhcpAPIScopesPutInput) GetTtl() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Ttl
}

// GetTtlOk returns a tuple with the Ttl field value
// and a boolean to check if the value has been set.
func (o *DhcpAPIScopesPutInput) GetTtlOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Ttl, true
}

// SetTtl sets field value
func (o *DhcpAPIScopesPutInput) SetTtl(v int32) {
	o.Ttl = v
}

func (o DhcpAPIScopesPutInput) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o DhcpAPIScopesPutInput) ToMap() (map[string]interface{}, error) {
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
	toSerialize["subnetCidr"] = o.SubnetCidr
	toSerialize["ttl"] = o.Ttl
	return toSerialize, nil
}

type NullableDhcpAPIScopesPutInput struct {
	value *DhcpAPIScopesPutInput
	isSet bool
}

func (v NullableDhcpAPIScopesPutInput) Get() *DhcpAPIScopesPutInput {
	return v.value
}

func (v *NullableDhcpAPIScopesPutInput) Set(val *DhcpAPIScopesPutInput) {
	v.value = val
	v.isSet = true
}

func (v NullableDhcpAPIScopesPutInput) IsSet() bool {
	return v.isSet
}

func (v *NullableDhcpAPIScopesPutInput) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDhcpAPIScopesPutInput(val *DhcpAPIScopesPutInput) *NullableDhcpAPIScopesPutInput {
	return &NullableDhcpAPIScopesPutInput{value: val, isSet: true}
}

func (v NullableDhcpAPIScopesPutInput) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDhcpAPIScopesPutInput) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
