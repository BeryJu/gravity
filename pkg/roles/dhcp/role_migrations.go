package dhcp

import (
	"context"
	"encoding/json"
	"strings"

	"beryju.io/gravity/pkg/instance/migrate"
	"beryju.io/gravity/pkg/roles/dhcp/types"
	"beryju.io/gravity/pkg/storage"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

func (r *Role) migrateMoveInitial(ctx context.Context) {
	r.log.Info("Running initial move migration")
	pfx := r.i.KV().Key(types.KeyRole, types.KeyLegacyLeases).Prefix(true).String()
	res, err := r.i.KV().Get(ctx, pfx, clientv3.WithPrefix())
	if err != nil {
		r.log.Warn("failed to get legacy leases", zap.Error(err))
		return
	}
	// Copy from legacy prefix (global scope) to new prefix (scope-scoped)
	for _, kv := range res.Kvs {
		l := &Lease{}
		ident := strings.TrimPrefix(string(kv.Key), pfx)
		err = json.Unmarshal(kv.Value, &l)
		if err != nil {
			r.log.Warn("failed to parse lease", zap.Error(err))
			continue
		}
		_, err = r.i.KV().Put(
			ctx,
			r.i.KV().Key(types.KeyRole, types.KeyScopes, l.ScopeKey, ident).String(),
			string(kv.Value),
		)
		if err != nil {
			r.log.Warn("failed to migrate lease", zap.Error(err))
			continue
		}
	}
}

func (r *Role) migrateMoveBackground(ctx context.Context) {
	watchChan := r.i.KV().Watch(
		ctx,
		r.i.KV().Key(types.KeyRole, types.KeyLegacyLeases).Prefix(true).String(),
		clientv3.WithPrefix(),
	)
	type partialLease struct {
		ScopeKey string `json:"scopeKey"`
	}
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			pl := partialLease{}
			err := json.Unmarshal(event.Kv.Value, &pl)
			if err != nil {
				r.log.Warn("failed to parse partial lease", zap.Error(err))
				continue
			}
			ident := strings.Split(string(event.Kv.Key), "/")[2]
			newKey := r.i.KV().Key(types.KeyRole, types.KeyScopes, pl.ScopeKey, ident).String()
			switch event.Type {
			case clientv3.EventTypePut:
				_, err = r.i.KV().Put(ctx, newKey, string(event.Kv.Value))
			case clientv3.EventTypeDelete:
				_, err = r.i.KV().Delete(ctx, newKey)
			}
			if err != nil {
				r.log.Warn("failed to mirror legacy lease operation", zap.Error(err))
				continue
			}
		}
	}
}

func (r *Role) RegisterMigrations() {
	r.i.Migrator().AddMigration(&migrate.InlineMigration{
		MigrationName:     "dhcp-move",
		ActivateOnVersion: migrate.MustParseConstraint("< 0.17.0"),
		CleanupFunc: func(ctx context.Context) error {
			res, err := r.i.KV().Delete(ctx,
				r.i.KV().Key(types.KeyRole, types.KeyLegacyLeases).Prefix(true).String(),
				clientv3.WithPrefix())
			if err != nil {
				return err
			}
			r.log.Info("Successfully cleaned up old DHCP leases", zap.Int64("count", res.Deleted))
			return nil
		},
		HookFunc: func(ctx context.Context) (*storage.Client, error) {
			pureKV := r.i.KV()
			leasePrefix := r.i.KV().Key(types.KeyRole, types.KeyScopes).Prefix(true).String()

			r.migrateMoveInitial(ctx)
			go r.migrateMoveBackground(ctx)

			return r.i.KV().WithHooks(storage.StorageHook{
				PutPost: func(ctx context.Context, key, val string, res *clientv3.PutResponse, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
					relKey := strings.TrimPrefix(key, leasePrefix)
					parts := strings.Split(relKey, "/")
					shouldIntercept := len(parts) == 2
					if shouldIntercept {
						r.log.Debug("hooking DHCP lease write for migration", zap.String("key", key))
						leaseKey := pureKV.Key(
							types.KeyRole,
							types.KeyLegacyLeases,
							parts[1],
						).String()
						r.log.Debug("Writing lease to legacy key", zap.String("key", leaseKey))
						_, err := pureKV.Put(
							ctx,
							leaseKey,
							val,
							opts...,
						)
						if err != nil {
							return nil, err
						}
					}
					return res, nil
				},
			}), nil
		},
	})

}
