# DhcpAPILease

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Address** | **string** |  | 
**AddressLeaseTime** | **string** |  | 
**DnsZone** | Pointer to **string** |  | [optional] 
**Hostname** | **string** |  | 
**Identifier** | **string** |  | 
**Info** | Pointer to [**DhcpAPILeaseInfo**](DhcpAPILeaseInfo.md) |  | [optional] 
**ScopeKey** | **string** |  | 

## Methods

### NewDhcpAPILease

`func NewDhcpAPILease(address string, addressLeaseTime string, hostname string, identifier string, scopeKey string, ) *DhcpAPILease`

NewDhcpAPILease instantiates a new DhcpAPILease object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDhcpAPILeaseWithDefaults

`func NewDhcpAPILeaseWithDefaults() *DhcpAPILease`

NewDhcpAPILeaseWithDefaults instantiates a new DhcpAPILease object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAddress

`func (o *DhcpAPILease) GetAddress() string`

GetAddress returns the Address field if non-nil, zero value otherwise.

### GetAddressOk

`func (o *DhcpAPILease) GetAddressOk() (*string, bool)`

GetAddressOk returns a tuple with the Address field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddress

`func (o *DhcpAPILease) SetAddress(v string)`

SetAddress sets Address field to given value.


### GetAddressLeaseTime

`func (o *DhcpAPILease) GetAddressLeaseTime() string`

GetAddressLeaseTime returns the AddressLeaseTime field if non-nil, zero value otherwise.

### GetAddressLeaseTimeOk

`func (o *DhcpAPILease) GetAddressLeaseTimeOk() (*string, bool)`

GetAddressLeaseTimeOk returns a tuple with the AddressLeaseTime field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAddressLeaseTime

`func (o *DhcpAPILease) SetAddressLeaseTime(v string)`

SetAddressLeaseTime sets AddressLeaseTime field to given value.


### GetDnsZone

`func (o *DhcpAPILease) GetDnsZone() string`

GetDnsZone returns the DnsZone field if non-nil, zero value otherwise.

### GetDnsZoneOk

`func (o *DhcpAPILease) GetDnsZoneOk() (*string, bool)`

GetDnsZoneOk returns a tuple with the DnsZone field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDnsZone

`func (o *DhcpAPILease) SetDnsZone(v string)`

SetDnsZone sets DnsZone field to given value.

### HasDnsZone

`func (o *DhcpAPILease) HasDnsZone() bool`

HasDnsZone returns a boolean if a field has been set.

### GetHostname

`func (o *DhcpAPILease) GetHostname() string`

GetHostname returns the Hostname field if non-nil, zero value otherwise.

### GetHostnameOk

`func (o *DhcpAPILease) GetHostnameOk() (*string, bool)`

GetHostnameOk returns a tuple with the Hostname field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHostname

`func (o *DhcpAPILease) SetHostname(v string)`

SetHostname sets Hostname field to given value.


### GetIdentifier

`func (o *DhcpAPILease) GetIdentifier() string`

GetIdentifier returns the Identifier field if non-nil, zero value otherwise.

### GetIdentifierOk

`func (o *DhcpAPILease) GetIdentifierOk() (*string, bool)`

GetIdentifierOk returns a tuple with the Identifier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIdentifier

`func (o *DhcpAPILease) SetIdentifier(v string)`

SetIdentifier sets Identifier field to given value.


### GetInfo

`func (o *DhcpAPILease) GetInfo() DhcpAPILeaseInfo`

GetInfo returns the Info field if non-nil, zero value otherwise.

### GetInfoOk

`func (o *DhcpAPILease) GetInfoOk() (*DhcpAPILeaseInfo, bool)`

GetInfoOk returns a tuple with the Info field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInfo

`func (o *DhcpAPILease) SetInfo(v DhcpAPILeaseInfo)`

SetInfo sets Info field to given value.

### HasInfo

`func (o *DhcpAPILease) HasInfo() bool`

HasInfo returns a boolean if a field has been set.

### GetScopeKey

`func (o *DhcpAPILease) GetScopeKey() string`

GetScopeKey returns the ScopeKey field if non-nil, zero value otherwise.

### GetScopeKeyOk

`func (o *DhcpAPILease) GetScopeKeyOk() (*string, bool)`

GetScopeKeyOk returns a tuple with the ScopeKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScopeKey

`func (o *DhcpAPILease) SetScopeKey(v string)`

SetScopeKey sets ScopeKey field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


