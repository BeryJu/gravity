package api

import (
	"encoding/base64"
	"fmt"

	"beryju.io/gravity/api"
	"beryju.io/gravity/pkg/instance"
	roleAPI "beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/roles/api/auth"
	"beryju.io/gravity/pkg/roles/api/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/gorilla/securecookie"
)

func APIClient(rootInst *instance.Instance) (*api.APIClient, func()) {
	inst := rootInst.ForRole("api")
	role := roleAPI.New(inst)
	ctx := tests.Context()
	role.Start(ctx, []byte{})

	token := base64.RawStdEncoding.EncodeToString(securecookie.GenerateRandomKey(64))

	inst.KV().Put(ctx, inst.KV().Key(
		types.KeyRole,
		types.KeyTokens,
		token,
	).String(), tests.MustJSON(auth.Token{
		Key:      token,
		Username: "foo",
	}))

	config := api.NewConfiguration()
	config.Debug = true
	config.Scheme = "http"
	config.Host = "localhost:8008"
	config.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	return api.NewAPIClient(config), func() {
		role.Stop()
	}
}
