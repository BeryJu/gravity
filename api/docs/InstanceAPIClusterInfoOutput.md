# InstanceAPIClusterInfoOutput

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ClusterVersion** | **string** |  | 
**ClusterVersionShort** | **string** |  | 
**Instances** | [**[]InstanceInstanceInfo**](InstanceInstanceInfo.md) |  | 

## Methods

### NewInstanceAPIClusterInfoOutput

`func NewInstanceAPIClusterInfoOutput(clusterVersion string, clusterVersionShort string, instances []InstanceInstanceInfo, ) *InstanceAPIClusterInfoOutput`

NewInstanceAPIClusterInfoOutput instantiates a new InstanceAPIClusterInfoOutput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewInstanceAPIClusterInfoOutputWithDefaults

`func NewInstanceAPIClusterInfoOutputWithDefaults() *InstanceAPIClusterInfoOutput`

NewInstanceAPIClusterInfoOutputWithDefaults instantiates a new InstanceAPIClusterInfoOutput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetClusterVersion

`func (o *InstanceAPIClusterInfoOutput) GetClusterVersion() string`

GetClusterVersion returns the ClusterVersion field if non-nil, zero value otherwise.

### GetClusterVersionOk

`func (o *InstanceAPIClusterInfoOutput) GetClusterVersionOk() (*string, bool)`

GetClusterVersionOk returns a tuple with the ClusterVersion field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClusterVersion

`func (o *InstanceAPIClusterInfoOutput) SetClusterVersion(v string)`

SetClusterVersion sets ClusterVersion field to given value.


### GetClusterVersionShort

`func (o *InstanceAPIClusterInfoOutput) GetClusterVersionShort() string`

GetClusterVersionShort returns the ClusterVersionShort field if non-nil, zero value otherwise.

### GetClusterVersionShortOk

`func (o *InstanceAPIClusterInfoOutput) GetClusterVersionShortOk() (*string, bool)`

GetClusterVersionShortOk returns a tuple with the ClusterVersionShort field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetClusterVersionShort

`func (o *InstanceAPIClusterInfoOutput) SetClusterVersionShort(v string)`

SetClusterVersionShort sets ClusterVersionShort field to given value.


### GetInstances

`func (o *InstanceAPIClusterInfoOutput) GetInstances() []InstanceInstanceInfo`

GetInstances returns the Instances field if non-nil, zero value otherwise.

### GetInstancesOk

`func (o *InstanceAPIClusterInfoOutput) GetInstancesOk() (*[]InstanceInstanceInfo, bool)`

GetInstancesOk returns a tuple with the Instances field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetInstances

`func (o *InstanceAPIClusterInfoOutput) SetInstances(v []InstanceInstanceInfo)`

SetInstances sets Instances field to given value.


### SetInstancesNil

`func (o *InstanceAPIClusterInfoOutput) SetInstancesNil(b bool)`

 SetInstancesNil sets the value for Instances to be an explicit nil

### UnsetInstances
`func (o *InstanceAPIClusterInfoOutput) UnsetInstances()`

UnsetInstances ensures that no value is present for Instances, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


