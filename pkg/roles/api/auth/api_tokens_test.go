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
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("api", ctx)
	role := api.New(inst)
	defer role.Stop()
	prov := auth.NewAuthProvider(role, inst)

	tests.PanicIfError(inst.KV().PutObj(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyTokens,
			tests.RandomString(),
		).String(),
		&types.Token{},
	))

	var output auth.APITokensGetOutput
	assert.NoError(t, prov.APITokensGet().Interact(ctx, struct{}{}, &output))
	assert.NotNil(t, output)
}

func TestAPITokensPut(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("api", ctx)
	role := api.New(inst)
	defer role.Stop()
	prov := auth.NewAuthProvider(role, inst)

	var output auth.APITokensPutOutput
	name := tests.RandomString()
	assert.NoError(t, prov.APITokensPut().Interact(ctx, auth.APITokensPutInput{
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
		&types.Token{
			Username: name,
		},
	)
}

func TestAPITokensDelete(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("api", ctx)
	role := api.New(inst)
	defer role.Stop()
	prov := auth.NewAuthProvider(role, inst)

	name := tests.RandomString()

	tests.PanicIfError(inst.KV().PutObj(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyTokens,
			name,
		).String(),
		&types.Token{},
	))

	assert.NoError(t, prov.APITokensDelete().Interact(ctx, auth.APITokensDeleteInput{
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
