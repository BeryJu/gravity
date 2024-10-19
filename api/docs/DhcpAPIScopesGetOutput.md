# DhcpAPIScopesGetOutput

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Scopes** | [**[]DhcpAPIScope**](DhcpAPIScope.md) |  | 
**Statistics** | [**DhcpAPIScopeStatistics**](DhcpAPIScopeStatistics.md) |  | 

## Methods

### NewDhcpAPIScopesGetOutput

`func NewDhcpAPIScopesGetOutput(scopes []DhcpAPIScope, statistics DhcpAPIScopeStatistics, ) *DhcpAPIScopesGetOutput`

NewDhcpAPIScopesGetOutput instantiates a new DhcpAPIScopesGetOutput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewDhcpAPIScopesGetOutputWithDefaults

`func NewDhcpAPIScopesGetOutputWithDefaults() *DhcpAPIScopesGetOutput`

NewDhcpAPIScopesGetOutputWithDefaults instantiates a new DhcpAPIScopesGetOutput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetScopes

`func (o *DhcpAPIScopesGetOutput) GetScopes() []DhcpAPIScope`

GetScopes returns the Scopes field if non-nil, zero value otherwise.

### GetScopesOk

`func (o *DhcpAPIScopesGetOutput) GetScopesOk() (*[]DhcpAPIScope, bool)`

GetScopesOk returns a tuple with the Scopes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScopes

`func (o *DhcpAPIScopesGetOutput) SetScopes(v []DhcpAPIScope)`

SetScopes sets Scopes field to given value.


### SetScopesNil

`func (o *DhcpAPIScopesGetOutput) SetScopesNil(b bool)`

 SetScopesNil sets the value for Scopes to be an explicit nil

### UnsetScopes
`func (o *DhcpAPIScopesGetOutput) UnsetScopes()`

UnsetScopes ensures that no value is present for Scopes, not even an explicit nil
### GetStatistics

`func (o *DhcpAPIScopesGetOutput) GetStatistics() DhcpAPIScopeStatistics`

GetStatistics returns the Statistics field if non-nil, zero value otherwise.

### GetStatisticsOk

`func (o *DhcpAPIScopesGetOutput) GetStatisticsOk() (*DhcpAPIScopeStatistics, bool)`

GetStatisticsOk returns a tuple with the Statistics field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatistics

`func (o *DhcpAPIScopesGetOutput) SetStatistics(v DhcpAPIScopeStatistics)`

SetStatistics sets Statistics field to given value.



[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


