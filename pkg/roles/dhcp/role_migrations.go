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

func (r *Role) migrateMoveInitial() {
	r.log.Info("Running initial move migration")
	res, err := r.i.KV().Get(
		r.ctx,
		r.i.KV().Key(types.KeyRole, types.KeyLegacyLeases).Prefix(true).String(),
		clientv3.WithPrefix(),
	)
	if err != nil {
		r.log.Warn("failed to get legacy leases", zap.Error(err))
		return
	}
	// Copy from legacy prefix (global scope) to new prefix (scope-scoped)
	for _, kv := range res.Kvs {
		l := &Lease{}
		err = json.Unmarshal(kv.Value, &l)
		if err != nil {
			r.log.Warn("failed to parse lease", zap.Error(err))
			continue
		}
		_, err = r.i.KV().Put(
			r.ctx,
			r.i.KV().Key(types.KeyRole, types.KeyScopes, l.ScopeKey, l.Identifier).String(),
			string(kv.Value),
		)
		if err != nil {
			r.log.Warn("failed to migrate lease", zap.Error(err))
			continue
		}
	}
}

// func (r *Role) migrateMoveBackground() {
// 	watchChan := r.i.KV().Watch(
// 		r.ctx,
// 		r.i.KV().Key(types.KeyRole, types.KeyLegacyLeases).Prefix(true).String(),
// 		clientv3.WithPrefix(),
// 	)
// 	for watchResp := range watchChan {
// 		for _, event := range watchResp.Events {
// 			switch event.Type {
// 			case clientv3.EventTypeDelete:
// 				r.i.KV().Delete(r.ctx)
// 			}
// 		}
// 	}
// }

func (r *Role) RegisterMigrations() {
	r.i.Migrator().AddMigration(&migrate.InlineMigration{
		MigrationName:     "dhcp-move",
		ActivateOnVersion: migrate.MustParseConstraint("< 0.17.0"),
		HookFunc: func(ctx context.Context) (*storage.Client, error) {
			pureKV := r.i.KV()
			leasePrefix := r.i.KV().Key(types.KeyRole, types.KeyScopes).Prefix(true).String()

			r.migrateMoveInitial()

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
