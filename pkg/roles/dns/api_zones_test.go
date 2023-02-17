package dns_test

import (
	"strings"
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dns"
	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestAPIZonesGet(t *testing.T) {
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	role := dns.New(inst)

	inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			"test.",
		).String(),
		tests.MustJSON(dns.Zone{}),
	)

	var output dns.APIZonesGetOutput
	assert.NoError(t, role.APIZonesGet().Interact(ctx, dns.APIZonesGetInput{}, &output))
	assert.NotNil(t, output)
}

func TestAPIZonesPut(t *testing.T) {
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	role := dns.New(inst)

	name := tests.RandomString() + "."
	assert.NoError(t, role.APIZonesPut().Interact(ctx, dns.APIZonesPutInput{
		Name:          strings.TrimSuffix(name, "."),
		Authoritative: true,
		HandlerConfigs: []map[string]string{
			{
				"type": "etcd",
			},
		},
	}, &struct{}{}))

	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			name,
		),
		dns.Zone{
			Authoritative: true,
			HandlerConfigs: []map[string]string{
				{
					"type": "etcd",
				},
			},
		},
	)
}

func TestAPIZonesDelete(t *testing.T) {
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	role := dns.New(inst)

	name := tests.RandomString() + "."

	inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			name,
		).String(),
		tests.MustJSON(dns.Zone{}),
	)

	assert.NoError(t, role.APIZonesDelete().Interact(ctx, dns.APIZonesDeleteInput{
		Zone: name,
	}, &struct{}{}))

	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			name,
		),
	)
}
