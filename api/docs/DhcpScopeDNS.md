# DhcpScopeDNS

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AddZoneInHostname** | Pointer to **bool** |  | [optional] 
**Search** | Pointer to **[]string** |  | [optional] 
**Zone** | Pointer to **string** |  | [optional] 

## Methods

### NewDhcpScopeDNS

`func NewDhcpScopeDNS() *DhcpScopeDNS`

NewDhcpScopeDNS instantiates a new DhcpScopeDNS object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDhcpScopeDNSWithDefaults

`func NewDhcpScopeDNSWithDefaults() *DhcpScopeDNS`

NewDhcpScopeDNSWithDefaults instantiates a new DhcpScopeDNS object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAddZoneInHostname

`func (o *DhcpScopeDNS) GetAddZoneInHostname() bool`

GetAddZoneInHostname returns the AddZoneInHostname field if non-nil, zero value otherwise.

### GetAddZoneInHostnameOk

`func (o *DhcpScopeDNS) GetAddZoneInHostnameOk() (*bool, bool)`

GetAddZoneInHostnameOk returns a tuple with the AddZoneInHostname field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddZoneInHostname

`func (o *DhcpScopeDNS) SetAddZoneInHostname(v bool)`

SetAddZoneInHostname sets AddZoneInHostname field to given value.

### HasAddZoneInHostname

`func (o *DhcpScopeDNS) HasAddZoneInHostname() bool`

HasAddZoneInHostname returns a boolean if a field has been set.

### GetSearch

`func (o *DhcpScopeDNS) GetSearch() []string`

GetSearch returns the Search field if non-nil, zero value otherwise.

### GetSearchOk

`func (o *DhcpScopeDNS) GetSearchOk() (*[]string, bool)`

GetSearchOk returns a tuple with the Search field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSearch

`func (o *DhcpScopeDNS) SetSearch(v []string)`

SetSearch sets Search field to given value.

### HasSearch

`func (o *DhcpScopeDNS) HasSearch() bool`

HasSearch returns a boolean if a field has been set.

### SetSearchNil

`func (o *DhcpScopeDNS) SetSearchNil(b bool)`

 SetSearchNil sets the value for Search to be an explicit nil

### UnsetSearch
`func (o *DhcpScopeDNS) UnsetSearch()`

UnsetSearch ensures that no value is present for Search, not even an explicit nil
### GetZone

`func (o *DhcpScopeDNS) GetZone() string`

GetZone returns the Zone field if non-nil, zero value otherwise.

### GetZoneOk

`func (o *DhcpScopeDNS) GetZoneOk() (*string, bool)`

GetZoneOk returns a tuple with the Zone field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetZone

`func (o *DhcpScopeDNS) SetZone(v string)`

SetZone sets Zone field to given value.

### HasZone

`func (o *DhcpScopeDNS) HasZone() bool`

HasZone returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


