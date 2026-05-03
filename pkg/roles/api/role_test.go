package api_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestRoleStartNoConfig(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("api", ctx)
	role := api.New(inst)
	defer role.Stop()
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, []byte{}))
}

func TestRoleStartEmptyConfig(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("api", ctx)
	role := api.New(inst)
	defer role.Stop()
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, []byte(tests.MustJSON(api.RoleConfig{}))))
}
