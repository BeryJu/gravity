package dns_test

import (
	"strings"
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dns"
	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/roles/dns/utils"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestAPIZonesGet(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	role := dns.New(inst)

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			"test.",
		).String(),
		tests.MustJSON(dns.Zone{}),
	))

	var output dns.APIZonesGetOutput
	assert.NoError(t, role.APIZonesGet().Interact(ctx, dns.APIZonesGetInput{}, &output))
	assert.NotNil(t, output)
}

func TestAPIZonesPut(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	role := dns.New(inst)

	name := utils.EnsureTrailingPeriod(tests.RandomString())
	assert.NoError(t, role.APIZonesPut().Interact(ctx, dns.APIZonesPutInput{
		Name:          strings.TrimSuffix(name, types.DNSSep),
		Authoritative: true,
		HandlerConfigs: []map[string]interface{}{
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
			HandlerConfigs: []map[string]interface{}{
				{
					"type": "etcd",
				},
			},
		},
	)
}

func TestAPIZonesDelete(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	role := dns.New(inst)

	name := utils.EnsureTrailingPeriod(tests.RandomString())

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			name,
		).String(),
		tests.MustJSON(dns.Zone{}),
	))

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
