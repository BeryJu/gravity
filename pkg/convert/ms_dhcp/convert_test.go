package ms_dhcp_test

import (
	"encoding/base64"
	"fmt"
	"testing"

	"beryju.io/gravity/api"
	"beryju.io/gravity/pkg/convert/ms_dhcp"
	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	roleAPI "beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/roles/api/auth"
	"beryju.io/gravity/pkg/roles/api/types"
	"beryju.io/gravity/pkg/roles/dhcp"
	"beryju.io/gravity/pkg/tests"
	"github.com/gorilla/securecookie"
	"github.com/stretchr/testify/assert"
)

func APIClient() *api.APIClient {
	rootInst := instance.New()
	inst := rootInst.ForRole("api")

	token := base64.RawStdEncoding.EncodeToString(securecookie.GenerateRandomKey(64))

	inst.KV().Put(tests.Context(), inst.KV().Key(
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
	config.Host = extconfig.Get().Listen(8008)
	config.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", token))
	return api.NewAPIClient(config)
}

func TestDHCPImport(t *testing.T) {
	rootInst := instance.New()
	ctx := tests.Context()
	// Create DHCP role to register API routes
	dhcp.New(rootInst.ForRole("dhcp"))
	inst := rootInst.ForRole("api")
	role := roleAPI.New(inst)
	role.Start(ctx, []byte{})
	defer role.Stop()

	files := []string{
		"./test_a.xml",
		"./test_b.xml",
		"./test_c.xml",
	}

	for _, file := range files {
		t.Run(file, func(t *testing.T) {
			c, err := ms_dhcp.New(APIClient(), file)
			assert.NoError(t, err)
			errors := c.Run(ctx)
			assert.Equal(t, []error{}, errors)
		})
	}
}
