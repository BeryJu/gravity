# EtcdAPIMemberJoinInput

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Identifier** | Pointer to **string** |  | [optional] 
**Peer** | Pointer to **string** |  | [optional] 
**Roles** | Pointer to **string** |  | [optional] 

## Methods

### NewEtcdAPIMemberJoinInput

`func NewEtcdAPIMemberJoinInput() *EtcdAPIMemberJoinInput`

NewEtcdAPIMemberJoinInput instantiates a new EtcdAPIMemberJoinInput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEtcdAPIMemberJoinInputWithDefaults

`func NewEtcdAPIMemberJoinInputWithDefaults() *EtcdAPIMemberJoinInput`

NewEtcdAPIMemberJoinInputWithDefaults instantiates a new EtcdAPIMemberJoinInput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetIdentifier

`func (o *EtcdAPIMemberJoinInput) GetIdentifier() string`

GetIdentifier returns the Identifier field if non-nil, zero value otherwise.

### GetIdentifierOk

`func (o *EtcdAPIMemberJoinInput) GetIdentifierOk() (*string, bool)`

GetIdentifierOk returns a tuple with the Identifier field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIdentifier

`func (o *EtcdAPIMemberJoinInput) SetIdentifier(v string)`

SetIdentifier sets Identifier field to given value.

### HasIdentifier

`func (o *EtcdAPIMemberJoinInput) HasIdentifier() bool`

HasIdentifier returns a boolean if a field has been set.

### GetPeer

`func (o *EtcdAPIMemberJoinInput) GetPeer() string`

GetPeer returns the Peer field if non-nil, zero value otherwise.

### GetPeerOk

`func (o *EtcdAPIMemberJoinInput) GetPeerOk() (*string, bool)`

GetPeerOk returns a tuple with the Peer field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPeer

`func (o *EtcdAPIMemberJoinInput) SetPeer(v string)`

SetPeer sets Peer field to given value.

### HasPeer

`func (o *EtcdAPIMemberJoinInput) HasPeer() bool`

HasPeer returns a boolean if a field has been set.

### GetRoles

`func (o *EtcdAPIMemberJoinInput) GetRoles() string`

GetRoles returns the Roles field if non-nil, zero value otherwise.

### GetRolesOk

`func (o *EtcdAPIMemberJoinInput) GetRolesOk() (*string, bool)`

GetRolesOk returns a tuple with the Roles field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRoles

`func (o *EtcdAPIMemberJoinInput) SetRoles(v string)`

SetRoles sets Roles field to given value.

### HasRoles

`func (o *EtcdAPIMemberJoinInput) HasRoles() bool`

HasRoles returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


