package dhcp_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dhcp"
	"beryju.io/gravity/pkg/roles/dhcp/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func testLease() dhcp.Lease {
	return dhcp.Lease{
		Identifier: tests.RandomString(),
		Address:    "192.0.2.1",
		Hostname:   "gravity.home.arpa",
	}
}

func TestAPILeasesGet(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)
	role := dhcp.New(inst)

	scope := testScope()
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyScopes,
			scope.Name,
		).String(),
		tests.MustJSON(scope),
	))
	lease := testLease()
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyLeases,
			lease.Identifier,
		).String(),
		tests.MustJSON(lease),
	))

	var output dhcp.APILeasesGetOutput
	assert.NoError(t, role.APILeasesGet().Interact(ctx, dhcp.APILeasesGetInput{
		ScopeName: scope.Name,
	}, &output))
	assert.NotNil(t, output)
}

func TestAPILeasesPut(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)
	role := dhcp.New(inst)

	scope := testScope()
	name := tests.RandomString()
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyScopes,
			scope.Name,
		).String(),
		tests.MustJSON(scope),
	))
	assert.NoError(t, role.APILeasesPut().Interact(ctx, dhcp.APILeasesPutInput{
		Identifier: name,
		Scope:      scope.Name,
		Address:    "192.0.2.1",
		Hostname:   "gravity.home.arpa",
	}, &struct{}{}))

	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			types.KeyRole,
			types.KeyLeases,
			name,
		),
		dhcp.Lease{
			ScopeKey: scope.Name,
			Address:  "192.0.2.1",
			Hostname: "gravity.home.arpa",
		},
	)
}

func TestAPILeasesDelete(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)
	role := dhcp.New(inst)

	scope := testScope()
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyScopes,
			scope.Name,
		).String(),
		tests.MustJSON(scope),
	))
	lease := testLease()
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyLeases,
			lease.Identifier,
		).String(),
		tests.MustJSON(lease),
	))

	assert.NoError(t, role.APILeasesDelete().Interact(ctx, dhcp.APILeasesDeleteInput{
		Scope:      scope.Name,
		Identifier: lease.Identifier,
	}, &struct{}{}))

	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			types.KeyRole,
			types.KeyLeases,
			lease.Identifier,
		),
	)
}
