# EtcdAPIMember

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Id** | Pointer to **string** |  | [optional] 
**IsLeader** | Pointer to **bool** |  | [optional] 
**IsLearner** | Pointer to **bool** |  | [optional] 
**Name** | Pointer to **string** |  | [optional] 

## Methods

### NewEtcdAPIMember

`func NewEtcdAPIMember() *EtcdAPIMember`

NewEtcdAPIMember instantiates a new EtcdAPIMember object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEtcdAPIMemberWithDefaults

`func NewEtcdAPIMemberWithDefaults() *EtcdAPIMember`

NewEtcdAPIMemberWithDefaults instantiates a new EtcdAPIMember object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetId

`func (o *EtcdAPIMember) GetId() string`

GetId returns the Id field if non-nil, zero value otherwise.

### GetIdOk

`func (o *EtcdAPIMember) GetIdOk() (*string, bool)`

GetIdOk returns a tuple with the Id field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetId

`func (o *EtcdAPIMember) SetId(v string)`

SetId sets Id field to given value.

### HasId

`func (o *EtcdAPIMember) HasId() bool`

HasId returns a boolean if a field has been set.

### GetIsLeader

`func (o *EtcdAPIMember) GetIsLeader() bool`

GetIsLeader returns the IsLeader field if non-nil, zero value otherwise.

### GetIsLeaderOk

`func (o *EtcdAPIMember) GetIsLeaderOk() (*bool, bool)`

GetIsLeaderOk returns a tuple with the IsLeader field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsLeader

`func (o *EtcdAPIMember) SetIsLeader(v bool)`

SetIsLeader sets IsLeader field to given value.

### HasIsLeader

`func (o *EtcdAPIMember) HasIsLeader() bool`

HasIsLeader returns a boolean if a field has been set.

### GetIsLearner

`func (o *EtcdAPIMember) GetIsLearner() bool`

GetIsLearner returns the IsLearner field if non-nil, zero value otherwise.

### GetIsLearnerOk

`func (o *EtcdAPIMember) GetIsLearnerOk() (*bool, bool)`

GetIsLearnerOk returns a tuple with the IsLearner field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetIsLearner

`func (o *EtcdAPIMember) SetIsLearner(v bool)`

SetIsLearner sets IsLearner field to given value.

### HasIsLearner

`func (o *EtcdAPIMember) HasIsLearner() bool`

HasIsLearner returns a boolean if a field has been set.

### GetName

`func (o *EtcdAPIMember) GetName() string`

GetName returns the Name field if non-nil, zero value otherwise.

### GetNameOk

`func (o *EtcdAPIMember) GetNameOk() (*string, bool)`

GetNameOk returns a tuple with the Name field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetName

`func (o *EtcdAPIMember) SetName(v string)`

SetName sets Name field to given value.

### HasName

`func (o *EtcdAPIMember) HasName() bool`

HasName returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


