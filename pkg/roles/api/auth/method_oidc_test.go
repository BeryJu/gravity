package auth_test

import (
	"fmt"
	"net/http"
	"strings"
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/roles/api/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestAuthOIDC(t *testing.T) {
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("api", ctx)

	role := api.New(inst)
	assert.NoError(t, role.Start(ctx, []byte(tests.MustJSON(api.RoleConfig{
		ListenOverride: tests.Listen(8008),
		OIDC: &types.OIDCConfig{
			Issuer:       "http://127.0.0.1:5556/dex",
			ClientID:     "gravity",
			ClientSecret: "08a8684b-db88-4b73-90a9-3cd1661f5466",
			RedirectURL:  "http://localhost:8008/auth/oidc/callback",
			Scopes:       []string{"openid", "email"},
		},
	}))))
	defer role.Stop()
	tests.WaitForPort(8008)

	res, err := http.DefaultClient.Get(fmt.Sprintf("http://%s/auth/oidc", tests.Listen(8008)))
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, res.StatusCode)
	assert.True(t, strings.HasPrefix(res.Request.URL.String(), "http://127.0.0.1:5556/dex/auth/local"))
}
