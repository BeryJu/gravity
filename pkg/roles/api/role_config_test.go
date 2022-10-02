package api_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestAPIRoleConfigGet(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("api")
	role := api.New(inst)
	ctx := tests.Context()
	role.Start(ctx, []byte{})
	defer role.Stop()

	var output api.APIRoleConfigOutput
	assert.NoError(t, role.APIRoleConfigGet().Interact(tests.Context(), struct{}{}, &output))
	assert.NotNil(t, output)
}

func TestAPIRoleConfigPut(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("api")
	role := api.New(inst)
	ctx := tests.Context()
	role.Start(ctx, []byte{})
	defer role.Stop()

	var output struct{}
	assert.NoError(t, role.APIRoleConfigPut().Interact(tests.Context(), api.APIRoleConfigInput{
		Config: api.RoleConfig{
			Port: 8013,
		},
	}, &output))
}
