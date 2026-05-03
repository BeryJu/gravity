package dhcp_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dhcp"
	"beryju.io/gravity/pkg/roles/dhcp/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func testScope() dhcp.Scope {
	return dhcp.Scope{
		Name:       tests.RandomString(),
		SubnetCIDR: "10.200.0.0/24",
		Default:    true,
		Options: []*types.DHCPOption{
			{
				TagName: types.TagNameRouter,
				Value:   types.OptionValue("10.200.0.1/24"),
			},
		},
		IPAM: map[string]string{
			"range_start": "10.200.0.100",
			"range_end":   "10.200.0.250",
		},
		DNS: &dhcp.ScopeDNS{},
	}
}

func TestAPIScopesGet(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)
	role := dhcp.New(inst)

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyScopes,
			"test",
		).String(),
		tests.MustJSON(testScope()),
	))
	lease := testLease()
	lease.ScopeKey = "test"
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyLeases,
			lease.Identifier,
		).String(),
		tests.MustJSON(lease),
	))

	var output dhcp.APIScopesGetOutput
	assert.NoError(t, role.APIScopesGet().Interact(ctx, dhcp.APIScopesGetInput{}, &output))
	assert.NotNil(t, output)
}

func TestAPIScopesPut(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)
	role := dhcp.New(inst)

	name := tests.RandomString()
	assert.NoError(t, role.APIScopesPut().Interact(ctx, dhcp.APIScopesPutInput{
		Name:       name,
		SubnetCIDR: "10.200.0.0/24",
		Default:    true,
		Options: []*types.DHCPOption{
			{
				TagName: types.TagNameRouter,
				Value:   types.OptionValue("10.200.0.1/24"),
			},
		},
		IPAM: map[string]string{
			"range_start": "10.200.0.100",
			"range_end":   "10.200.0.150",
		},
		DNS: &dhcp.ScopeDNS{},
	}, &struct{}{}))

	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			types.KeyRole,
			types.KeyScopes,
			name,
		),
		dhcp.Scope{
			SubnetCIDR: "10.200.0.0/24",
			Default:    true,
			Options: []*types.DHCPOption{
				{
					TagName: types.TagNameRouter,
					Value:   types.OptionValue("10.200.0.1/24"),
				},
			},
			IPAM: map[string]string{
				"range_start": "10.200.0.100",
				"range_end":   "10.200.0.150",
			},
			DNS: &dhcp.ScopeDNS{},
		},
	)
}

func TestAPIScopesDelete(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)
	role := dhcp.New(inst)

	name := tests.RandomString()

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyScopes,
			name,
		).String(),
		tests.MustJSON(testScope()),
	))

	assert.NoError(t, role.APIScopesDelete().Interact(ctx, dhcp.APIScopesDeleteInput{
		Scope: name,
	}, &struct{}{}))

	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			types.KeyRole,
			types.KeyScopes,
			name,
		),
	)
}
