package auth

import (
	"context"
	"encoding/base64"
	"fmt"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/gorilla/securecookie"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (ap *AuthProvider) createDefaultUser() {
	users, err := ap.inst.KV().Get(
		context.Background(),
		ap.inst.KV().Key(
			types.KeyRole,
			types.KeyUsers,
		).Prefix(true).String(),
		clientv3.WithPrefix(),
	)
	if err != nil {
		ap.log.WithError(err).Warning("failed to check for users")
		return
	}
	if len(users.Kvs) > 0 {
		return
	}
	password := base64.RawStdEncoding.EncodeToString(securecookie.GenerateRandomKey(32))
	err = ap.CreateUser(context.Background(), "admin", password)
	if err != nil {
		ap.log.WithError(err).Warning("failed to create default user")
		return
	}
	text := `------------------------------------------------------------
Welcome to gravity! An admin user has been created for you.
Username: '%s'
Password: '%s'
Open 'http://%s:8008/' to start using Gravity!
--------------------------------------------------------------------
`
	fmt.Printf(text, "admin", password, extconfig.Get().Instance.IP)
}
