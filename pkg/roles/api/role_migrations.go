package api

import (
	"context"
	"net/http"
	"strings"

	"beryju.io/gravity/pkg/instance/migrate"
	"beryju.io/gravity/pkg/roles/api/types"
	"beryju.io/gravity/pkg/storage"
	"github.com/Masterminds/semver/v3"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (r *Role) RegisterMigrations() {
	r.i.Migrator().AddMigration(&migrate.InlineMigration{
		MigrationName: "api-add-default-perms",
		ActivateFunc:  func(v *semver.Version) bool { return true },
		HookFunc: func(ctx context.Context) (*storage.Client, error) {
			userPrefix := r.i.KV().Key(types.KeyRole, types.KeyUsers).Prefix(true).String()
			defaultPerms := []types.Permission{
				{
					Path:    "/*",
					Methods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodHead, http.MethodDelete},
				},
			}
			return r.i.KV().WithHooks(storage.StorageHook{
				GetPost: func(ctx context.Context, key string, res *clientv3.GetResponse, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
					shouldIntercept := res != nil && len(res.Kvs) > 0 && strings.HasPrefix(key, userPrefix)
					// If we're fetching a user, intercept the response
					if shouldIntercept {
						u := map[string]any{}
						err := r.i.KV().Unmarshal(res.Kvs[0].Value, &u)
						if err != nil {
							return res, nil
						}
						if _, set := u["permissions"]; !set {
							u["permissions"] = defaultPerms
						}
						v, err := r.i.KV().Marshal(u)
						if err != nil {
							return res, nil
						}
						res.Kvs[0].Value = v
					}
					return res, nil
				},
			}), nil
		},
	})

}
