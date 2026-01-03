package dhcp_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dhcp"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestAPIRoleConfigGet(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)
	role := dhcp.New(inst)
	tests.PanicIfError(role.Start(ctx, RoleConfig()))
	defer role.Stop()

	var output dhcp.APIRoleConfigOutput
	assert.NoError(t, role.APIRoleConfigGet().Interact(ctx, struct{}{}, &output))
	assert.NotNil(t, output)
}

func TestAPIRoleConfigPut(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)
	role := dhcp.New(inst)
	tests.PanicIfError(role.Start(ctx, RoleConfig()))
	defer role.Stop()

	assert.NoError(t, role.APIRoleConfigPut().Interact(ctx, dhcp.APIRoleConfigInput{
		Config: dhcp.RoleConfig{
			Port: 613,
		},
	}, &struct{}{}))
}
