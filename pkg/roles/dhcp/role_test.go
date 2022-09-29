package dhcp_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dhcp"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestRole_Start_NoConfig(t *testing.T) {
	rootInst := instance.NewInstance()
	inst := rootInst.ForRole("dhcp")
	role := dhcp.New(inst)
	assert.NotNil(t, role)
	ctx := tests.Context()
	role.Start(ctx, []byte{})
	role.Stop()
}
