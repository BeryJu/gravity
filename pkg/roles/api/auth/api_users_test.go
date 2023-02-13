package auth_test

import (
	"context"
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/roles/api/auth"
	"beryju.io/gravity/pkg/roles/api/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/gorilla/sessions"
	"github.com/stretchr/testify/assert"
)

func TestAPIUsersGet(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("api")
	role := api.New(inst)
	ctx := tests.Context()
	prov := auth.NewAuthProvider(role, inst)

	inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyUsers,
			tests.RandomString(),
		).String(),
		tests.MustJSON(auth.User{}),
	)

	var output auth.APIUsersGetOutput
	assert.NoError(t, prov.APIUsersGet().Interact(ctx, auth.APIUsersGetInput{}, &output))
	assert.NotNil(t, output)
}

func TestAPIUsersPut(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("api")
	role := api.New(inst)
	prov := auth.NewAuthProvider(role, inst)

	name := tests.RandomString()
	password := tests.RandomString()

	assert.NoError(t, prov.APIUsersPut().Interact(tests.Context(), auth.APIUsersPutInput{
		Username: name,
		Password: password,
	}, &struct{}{}))

	_, err := inst.KV().Get(
		tests.Context(),
		inst.KV().Key(
			types.KeyRole,
			types.KeyUsers,
			name,
		).String(),
	)
	assert.NoError(t, err)

	var loginOutput auth.APILoginOutput
	sess := role.SessionStore()
	ctx := context.WithValue(tests.Context(), types.RequestSession, sessions.NewSession(sess, types.SessionName))
	assert.NoError(t, prov.APILogin().Interact(ctx, &auth.APILoginInput{
		Username: name,
		Password: password,
	}, &loginOutput))
	assert.True(t, loginOutput.Successful)
}

func TestAPIUsersDelete(t *testing.T) {
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
			types.KeyUsers,
			name,
		).String(),
		tests.MustJSON(auth.User{}),
	)

	assert.NoError(t, prov.APIUsersDelete().Interact(tests.Context(), auth.APIUsersDeleteInput{
		Username: name,
	}, &struct{}{}))

	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			types.KeyRole,
			types.KeyUsers,
			name,
		),
	)
}
