package auth

import (
	"context"
	"strings"

	"beryju.io/gravity/pkg/roles/api/types"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (ap *AuthProvider) userFromKV(raw *mvccpb.KeyValue) (*types.User, error) {
	user := &types.User{}
	prefix := ap.inst.KV().Key(
		types.KeyRole,
		types.KeyUsers,
	).Prefix(true).String()
	err := ap.inst.KV().Unmarshal(raw.Value, user)
	if err != nil {
		return user, err
	}
	user.Username = strings.TrimPrefix(string(raw.Key), prefix)
	return user, nil
}

func (ap *AuthProvider) putUser(u *types.User, ctx context.Context, opts ...clientv3.OpOption) error {
	fullKey := ap.inst.KV().Key(
		types.KeyRole,
		types.KeyUsers,
		u.Username,
	).String()
	_, err := ap.inst.KV().PutObj(ctx, fullKey, u, opts...)
	return err
}
