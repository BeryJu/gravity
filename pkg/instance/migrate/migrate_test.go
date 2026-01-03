package migrate_test

import (
	"context"
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/instance/migrate"
	"beryju.io/gravity/pkg/instance/types"
	"beryju.io/gravity/pkg/storage"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestMigrate_ClusterVersion(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	ri := rootInst.ForRole("migrate", ctx)

	_, err := ri.KV().Put(
		ctx,
		ri.KV().Key(types.KeyInstance, "foo").String(),
		`{"version":"0.1.0"}`,
	)
	assert.NoError(t, err)
	_, err = ri.KV().Put(
		ctx,
		ri.KV().Key(types.KeyInstance, "bar").String(),
		`{"version":"0.15.0+foo"}`,
	)
	assert.NoError(t, err)
	// Invalid JSON
	_, err = ri.KV().Put(
		ctx,
		ri.KV().Key(types.KeyInstance, "baz").String(),
		`{`,
	)
	assert.NoError(t, err)
	// Invalid Version
	_, err = ri.KV().Put(
		ctx,
		ri.KV().Key(types.KeyInstance, "baz").String(),
		`{"version":"0.15.0++foo"}`,
	)
	assert.NoError(t, err)

	ct := 0
	ri.Migrator().AddMigration(&migrate.InlineMigration{
		MigrationName:     "test",
		ActivateOnVersion: migrate.MustParseConstraint("< 0.14.0"),
		HookFunc: func(ctx context.Context) (*storage.Client, error) {
			ct = 1
			return ri.KV(), nil
		},
	})
	_, err = ri.Migrator().Run(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 1, ct)
}

func TestMigrate(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	ri := rootInst.ForRole("migrate", ctx)
	ct := 0
	ri.Migrator().AddMigration(&migrate.InlineMigration{
		MigrationName:     "test",
		ActivateOnVersion: migrate.MustParseConstraint("> 0.0.0"),
		HookFunc: func(ctx context.Context) (*storage.Client, error) {
			ct = 1
			return ri.KV(), nil
		},
	})
	_, err := ri.Migrator().Run(ctx)
	assert.NoError(t, err)
	_, err = ri.KV().Put(
		ctx,
		ri.KV().Key("foo").String(),
		"bar",
	)
	assert.NoError(t, err)
	tests.AssertEtcd(t, ri.KV(), ri.KV().Key("foo"), "bar")
	assert.Equal(t, 1, ct)
}

func TestMigrate_Hook(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	ri := rootInst.ForRole("migrate", ctx)
	ct := 0
	ri.Migrator().AddMigration(&migrate.InlineMigration{
		MigrationName:     "test",
		ActivateOnVersion: migrate.MustParseConstraint("> 0.0.0"),
		HookFunc: func(ctx context.Context) (*storage.Client, error) {
			return ri.KV().WithHooks(storage.StorageHook{
				GetPre: func(ctx context.Context, key string, opts ...clientv3.OpOption) error {
					ct += 1
					return nil
				},
				PutPre: func(ctx context.Context, key, val string, opts ...clientv3.OpOption) error {
					ct += 1
					return nil
				},
			}), nil
		},
	})
	kv, err := ri.Migrator().Run(ctx)
	assert.NoError(t, err)
	_, err = kv.Put(
		ctx,
		kv.Key("foo").String(),
		"bar",
	)
	assert.NoError(t, err)
	tests.AssertEtcd(t, kv, ri.KV().Key("foo"), "bar")
	assert.Equal(t, 2, ct)
}

func TestMigrate_Cleanup(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	ri := rootInst.ForRole("migrate", ctx)
	ct := 0
	ri.Migrator().AddMigration(&migrate.InlineMigration{
		MigrationName:     "test",
		ActivateOnVersion: migrate.MustParseConstraint("< 0.1.0"),
		HookFunc: func(ctx context.Context) (*storage.Client, error) {
			return ri.KV(), nil
		},
		CleanupFunc: func(ctx context.Context) error {
			ct += 1
			return nil
		},
	})
	_, err := ri.Migrator().Run(ctx)
	assert.NoError(t, err)
	assert.Equal(t, 1, ct)
}
