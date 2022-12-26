package instance_test

import (
	"encoding/base64"
	"testing"
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/instance/types"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestStart(t *testing.T) {
	called := false
	rootInst := instance.New()
	inst := rootInst.ForRole("test")
	inst.AddEventListener(types.EventTopicInstanceBootstrapped, func(ev *roles.Event) {
		defer rootInst.Stop()

		// Start API to trigger EventTopicAPIMuxSetup
		role := rootInst.Role("api").(*api.Role)
		ctx := tests.Context()
		assert.Nil(t, role.Start(ctx, []byte{}))
		role.Stop()

		called = true
	})
	rootInst.Start()
	assert.True(t, called)
}

func TestFirstStart(t *testing.T) {
	called := false
	rootInst := instance.New()
	inst := rootInst.ForRole("test")
	inst.KV().Delete(
		tests.Context(),
		inst.KV().Key(
			types.KeyCluster,
		).Prefix(true).String(),
		clientv3.WithPrefix(),
	)
	inst.AddEventListener(types.EventTopicInstanceFirstStart, func(ev *roles.Event) {
		defer rootInst.Stop()

		called = true
	})
	rootInst.Start()
	assert.True(t, called)
}

func TestWatch(t *testing.T) {
	called := 0
	rootInst := instance.New()
	inst := rootInst.ForRole("test")

	inst.KV().Put(
		tests.Context(),
		inst.KV().Key(
			types.KeyInstance,
			extconfig.Get().Instance.Identifier,
			"roles",
		).String(),
		"api",
	)

	inst.AddEventListener(types.EventTopicRoleStarted, func(ev *roles.Event) {
		if ev.Payload.Data["role"] != "api" {
			return
		}
		called += 1
		if called == 1 {
			go func() {
				// Yes this is a bit hacky, but we need to wait for the role to finish starting
				// and the etcd watcher to start
				time.Sleep(5 * time.Second)
				inst.KV().Put(
					tests.Context(),
					inst.KV().Key(
						types.KeyInstance,
						types.KeyRole,
						"api",
					).String(),
					tests.MustJSON(api.RoleConfig{
						CookieSecret: base64.StdEncoding.EncodeToString([]byte("bar")),
						Port:         8008,
					}),
				)
			}()
		}
		if called == 2 {
			defer rootInst.Stop()
		}
	})
	inst.KV().Put(
		tests.Context(),
		inst.KV().Key(
			types.KeyInstance,
			types.KeyRole,
			"api",
		).String(),
		tests.MustJSON(api.RoleConfig{
			CookieSecret: base64.StdEncoding.EncodeToString([]byte("foo")),
			Port:         8008,
		}),
	)

	rootInst.Start()
	assert.Equal(t, 2, called)
}
