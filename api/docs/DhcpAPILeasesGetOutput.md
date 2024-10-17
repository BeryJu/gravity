# DhcpAPILeasesGetOutput

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Leases** | [**[]DhcpAPILease**](DhcpAPILease.md) |  | 

## Methods

### NewDhcpAPILeasesGetOutput

`func NewDhcpAPILeasesGetOutput(leases []DhcpAPILease, ) *DhcpAPILeasesGetOutput`

NewDhcpAPILeasesGetOutput instantiates a new DhcpAPILeasesGetOutput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDhcpAPILeasesGetOutputWithDefaults

`func NewDhcpAPILeasesGetOutputWithDefaults() *DhcpAPILeasesGetOutput`

NewDhcpAPILeasesGetOutputWithDefaults instantiates a new DhcpAPILeasesGetOutput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetLeases

`func (o *DhcpAPILeasesGetOutput) GetLeases() []DhcpAPILease`

GetLeases returns the Leases field if non-nil, zero value otherwise.

### GetLeasesOk

`func (o *DhcpAPILeasesGetOutput) GetLeasesOk() (*[]DhcpAPILease, bool)`

GetLeasesOk returns a tuple with the Leases field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetLeases

`func (o *DhcpAPILeasesGetOutput) SetLeases(v []DhcpAPILease)`

SetLeases sets Leases field to given value.


### SetLeasesNil

`func (o *DhcpAPILeasesGetOutput) SetLeasesNil(b bool)`

 SetLeasesNil sets the value for Leases to be an explicit nil

### UnsetLeases
`func (o *DhcpAPILeasesGetOutput) UnsetLeases()`

UnsetLeases ensures that no value is present for Leases, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


