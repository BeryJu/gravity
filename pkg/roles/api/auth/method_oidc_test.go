package auth_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/roles/api/auth"
	"beryju.io/gravity/pkg/roles/api/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestAuthOIDC(t *testing.T) {
	tests.Setup(t)
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

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/auth/oidc", nil)
	role.Mux().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusFound, rr.Result().StatusCode)
	loc, _ := rr.Result().Location()
	assert.True(t, strings.HasPrefix(loc.String(), "http://127.0.0.1:5556/dex/auth"), loc.String())
}

func TestAuthOIDC_Token(t *testing.T) {
	tests.Setup(t)

	// get initial token from Dex
	// https://dexidp.io/docs/connectors/local/

	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("scope", "openid profile email")
	data.Set("username", "admin@example.com")
	data.Set("password", "password")

	req, err := http.NewRequest(http.MethodPost, "http://127.0.0.1:5556/dex/token", strings.NewReader(data.Encode()))
	if err != nil {
		panic(err)
	}
	req.Header.Set("Authorization", "Basic Z3Jhdml0eTowOGE4Njg0Yi1kYjg4LTRiNzMtOTBhOS0zY2QxNjYxZjU0NjY=")
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		panic(err)
	}
	type b struct {
		IDToken string `json:"id_token"`
	}
	bo := b{}
	err = json.NewDecoder(res.Body).Decode(&bo)
	if err != nil {
		panic(err)
	}

	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("api", ctx)

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyUsers,
			"admin@example.com",
		).String(),
		tests.MustJSON(auth.User{}),
	))

	role := api.New(inst)
	assert.NoError(t, role.Start(ctx, []byte(tests.MustJSON(api.RoleConfig{
		ListenOverride: tests.Listen(8008),
		OIDC: &types.OIDCConfig{
			Issuer:             "http://127.0.0.1:5556/dex",
			ClientID:           "gravity",
			ClientSecret:       "08a8684b-db88-4b73-90a9-3cd1661f5466",
			RedirectURL:        "http://localhost:8008/auth/oidc/callback",
			Scopes:             []string{"openid", "email"},
			TokenUsernameField: "email",
		},
	}))))
	defer role.Stop()

	// Actual test with the token we just got

	rr := httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+bo.IDToken)
	role.Mux().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	assert.JSONEq(t, tests.MustJSON(auth.APIMeOutput{
		Username:      "admin@example.com",
		Authenticated: true,
	}), rr.Body.String())

	// test with invalid token

	rr = httptest.NewRecorder()
	req, _ = http.NewRequest(http.MethodGet, "/api/v1/auth/me", nil)
	req.Header.Set("Authorization", "Bearer "+bo.IDToken+"foo")
	role.Mux().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
	assert.JSONEq(t, tests.MustJSON(auth.APIMeOutput{
		Username:      "",
		Authenticated: false,
	}), rr.Body.String())
}
