# TsdbRoleConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Enabled** | Pointer to **bool** |  | [optional] 
**Expire** | Pointer to **int32** |  | [optional] 
**Scrape** | Pointer to **int32** |  | [optional] 

## Methods

### NewTsdbRoleConfig

`func NewTsdbRoleConfig() *TsdbRoleConfig`

NewTsdbRoleConfig instantiates a new TsdbRoleConfig object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTsdbRoleConfigWithDefaults

`func NewTsdbRoleConfigWithDefaults() *TsdbRoleConfig`

NewTsdbRoleConfigWithDefaults instantiates a new TsdbRoleConfig object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetEnabled

`func (o *TsdbRoleConfig) GetEnabled() bool`

GetEnabled returns the Enabled field if non-nil, zero value otherwise.

### GetEnabledOk

`func (o *TsdbRoleConfig) GetEnabledOk() (*bool, bool)`

GetEnabledOk returns a tuple with the Enabled field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetEnabled

`func (o *TsdbRoleConfig) SetEnabled(v bool)`

SetEnabled sets Enabled field to given value.

### HasEnabled

`func (o *TsdbRoleConfig) HasEnabled() bool`

HasEnabled returns a boolean if a field has been set.

### GetExpire

`func (o *TsdbRoleConfig) GetExpire() int32`

GetExpire returns the Expire field if non-nil, zero value otherwise.

### GetExpireOk

`func (o *TsdbRoleConfig) GetExpireOk() (*int32, bool)`

GetExpireOk returns a tuple with the Expire field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetExpire

`func (o *TsdbRoleConfig) SetExpire(v int32)`

SetExpire sets Expire field to given value.

### HasExpire

`func (o *TsdbRoleConfig) HasExpire() bool`

HasExpire returns a boolean if a field has been set.

### GetScrape

`func (o *TsdbRoleConfig) GetScrape() int32`

GetScrape returns the Scrape field if non-nil, zero value otherwise.

### GetScrapeOk

`func (o *TsdbRoleConfig) GetScrapeOk() (*int32, bool)`

GetScrapeOk returns a tuple with the Scrape field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScrape

`func (o *TsdbRoleConfig) SetScrape(v int32)`

SetScrape sets Scrape field to given value.

### HasScrape

`func (o *TsdbRoleConfig) HasScrape() bool`

HasScrape returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


