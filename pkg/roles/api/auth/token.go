package auth

import (
	"context"
	"encoding/json"
	"strings"

	"beryju.io/gravity/pkg/roles/api/types"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

const (
	AuthorizationHeader = "Authorization"
	BearerType          = "bearer"
)

type Token struct {
	Key string `json:"-"`

	Username string `json:"username"`

	ap *AuthProvider
}

func (token *Token) put(ctx context.Context, opts ...clientv3.OpOption) error {
	raw, err := json.Marshal(&token)
	if err != nil {
		return err
	}
	fullKey := token.ap.inst.KV().Key(
		types.KeyRole,
		types.KeyTokens,
		token.Key,
	).String()
	_, err = token.ap.inst.KV().Put(ctx, fullKey, string(raw), opts...)
	return err
}

func (ap *AuthProvider) tokenFromKV(raw *mvccpb.KeyValue) (*Token, error) {
	token := &Token{
		ap: ap,
	}
	prefix := ap.inst.KV().Key(
		types.KeyRole,
		types.KeyTokens,
	).Prefix(true).String()
	token.Key = strings.TrimPrefix(string(raw.Key), prefix)
	err := json.Unmarshal(raw.Value, &token)
	if err != nil {
		return token, err
	}
	return token, nil
}
