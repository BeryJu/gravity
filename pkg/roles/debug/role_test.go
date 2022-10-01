package debug_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/debug"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestRole_Start_NoConfig(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("debug")
	role := debug.New(inst)
	assert.NotNil(t, role)
	ctx := tests.Context()
	assert.Equal(t, roles.ErrRoleNotConfigured, role.Start(ctx, []byte{}))
	defer role.Stop()
}

func TestRole_Start_EmptyConfig(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("debug")
	role := debug.New(inst)
	assert.NotNil(t, role)
	ctx := tests.Context()
	assert.Equal(t, roles.ErrRoleNotConfigured, role.Start(ctx, []byte("{}")))
	defer role.Stop()
}
