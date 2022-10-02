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
	inst := rootInst.ForRole("dns")
	role := dns.New(inst)
	ctx := tests.Context()

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
	assert.NoError(t, role.APIZonesGet().Interact(ctx, struct{}{}, &output))
	assert.NotNil(t, output)
}

func TestAPIZonesPut(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("dns")
	role := dns.New(inst)

	var output struct{}
	name := tests.RandomString() + "."
	assert.NoError(t, role.APIZonesPut().Interact(tests.Context(), dns.APIZonesPutInput{
		Name:          strings.TrimSuffix(name, "."),
		Authoritative: true,
		HandlerConfigs: []map[string]string{
			{
				"type": "etcd",
			},
		},
	}, &output))

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
	inst := rootInst.ForRole("dns")
	role := dns.New(inst)
	ctx := tests.Context()

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

	assert.NoError(t, role.APIZonesDelete().Interact(tests.Context(), dns.APIZonesDeleteInput{
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
