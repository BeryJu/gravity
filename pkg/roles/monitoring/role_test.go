package monitoring_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/monitoring"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestRole_Start_NoConfig(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("monitoring")
	role := monitoring.New(inst)
	assert.NotNil(t, role)
	ctx := tests.Context()
	assert.Nil(t, role.Start(ctx, []byte{}))
	role.Stop()
}

func TestRole_Start_EmptyConfig(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("monitoring")
	role := monitoring.New(inst)
	assert.NotNil(t, role)
	ctx := tests.Context()
	assert.Nil(t, role.Start(ctx, []byte("{}")))
	role.Stop()
}
