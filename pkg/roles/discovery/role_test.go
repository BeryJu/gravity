package discovery_test

import (
	"testing"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/discovery"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestRoleStartNoConfig(t *testing.T) {
	tests.Setup(t)
	extconfig.Get().ListenOnlyMode = false
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("discovery", ctx)
	role := discovery.New(inst)
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, []byte{}))
	role.Stop()
}

func TestRoleStartEmptyConfig(t *testing.T) {
	tests.Setup(t)
	extconfig.Get().ListenOnlyMode = false
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("discovery", ctx)
	role := discovery.New(inst)
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, []byte("{}")))
	role.Stop()
}
