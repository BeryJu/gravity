package api

import (
	"encoding/base64"
	"fmt"
	"math/rand"

	"beryju.io/gravity/api"
	"beryju.io/gravity/pkg/instance"
	roleAPI "beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/roles/api/auth"
	"beryju.io/gravity/pkg/roles/api/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/gorilla/securecookie"
)

func APIClient(rootInst *instance.Instance) (*api.APIClient, func()) {
	port := rand.Intn(65535-1024) + 1024
	listen := fmt.Sprintf("localhost:%d", port)

	inst := rootInst.ForRole("api")
	role := roleAPI.New(inst)
	ctx := tests.Context()
	role.Start(ctx, []byte(tests.MustJSON(roleAPI.RoleConfig{
		ListenOverride: listen,
	})))

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
	config.Host = listen
	config.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	return api.NewAPIClient(config), func() {
		role.Stop()
	}
}
