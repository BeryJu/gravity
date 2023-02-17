package debug_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/debug"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestRoleStartNoConfig(t *testing.T) {
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("debug", ctx)
	role := debug.New(inst)
	assert.NotNil(t, role)
	assert.NoError(t, role.Start(ctx, []byte{}))
	defer role.Stop()
}

func TestRoleStartEmptyConfig(t *testing.T) {
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("debug", ctx)
	role := debug.New(inst)
	assert.NotNil(t, role)
	assert.NoError(t, role.Start(ctx, []byte("{}")))
	defer role.Stop()
}
