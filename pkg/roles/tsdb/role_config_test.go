package tsdb_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/tsdb"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestAPIRoleConfigGet(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("tsdb", ctx)
	role := tsdb.New(inst)
	tests.PanicIfError(role.Start(ctx, []byte{}))
	defer role.Stop()

	var output tsdb.APIRoleConfigOutput
	assert.NoError(t, role.APIRoleConfigGet().Interact(ctx, struct{}{}, &output))
	assert.NotNil(t, output)
}

func TestAPIRoleConfigPut(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("tsdb", ctx)
	role := tsdb.New(inst)
	tests.PanicIfError(role.Start(ctx, []byte{}))
	defer role.Stop()

	assert.NoError(t, role.APIRoleConfigPut().Interact(ctx, tsdb.APIRoleConfigInput{
		Config: tsdb.RoleConfig{
			Enabled: false,
		},
	}, &struct{}{}))
}
