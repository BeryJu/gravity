package instance_test

import (
	"testing"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestAPIInstanceInfo(t *testing.T) {
	defer tests.Setup(t)()
	rootInst := instance.New()

	var output instance.APIInstanceInfo
	assert.NoError(t, rootInst.APIInstanceGet().Interact(tests.Context(), struct{}{}, &output))
	assert.NotNil(t, output)
	assert.Equal(t, output.Version, extconfig.Version)
}

func TestAPIInstances(t *testing.T) {
	defer tests.Setup(t)()
	rootInst := instance.New()

	var output instance.APIClusterInfoOutput
	assert.NoError(t, rootInst.APIClusterInfo().Interact(tests.Context(), struct{}{}, &output))
	assert.NotNil(t, output)
}
