package tsdb_test

import (
	"testing"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/roles/debug"
	"beryju.io/gravity/pkg/roles/tsdb"
	"beryju.io/gravity/pkg/roles/tsdb/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestRoleStartNoConfig(t *testing.T) {
	rootInst := instance.New()
	ctx := tests.Context()

	apiRole := api.New(rootInst.ForRole("api"))
	assert.NoError(t, apiRole.Start(ctx, []byte{}))
	defer apiRole.Stop()
	debugRole := debug.New(rootInst.ForRole("debug"))
	assert.NoError(t, debugRole.Start(ctx, []byte{}))
	defer debugRole.Stop()

	inst := rootInst.ForRole("tsdb")
	role := tsdb.New(inst)
	assert.NotNil(t, role)
	assert.NoError(t, role.Start(ctx, []byte{}))
	defer role.Stop()
}

func TestRoleStartEmptyConfig(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("tsdb")
	role := tsdb.New(inst)
	assert.NotNil(t, role)
	ctx := tests.Context()
	assert.NoError(t, role.Start(ctx, []byte("{}")))
	defer role.Stop()
}

func TestRoleStartNotEnabled(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("tsdb")
	role := tsdb.New(inst)
	assert.NotNil(t, role)
	ctx := tests.Context()
	assert.Error(t, roles.ErrRoleNotConfigured, role.Start(ctx, []byte(tests.MustJSON(tsdb.RoleConfig{
		Enabled: false,
	}))))
	defer role.Stop()
}

func TestRoleWrite(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("tsdb")
	role := tsdb.New(inst)
	ctx := tests.Context()
	assert.NoError(t, role.Start(ctx, []byte("{}")))
	defer role.Stop()
	nameBeforeWrite := tests.RandomString()
	nameSet := tests.RandomString()
	nameInc := tests.RandomString()

	inst.KV().Delete(
		ctx,
		inst.KV().Key(
			types.KeyRole,
		).Prefix(true).String(),
		clientv3.WithPrefix(),
	)

	inst.AddEventListener(types.EventTopicTSDBBeforeWrite, func(ev *roles.Event) {
		role.SetMetric(
			inst.KV().Key(nameBeforeWrite).String(),
			types.Metric{
				Value: 42,
			},
		)
	})
	inst.DispatchEvent(types.EventTopicTSDBSet, roles.NewEvent(ctx, map[string]interface{}{
		"key": nameSet,
		"value": types.Metric{
			Value: 43,
		},
	}))
	inst.DispatchEvent(types.EventTopicTSDBInc, roles.NewEvent(ctx, map[string]interface{}{
		"key": nameInc,
		"default": types.Metric{
			ResetOnWrite: true,
		},
	}))
	inst.DispatchEvent(types.EventTopicTSDBWrite, roles.NewEvent(ctx, map[string]interface{}{}))

	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			types.KeyRole, nameBeforeWrite,
			extconfig.Get().Instance.Identifier,
		).Prefix(true),
		42,
	)
	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			types.KeyRole, nameSet,
			extconfig.Get().Instance.Identifier,
		).Prefix(true),
		43,
	)
	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			types.KeyRole, nameInc,
			extconfig.Get().Instance.Identifier,
		).Prefix(true),
		1,
	)
}
