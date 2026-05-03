package instance_test

import (
	"encoding/base64"
	"testing"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/instance/types"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestStart(t *testing.T) {
	tests.Setup(t)
	called := false
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("test", ctx)
	inst.AddEventListener(types.EventTopicInstanceBootstrapped, func(ev *roles.Event) {
		defer rootInst.Stop()

		// Start API to trigger EventTopicAPIMuxSetup
		role := rootInst.Role("api").(*api.Role)
		assert.Nil(t, role.Start(ctx, []byte{}))
		role.Stop()

		called = true
	})
	rootInst.Start()
	assert.True(t, called)
}

func TestFirstStart(t *testing.T) {
	tests.Setup(t)
	called := false
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("test", ctx)
	inst.AddEventListener(types.EventTopicInstanceFirstStart, func(ev *roles.Event) {
		defer rootInst.Stop()

		called = true
	})
	rootInst.Start()
	assert.True(t, called)
}

func TestWatch(t *testing.T) {
	tests.Setup(t)
	called := 0
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("test", ctx)

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyInstance,
			extconfig.Get().Instance.Identifier,
			"roles",
		).String(),
		"api",
	))

	inst.AddEventListener(types.EventTopicRolesStarted, func(ev *roles.Event) {
		// once all roles are started, write a new config to trigger a restart
		tests.PanicIfError(inst.KV().Put(
			ctx,
			inst.KV().Key(
				types.KeyInstance,
				types.KeyRole,
				"api",
			).String(),
			tests.MustJSON(api.RoleConfig{
				CookieSecret: base64.StdEncoding.EncodeToString([]byte("bar")),
				Port:         8008,
			}),
		))
	})
	inst.AddEventListener(types.EventTopicRoleStarted, func(ev *roles.Event) {
		if ev.Payload.Data["role"] != "api" {
			return
		}
		called += 1
		if called == 2 {
			rootInst.Stop()
		}
	})
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyInstance,
			types.KeyRole,
			"api",
		).String(),
		tests.MustJSON(api.RoleConfig{
			CookieSecret: base64.StdEncoding.EncodeToString([]byte("foo")),
			Port:         8008,
		}),
	))

	rootInst.Start()
	assert.Equal(t, 2, called)
}
