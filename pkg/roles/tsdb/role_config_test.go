package tsdb_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/tsdb"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestAPIRoleConfigGet(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("tsdb")
	role := tsdb.New(inst)
	ctx := tests.Context()
	role.Start(ctx, []byte{})
	defer role.Stop()

	var output tsdb.APIRoleConfigOutput
	assert.NoError(t, role.APIRoleConfigGet().Interact(tests.Context(), struct{}{}, &output))
	assert.NotNil(t, output)
}

func TestAPIRoleConfigPut(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("tsdb")
	role := tsdb.New(inst)
	ctx := tests.Context()
	role.Start(ctx, []byte{})
	defer role.Stop()

	assert.NoError(t, role.APIRoleConfigPut().Interact(tests.Context(), tsdb.APIRoleConfigInput{
		Config: tsdb.RoleConfig{
			Enabled: false,
		},
	}, &struct{}{}))
}
