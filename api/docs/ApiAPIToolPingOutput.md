# ApiAPIToolPingOutput

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AvgRtt** | Pointer to **interface{}** |  | [optional] 
**MaxRtt** | Pointer to **interface{}** |  | [optional] 
**MinRtt** | Pointer to **int32** |  | [optional] 
**PacketLoss** | Pointer to **float32** |  | [optional] 
**PacketsRecv** | Pointer to **int32** |  | [optional] 
**PacketsRecvDuplicates** | Pointer to **int32** |  | [optional] 
**PacketsSent** | Pointer to **int32** |  | [optional] 
**StdDevRtt** | Pointer to **interface{}** |  | [optional] 

## Methods

### NewApiAPIToolPingOutput

`func NewApiAPIToolPingOutput() *ApiAPIToolPingOutput`

NewApiAPIToolPingOutput instantiates a new ApiAPIToolPingOutput object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewApiAPIToolPingOutputWithDefaults

`func NewApiAPIToolPingOutputWithDefaults() *ApiAPIToolPingOutput`

NewApiAPIToolPingOutputWithDefaults instantiates a new ApiAPIToolPingOutput object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAvgRtt

`func (o *ApiAPIToolPingOutput) GetAvgRtt() interface{}`

GetAvgRtt returns the AvgRtt field if non-nil, zero value otherwise.

### GetAvgRttOk

`func (o *ApiAPIToolPingOutput) GetAvgRttOk() (*interface{}, bool)`

GetAvgRttOk returns a tuple with the AvgRtt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvgRtt

`func (o *ApiAPIToolPingOutput) SetAvgRtt(v interface{})`

SetAvgRtt sets AvgRtt field to given value.

### HasAvgRtt

`func (o *ApiAPIToolPingOutput) HasAvgRtt() bool`

HasAvgRtt returns a boolean if a field has been set.

### SetAvgRttNil

`func (o *ApiAPIToolPingOutput) SetAvgRttNil(b bool)`

 SetAvgRttNil sets the value for AvgRtt to be an explicit nil

### UnsetAvgRtt
`func (o *ApiAPIToolPingOutput) UnsetAvgRtt()`

UnsetAvgRtt ensures that no value is present for AvgRtt, not even an explicit nil
### GetMaxRtt

`func (o *ApiAPIToolPingOutput) GetMaxRtt() interface{}`

GetMaxRtt returns the MaxRtt field if non-nil, zero value otherwise.

### GetMaxRttOk

`func (o *ApiAPIToolPingOutput) GetMaxRttOk() (*interface{}, bool)`

GetMaxRttOk returns a tuple with the MaxRtt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMaxRtt

`func (o *ApiAPIToolPingOutput) SetMaxRtt(v interface{})`

SetMaxRtt sets MaxRtt field to given value.

### HasMaxRtt

`func (o *ApiAPIToolPingOutput) HasMaxRtt() bool`

HasMaxRtt returns a boolean if a field has been set.

### SetMaxRttNil

`func (o *ApiAPIToolPingOutput) SetMaxRttNil(b bool)`

 SetMaxRttNil sets the value for MaxRtt to be an explicit nil

### UnsetMaxRtt
`func (o *ApiAPIToolPingOutput) UnsetMaxRtt()`

UnsetMaxRtt ensures that no value is present for MaxRtt, not even an explicit nil
### GetMinRtt

`func (o *ApiAPIToolPingOutput) GetMinRtt() int32`

GetMinRtt returns the MinRtt field if non-nil, zero value otherwise.

### GetMinRttOk

`func (o *ApiAPIToolPingOutput) GetMinRttOk() (*int32, bool)`

GetMinRttOk returns a tuple with the MinRtt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMinRtt

`func (o *ApiAPIToolPingOutput) SetMinRtt(v int32)`

SetMinRtt sets MinRtt field to given value.

### HasMinRtt

`func (o *ApiAPIToolPingOutput) HasMinRtt() bool`

HasMinRtt returns a boolean if a field has been set.

### GetPacketLoss

`func (o *ApiAPIToolPingOutput) GetPacketLoss() float32`

GetPacketLoss returns the PacketLoss field if non-nil, zero value otherwise.

### GetPacketLossOk

`func (o *ApiAPIToolPingOutput) GetPacketLossOk() (*float32, bool)`

GetPacketLossOk returns a tuple with the PacketLoss field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPacketLoss

`func (o *ApiAPIToolPingOutput) SetPacketLoss(v float32)`

SetPacketLoss sets PacketLoss field to given value.

### HasPacketLoss

`func (o *ApiAPIToolPingOutput) HasPacketLoss() bool`

HasPacketLoss returns a boolean if a field has been set.

### GetPacketsRecv

`func (o *ApiAPIToolPingOutput) GetPacketsRecv() int32`

GetPacketsRecv returns the PacketsRecv field if non-nil, zero value otherwise.

### GetPacketsRecvOk

`func (o *ApiAPIToolPingOutput) GetPacketsRecvOk() (*int32, bool)`

GetPacketsRecvOk returns a tuple with the PacketsRecv field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPacketsRecv

`func (o *ApiAPIToolPingOutput) SetPacketsRecv(v int32)`

SetPacketsRecv sets PacketsRecv field to given value.

### HasPacketsRecv

`func (o *ApiAPIToolPingOutput) HasPacketsRecv() bool`

HasPacketsRecv returns a boolean if a field has been set.

### GetPacketsRecvDuplicates

`func (o *ApiAPIToolPingOutput) GetPacketsRecvDuplicates() int32`

GetPacketsRecvDuplicates returns the PacketsRecvDuplicates field if non-nil, zero value otherwise.

### GetPacketsRecvDuplicatesOk

`func (o *ApiAPIToolPingOutput) GetPacketsRecvDuplicatesOk() (*int32, bool)`

GetPacketsRecvDuplicatesOk returns a tuple with the PacketsRecvDuplicates field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPacketsRecvDuplicates

`func (o *ApiAPIToolPingOutput) SetPacketsRecvDuplicates(v int32)`

SetPacketsRecvDuplicates sets PacketsRecvDuplicates field to given value.

### HasPacketsRecvDuplicates

`func (o *ApiAPIToolPingOutput) HasPacketsRecvDuplicates() bool`

HasPacketsRecvDuplicates returns a boolean if a field has been set.

### GetPacketsSent

`func (o *ApiAPIToolPingOutput) GetPacketsSent() int32`

GetPacketsSent returns the PacketsSent field if non-nil, zero value otherwise.

### GetPacketsSentOk

`func (o *ApiAPIToolPingOutput) GetPacketsSentOk() (*int32, bool)`

GetPacketsSentOk returns a tuple with the PacketsSent field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPacketsSent

`func (o *ApiAPIToolPingOutput) SetPacketsSent(v int32)`

SetPacketsSent sets PacketsSent field to given value.

### HasPacketsSent

`func (o *ApiAPIToolPingOutput) HasPacketsSent() bool`

HasPacketsSent returns a boolean if a field has been set.

### GetStdDevRtt

`func (o *ApiAPIToolPingOutput) GetStdDevRtt() interface{}`

GetStdDevRtt returns the StdDevRtt field if non-nil, zero value otherwise.

### GetStdDevRttOk

`func (o *ApiAPIToolPingOutput) GetStdDevRttOk() (*interface{}, bool)`

GetStdDevRttOk returns a tuple with the StdDevRtt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStdDevRtt

`func (o *ApiAPIToolPingOutput) SetStdDevRtt(v interface{})`

SetStdDevRtt sets StdDevRtt field to given value.

### HasStdDevRtt

`func (o *ApiAPIToolPingOutput) HasStdDevRtt() bool`

HasStdDevRtt returns a boolean if a field has been set.

### SetStdDevRttNil

`func (o *ApiAPIToolPingOutput) SetStdDevRttNil(b bool)`

 SetStdDevRttNil sets the value for StdDevRtt to be an explicit nil

### UnsetStdDevRtt
`func (o *ApiAPIToolPingOutput) UnsetStdDevRtt()`

UnsetStdDevRtt ensures that no value is present for StdDevRtt, not even an explicit nil

[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


