package dhcp_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dhcp"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func RoleConfig() []byte {
	return []byte(tests.MustJSON(dhcp.RoleConfig{
		Port: 0,
	}))
}

func TestRoleStartNoConfig(t *testing.T) {
	defer tests.Setup(t)()
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)
	role := dhcp.New(inst)
	assert.NotNil(t, role)
	cfg := tests.MustJSON(&dhcp.RoleConfig{
		Port: 1067,
	})
	assert.Nil(t, role.Start(ctx, []byte(cfg)))
	defer role.Stop()
}
