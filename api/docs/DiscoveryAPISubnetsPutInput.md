# DiscoveryAPISubnetsPutInput

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**DiscoveryTTL** | **int32** |  | 
**DnsResolver** | **string** |  | 
**SubnetCidr** | **string** |  | 

## Methods

### NewDiscoveryAPISubnetsPutInput

`func NewDiscoveryAPISubnetsPutInput(discoveryTTL int32, dnsResolver string, subnetCidr string, ) *DiscoveryAPISubnetsPutInput`

NewDiscoveryAPISubnetsPutInput instantiates a new DiscoveryAPISubnetsPutInput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDiscoveryAPISubnetsPutInputWithDefaults

`func NewDiscoveryAPISubnetsPutInputWithDefaults() *DiscoveryAPISubnetsPutInput`

NewDiscoveryAPISubnetsPutInputWithDefaults instantiates a new DiscoveryAPISubnetsPutInput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDiscoveryTTL

`func (o *DiscoveryAPISubnetsPutInput) GetDiscoveryTTL() int32`

GetDiscoveryTTL returns the DiscoveryTTL field if non-nil, zero value otherwise.

### GetDiscoveryTTLOk

`func (o *DiscoveryAPISubnetsPutInput) GetDiscoveryTTLOk() (*int32, bool)`

GetDiscoveryTTLOk returns a tuple with the DiscoveryTTL field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDiscoveryTTL

`func (o *DiscoveryAPISubnetsPutInput) SetDiscoveryTTL(v int32)`

SetDiscoveryTTL sets DiscoveryTTL field to given value.


### GetDnsResolver

`func (o *DiscoveryAPISubnetsPutInput) GetDnsResolver() string`

GetDnsResolver returns the DnsResolver field if non-nil, zero value otherwise.

### GetDnsResolverOk

`func (o *DiscoveryAPISubnetsPutInput) GetDnsResolverOk() (*string, bool)`

GetDnsResolverOk returns a tuple with the DnsResolver field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDnsResolver

`func (o *DiscoveryAPISubnetsPutInput) SetDnsResolver(v string)`

SetDnsResolver sets DnsResolver field to given value.


### GetSubnetCidr

`func (o *DiscoveryAPISubnetsPutInput) GetSubnetCidr() string`

GetSubnetCidr returns the SubnetCidr field if non-nil, zero value otherwise.

### GetSubnetCidrOk

`func (o *DiscoveryAPISubnetsPutInput) GetSubnetCidrOk() (*string, bool)`

GetSubnetCidrOk returns a tuple with the SubnetCidr field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSubnetCidr

`func (o *DiscoveryAPISubnetsPutInput) SetSubnetCidr(v string)`

SetSubnetCidr sets SubnetCidr field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


