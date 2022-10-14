package auth_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/roles/api/auth"
	"beryju.io/gravity/pkg/roles/api/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestAPITokensGet(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("api")
	role := api.New(inst)
	ctx := tests.Context()
	prov := auth.NewAuthProvider(role, inst)

	inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyTokens,
			tests.RandomString(),
		).String(),
		tests.MustJSON(auth.Token{}),
	)

	var output auth.APITokensGetOutput
	assert.NoError(t, prov.APITokensGet().Interact(ctx, struct{}{}, &output))
	assert.NotNil(t, output)
}

func TestAPITokensPut(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("api")
	role := api.New(inst)
	prov := auth.NewAuthProvider(role, inst)

	var output auth.APITokensPutOutput
	name := tests.RandomString()
	assert.NoError(t, prov.APITokensPut().Interact(tests.Context(), auth.APITokensPutInput{
		Username: name,
	}, &output))

	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			types.KeyRole,
			types.KeyTokens,
			output.Key,
		),
		auth.Token{
			Username: name,
		},
	)
}

func TestAPITokensDelete(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("api")
	role := api.New(inst)
	ctx := tests.Context()
	prov := auth.NewAuthProvider(role, inst)

	name := tests.RandomString()

	inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyTokens,
			name,
		).String(),
		tests.MustJSON(auth.Token{}),
	)

	assert.NoError(t, prov.APITokensDelete().Interact(tests.Context(), auth.APITokensDeleteInput{
		Key: name,
	}, &struct{}{}))

	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			types.KeyRole,
			types.KeyTokens,
			name,
		),
	)
}
