package instance_test

import (
	"testing"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestAPIInstanceInfo(t *testing.T) {
	rootInst := instance.New()

	var output instance.APIInstanceInfo
	assert.NoError(t, rootInst.APIInstanceInfo().Interact(tests.Context(), struct{}{}, &output))
	assert.NotNil(t, output)
	assert.Equal(t, output.Version, extconfig.Version)
}

func TestAPIInstances(t *testing.T) {
	rootInst := instance.New()

	var output instance.APIInstancesOutput
	assert.NoError(t, rootInst.APIInstances().Interact(tests.Context(), struct{}{}, &output))
	assert.NotNil(t, output)
}
