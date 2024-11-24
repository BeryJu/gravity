package auth

import (
	"context"
	"encoding/json"
	"strings"

	"beryju.io/gravity/pkg/roles/api/types"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Permission struct {
	Path    string   `json:"path"`
	Methods []string `json:"methods"`
}

type User struct {
	ap *AuthProvider

	Username    string       `json:"-"`
	Password    string       `json:"password"`
	Permissions []Permission `json:"permissions"`
}

func (u *User) String() string {
	return u.Username
}

func (ap *AuthProvider) userFromKV(raw *mvccpb.KeyValue) (*User, error) {
	user := &User{
		ap: ap,
	}
	prefix := ap.inst.KV().Key(
		types.KeyRole,
		types.KeyUsers,
	).Prefix(true).String()
	user.Username = strings.TrimPrefix(string(raw.Key), prefix)
	err := json.Unmarshal(raw.Value, &user)
	if err != nil {
		return user, err
	}
	return user, nil
}

func (u *User) put(ctx context.Context, opts ...clientv3.OpOption) error {
	raw, err := json.Marshal(&u)
	if err != nil {
		return err
	}
	fullKey := u.ap.inst.KV().Key(
		types.KeyRole,
		types.KeyUsers,
		u.Username,
	).String()
	_, err = u.ap.inst.KV().Put(ctx, fullKey, string(raw), opts...)
	return err
}
