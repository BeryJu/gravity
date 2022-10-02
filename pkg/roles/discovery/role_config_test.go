package discovery_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/discovery"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestAPIRoleConfigGet(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("discovery")
	role := discovery.New(inst)
	ctx := tests.Context()
	role.Start(ctx, []byte{})
	defer role.Stop()

	var output discovery.APIRoleConfigOutput
	assert.NoError(t, role.APIRoleConfigGet().Interact(tests.Context(), struct{}{}, &output))
	assert.NotNil(t, output)
}

func TestAPIRoleConfigPut(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("discovery")
	role := discovery.New(inst)
	ctx := tests.Context()
	role.Start(ctx, []byte{})
	defer role.Stop()

	var output struct{}
	assert.NoError(t, role.APIRoleConfigPut().Interact(tests.Context(), discovery.APIRoleConfigInput{
		Config: discovery.RoleConfig{},
	}, &output))
}
