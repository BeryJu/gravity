package tsdb_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/tsdb"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestRoleStartNoConfig(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("tsdb")
	role := tsdb.New(inst)
	assert.NotNil(t, role)
	ctx := tests.Context()
	assert.Equal(t, roles.ErrRoleNotConfigured, role.Start(ctx, []byte{}))
	defer role.Stop()
}

func TestRoleStartEmptyConfig(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("tsdb")
	role := tsdb.New(inst)
	assert.NotNil(t, role)
	ctx := tests.Context()
	assert.Equal(t, roles.ErrRoleNotConfigured, role.Start(ctx, []byte("{}")))
	defer role.Stop()
}
