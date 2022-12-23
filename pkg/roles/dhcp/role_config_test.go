package dhcp_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dhcp"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestAPIRoleConfigGet(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("dhcp")
	role := dhcp.New(inst)
	ctx := tests.Context()
	role.Start(ctx, []byte{})
	defer role.Stop()

	var output dhcp.APIRoleConfigOutput
	assert.NoError(t, role.APIRoleConfigGet().Interact(tests.Context(), struct{}{}, &output))
	assert.NotNil(t, output)
}

func TestAPIRoleConfigPut(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("dhcp")
	role := dhcp.New(inst)
	ctx := tests.Context()
	role.Start(ctx, []byte{})
	defer role.Stop()

	assert.NoError(t, role.APIRoleConfigPut().Interact(tests.Context(), dhcp.APIRoleConfigInput{
		Config: dhcp.RoleConfig{
			Port: 613,
		},
	}, &struct{}{}))
}
