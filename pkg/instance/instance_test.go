package instance_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/instance/types"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func Test_Start(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("test")
	called := false
	inst.AddEventListener(types.EventTopicInstanceBootstrapped, func(ev *roles.Event) {
		defer rootInst.Stop()

		// Start API to trigger EventTopicAPIMuxSetup
		inst := rootInst.ForRole("api")
		role := api.New(inst)
		ctx := tests.Context()
		assert.Nil(t, role.Start(ctx, []byte{}))
		role.Stop()

		called = true
	})
	rootInst.Start()
	assert.True(t, called)
}
