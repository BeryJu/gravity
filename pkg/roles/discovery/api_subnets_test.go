package discovery_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/discovery"
	"beryju.io/gravity/pkg/roles/discovery/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestAPISubnetsGet(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("discovery", ctx)
	role := discovery.New(inst)

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeySubnets,
			tests.RandomString(),
		).String(),
		tests.MustJSON(discovery.Subnet{}),
	))

	var output discovery.APISubnetsGetOutput
	assert.NoError(t, role.APISubnetsGet().Interact(ctx, discovery.APISubnetsGetInput{}, &output))
	assert.NotNil(t, output)
}

func TestAPISubnetsPut(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("discovery", ctx)
	role := discovery.New(inst)

	name := tests.RandomString()
	assert.NoError(t, role.APISubnetsPut().Interact(ctx, discovery.APISubnetsPutInput{
		Name: name,
	}, &struct{}{}))

	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			types.KeyRole,
			types.KeySubnets,
			name,
		),
		discovery.Subnet{},
	)
}

func TestAPISubnetsDelete(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("discovery", ctx)
	role := discovery.New(inst)

	name := tests.RandomString()

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeySubnets,
			name,
		).String(),
		tests.MustJSON(discovery.Subnet{}),
	))

	assert.NoError(t, role.APISubnetsDelete().Interact(ctx, discovery.APISubnetsDeleteInput{
		Name: name,
	}, &struct{}{}))

	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			types.KeyRole,
			types.KeySubnets,
			name,
		),
	)
}
