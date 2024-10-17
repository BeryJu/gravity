# TypesOIDCConfig

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ClientID** | Pointer to **string** |  | [optional] 
**ClientSecret** | Pointer to **string** |  | [optional] 
**Issuer** | Pointer to **string** |  | [optional] 
**RedirectURL** | Pointer to **string** |  | [optional] 
**Scopes** | Pointer to **[]string** |  | [optional] 

## Methods

### NewTypesOIDCConfig

`func NewTypesOIDCConfig() *TypesOIDCConfig`

NewTypesOIDCConfig instantiates a new TypesOIDCConfig object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewTypesOIDCConfigWithDefaults

`func NewTypesOIDCConfigWithDefaults() *TypesOIDCConfig`

NewTypesOIDCConfigWithDefaults instantiates a new TypesOIDCConfig object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetClientID

`func (o *TypesOIDCConfig) GetClientID() string`

GetClientID returns the ClientID field if non-nil, zero value otherwise.

### GetClientIDOk

`func (o *TypesOIDCConfig) GetClientIDOk() (*string, bool)`

GetClientIDOk returns a tuple with the ClientID field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientID

`func (o *TypesOIDCConfig) SetClientID(v string)`

SetClientID sets ClientID field to given value.

### HasClientID

`func (o *TypesOIDCConfig) HasClientID() bool`

HasClientID returns a boolean if a field has been set.

### GetClientSecret

`func (o *TypesOIDCConfig) GetClientSecret() string`

GetClientSecret returns the ClientSecret field if non-nil, zero value otherwise.

### GetClientSecretOk

`func (o *TypesOIDCConfig) GetClientSecretOk() (*string, bool)`

GetClientSecretOk returns a tuple with the ClientSecret field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClientSecret

`func (o *TypesOIDCConfig) SetClientSecret(v string)`

SetClientSecret sets ClientSecret field to given value.

### HasClientSecret

`func (o *TypesOIDCConfig) HasClientSecret() bool`

HasClientSecret returns a boolean if a field has been set.

### GetIssuer

`func (o *TypesOIDCConfig) GetIssuer() string`

GetIssuer returns the Issuer field if non-nil, zero value otherwise.

### GetIssuerOk

`func (o *TypesOIDCConfig) GetIssuerOk() (*string, bool)`

GetIssuerOk returns a tuple with the Issuer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIssuer

`func (o *TypesOIDCConfig) SetIssuer(v string)`

SetIssuer sets Issuer field to given value.

### HasIssuer

`func (o *TypesOIDCConfig) HasIssuer() bool`

HasIssuer returns a boolean if a field has been set.

### GetRedirectURL

`func (o *TypesOIDCConfig) GetRedirectURL() string`

GetRedirectURL returns the RedirectURL field if non-nil, zero value otherwise.

### GetRedirectURLOk

`func (o *TypesOIDCConfig) GetRedirectURLOk() (*string, bool)`

GetRedirectURLOk returns a tuple with the RedirectURL field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRedirectURL

`func (o *TypesOIDCConfig) SetRedirectURL(v string)`

SetRedirectURL sets RedirectURL field to given value.

### HasRedirectURL

`func (o *TypesOIDCConfig) HasRedirectURL() bool`

HasRedirectURL returns a boolean if a field has been set.

### GetScopes

`func (o *TypesOIDCConfig) GetScopes() []string`

GetScopes returns the Scopes field if non-nil, zero value otherwise.

### GetScopesOk

`func (o *TypesOIDCConfig) GetScopesOk() (*[]string, bool)`

GetScopesOk returns a tuple with the Scopes field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetScopes

`func (o *TypesOIDCConfig) SetScopes(v []string)`

SetScopes sets Scopes field to given value.

### HasScopes

`func (o *TypesOIDCConfig) HasScopes() bool`

HasScopes returns a boolean if a field has been set.

### SetScopesNil

`func (o *TypesOIDCConfig) SetScopesNil(b bool)`

 SetScopesNil sets the value for Scopes to be an explicit nil

### UnsetScopes
`func (o *TypesOIDCConfig) UnsetScopes()`

UnsetScopes ensures that no value is present for Scopes, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


