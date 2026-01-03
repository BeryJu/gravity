package monitoring_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/monitoring"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestAPIRoleConfigGet(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("monitoring", ctx)
	role := monitoring.New(inst)
	tests.PanicIfError(role.Start(ctx, []byte{}))
	defer role.Stop()

	var output monitoring.APIRoleConfigOutput
	assert.NoError(t, role.APIRoleConfigGet().Interact(ctx, struct{}{}, &output))
	assert.NotNil(t, output)
}

func TestAPIRoleConfigPut(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("monitoring", ctx)
	role := monitoring.New(inst)
	tests.PanicIfError(role.Start(ctx, []byte{}))
	defer role.Stop()

	assert.NoError(t, role.APIRoleConfigPut().Interact(ctx, monitoring.APIRoleConfigInput{
		Config: monitoring.RoleConfig{
			Port: 1234,
		},
	}, &struct{}{}))
}
