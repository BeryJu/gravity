package auth

import (
	"context"
	"strings"

	"beryju.io/gravity/pkg/roles/api/types"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/protobuf/proto"
)

func (ap *AuthProvider) userFromKV(raw *mvccpb.KeyValue) (*types.User, error) {
	user := &types.User{}
	prefix := ap.inst.KV().Key(
		types.KeyRole,
		types.KeyUsers,
	).Prefix(true).String()
	user.Username = strings.TrimPrefix(string(raw.Key), prefix)
	err := proto.Unmarshal(raw.Value, user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (ap *AuthProvider) putUser(u *types.User, ctx context.Context, opts ...clientv3.OpOption) error {
	raw, err := proto.Marshal(u)
	if err != nil {
		return err
	}
	fullKey := ap.inst.KV().Key(
		types.KeyRole,
		types.KeyUsers,
		u.Username,
	).String()
	_, err = ap.inst.KV().Put(ctx, fullKey, string(raw), opts...)
	return err
}
