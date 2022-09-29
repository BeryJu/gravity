package monitoring_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/monitoring"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func Test_APIHandlerRoleConfigGet(t *testing.T) {
	rootInst := instance.NewInstance()
	inst := rootInst.ForRole("monitoring")
	role := monitoring.New(inst)
	ctx := tests.Context()
	role.Start(ctx, []byte{})

	var output monitoring.RoleMonitoringConfigOutput
	assert.NoError(t, role.APIHandlerRoleConfigGet().Interact(tests.Context(), struct{}{}, &output))
	assert.NotNil(t, output)
}

func Test_APIHandlerRoleConfigPut(t *testing.T) {
	rootInst := instance.NewInstance()
	inst := rootInst.ForRole("monitoring")
	role := monitoring.New(inst)
	ctx := tests.Context()
	role.Start(ctx, []byte{})

	var output struct{}
	assert.NoError(t, role.APIHandlerRoleConfigPut().Interact(tests.Context(), monitoring.RoleMonitoringConfigInput{
		Config: monitoring.RoleConfig{
			Port: 1234,
		},
	}, &output))
}
