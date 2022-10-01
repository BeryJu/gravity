package monitoring_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/monitoring"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func Test_APIRoleConfigGet(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("monitoring")
	role := monitoring.New(inst)
	ctx := tests.Context()
	role.Start(ctx, []byte{})
	defer role.Stop()

	var output monitoring.APIRoleMonitoringConfigOutput
	assert.NoError(t, role.APIRoleConfigGet().Interact(tests.Context(), struct{}{}, &output))
	assert.NotNil(t, output)
}

func Test_APIRoleConfigPut(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("monitoring")
	role := monitoring.New(inst)
	ctx := tests.Context()
	role.Start(ctx, []byte{})
	defer role.Stop()

	var output struct{}
	assert.NoError(t, role.APIRoleConfigPut().Interact(tests.Context(), monitoring.APIRoleMonitoringConfigInput{
		Config: monitoring.RoleConfig{
			Port: 1234,
		},
	}, &output))
}
