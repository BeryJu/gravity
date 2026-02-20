package auth

import (
	"context"
	"strings"

	"beryju.io/gravity/pkg/roles/api/types"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	AuthorizationHeader = "Authorization"
	BearerType          = "bearer"
)

func (ap *AuthProvider) putToken(t *types.Token, ctx context.Context, opts ...clientv3.OpOption) error {
	fullKey := ap.inst.KV().Key(
		types.KeyRole,
		types.KeyTokens,
		t.Key,
	).String()
	_, err := ap.inst.KV().PutObj(ctx, fullKey, t, opts...)
	return err
}

func (ap *AuthProvider) tokenFromKV(raw *mvccpb.KeyValue) (*types.Token, error) {
	token := &types.Token{}
	prefix := ap.inst.KV().Key(
		types.KeyRole,
		types.KeyTokens,
	).Prefix(true).String()
	err := ap.inst.KV().Unmarshal(raw.Value, &token)
	if err != nil {
		return token, err
	}
	token.Key = strings.TrimPrefix(string(raw.Key), prefix)
	return token, nil
}
