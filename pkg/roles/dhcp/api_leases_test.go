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
	lease.ScopeKey = scope.Name
	lease.Expiry = -1
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyLeases,
			types.KeyReservations,
			scope.Name,
			lease.Identifier,
		).String(),
		tests.MustJSON(lease),
	))
	tests.PanicIfError(role.Start(ctx, RoleConfig()))
	defer role.Stop()

	var output dhcp.APILeasesGetOutput
	assert.NoError(t, role.APILeasesGet().Interact(ctx, dhcp.APILeasesGetInput{
		ScopeName:  scope.Name,
		Identifier: lease.Identifier,
	}, &output))
	assert.Len(t, output.Leases, 1)
	assert.Equal(t, lease.Identifier, output.Leases[0].Identifier)
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

func TestAPILeasesPutReservation(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)
	role := dhcp.New(inst)

	scope := testScope()
	identifier := tests.RandomString()
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(types.KeyRole, types.KeyScopes, scope.Name).String(),
		tests.MustJSON(scope),
	))
	assert.NoError(t, role.APILeasesPut().Interact(ctx, dhcp.APILeasesPutInput{
		Identifier: identifier,
		Scope:      scope.Name,
		Address:    "10.200.0.150",
		Expiry:     -1,
	}, &struct{}{}))

	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(types.KeyRole, types.KeyLeases, types.KeyReservations, scope.Name, identifier),
		dhcp.Lease{
			ScopeKey: scope.Name,
			Address:  "10.200.0.150",
			Expiry:   -1,
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
	lease.ScopeKey = scope.Name
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

func TestAPILeasesDeleteReservation(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)
	role := dhcp.New(inst)

	scope := testScope()
	lease := testLease()
	lease.ScopeKey = scope.Name
	lease.Expiry = -1
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(types.KeyRole, types.KeyScopes, scope.Name).String(),
		tests.MustJSON(scope),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(types.KeyRole, types.KeyLeases, types.KeyReservations, scope.Name, lease.Identifier).String(),
		tests.MustJSON(lease),
	))

	assert.NoError(t, role.APILeasesDelete().Interact(ctx, dhcp.APILeasesDeleteInput{
		Scope:      scope.Name,
		Identifier: lease.Identifier,
	}, &struct{}{}))

	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(types.KeyRole, types.KeyLeases, types.KeyReservations, scope.Name, lease.Identifier),
	)
}
