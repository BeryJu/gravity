package tests

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"beryju.io/gravity/api"
	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	roleAPI "beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/roles/api/auth"
	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/gorilla/securecookie"
)

func APIClient(rootInst *instance.Instance) (*api.APIClient, func()) {
	ctx := Context()
	inst := rootInst.ForRole("api", ctx)
	role := roleAPI.New(inst)
	PanicIfError(role.Start(ctx, []byte(MustJSON(roleAPI.RoleConfig{
		ListenOverride: "localhost:8008",
	}))))

	username := base64.RawStdEncoding.EncodeToString(securecookie.GenerateRandomKey(32))
	token := base64.RawStdEncoding.EncodeToString(securecookie.GenerateRandomKey(64))

	PanicIfError(inst.KV().Put(ctx, inst.KV().Key(
		types.KeyRole,
		types.KeyUsers,
		username,
	).String(), MustJSON(&types.User{
		Username: username,
		Permissions: []*types.Permission{
			{
				Path:    "/*",
				Methods: []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodHead, http.MethodDelete},
			},
		},
	})))
	PanicIfError(inst.KV().Put(ctx, inst.KV().Key(
		types.KeyRole,
		types.KeyTokens,
		token,
	).String(), MustJSON(auth.Token{
		Key:      token,
		Username: username,
	})))

	config := api.NewConfiguration()
	config.Debug = true
	config.Scheme = "http"
	config.Host = "localhost:8008"
	config.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	config.UserAgent = fmt.Sprintf("gravity-testing/%s", extconfig.FullVersion())
	return api.NewAPIClient(config), func() {
		role.Stop()
	}
}
