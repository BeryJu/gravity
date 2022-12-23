# ApiRoleConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**CookieSecret** | Pointer to **string** |  | [optional] 
**Oidc** | Pointer to [**TypesOIDCConfig**](TypesOIDCConfig.md) |  | [optional] 
**Port** | Pointer to **int32** |  | [optional] 

## Methods

### NewApiRoleConfig

`func NewApiRoleConfig() *ApiRoleConfig`

NewApiRoleConfig instantiates a new ApiRoleConfig object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewApiRoleConfigWithDefaults

`func NewApiRoleConfigWithDefaults() *ApiRoleConfig`

NewApiRoleConfigWithDefaults instantiates a new ApiRoleConfig object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCookieSecret

`func (o *ApiRoleConfig) GetCookieSecret() string`

GetCookieSecret returns the CookieSecret field if non-nil, zero value otherwise.

### GetCookieSecretOk

`func (o *ApiRoleConfig) GetCookieSecretOk() (*string, bool)`

GetCookieSecretOk returns a tuple with the CookieSecret field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCookieSecret

`func (o *ApiRoleConfig) SetCookieSecret(v string)`

SetCookieSecret sets CookieSecret field to given value.

### HasCookieSecret

`func (o *ApiRoleConfig) HasCookieSecret() bool`

HasCookieSecret returns a boolean if a field has been set.

### GetOidc

`func (o *ApiRoleConfig) GetOidc() TypesOIDCConfig`

GetOidc returns the Oidc field if non-nil, zero value otherwise.

### GetOidcOk

`func (o *ApiRoleConfig) GetOidcOk() (*TypesOIDCConfig, bool)`

GetOidcOk returns a tuple with the Oidc field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetOidc

`func (o *ApiRoleConfig) SetOidc(v TypesOIDCConfig)`

SetOidc sets Oidc field to given value.

### HasOidc

`func (o *ApiRoleConfig) HasOidc() bool`

HasOidc returns a boolean if a field has been set.

### GetPort

`func (o *ApiRoleConfig) GetPort() int32`

GetPort returns the Port field if non-nil, zero value otherwise.

### GetPortOk

`func (o *ApiRoleConfig) GetPortOk() (*int32, bool)`

GetPortOk returns a tuple with the Port field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPort

`func (o *ApiRoleConfig) SetPort(v int32)`

SetPort sets Port field to given value.

### HasPort

`func (o *ApiRoleConfig) HasPort() bool`

HasPort returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


