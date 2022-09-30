package instance_test

import (
	"testing"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func Test_APIHandlerInfo(t *testing.T) {
	rootInst := instance.New()

	var output instance.APISystemInfo
	assert.NoError(t, rootInst.APIHandlerInfo().Interact(tests.Context(), struct{}{}, &output))
	assert.NotNil(t, output)
	assert.Equal(t, output.Version, extconfig.Version)
}

func Test_APIHandlerInstances(t *testing.T) {
	rootInst := instance.New()

	var output instance.APIInstancesOutput
	assert.NoError(t, rootInst.APIHandlerInstances().Interact(tests.Context(), struct{}{}, &output))
	assert.NotNil(t, output)
}
