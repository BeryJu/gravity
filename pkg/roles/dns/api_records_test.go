package dns_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dns"
	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/roles/dns/utils"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestAPIRecordsGet(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	role := dns.New(inst)

	zone := utils.EnsureTrailingPeriod(tests.RandomString())
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			zone,
		).String(),
		tests.MustJSON(dns.Zone{}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			zone,
			"test",
			"A",
		).String(),
		tests.MustJSON(dns.Record{
			Data: "192.0.2.1",
		}),
	))

	var output dns.APIRecordsGetOutput
	assert.NoError(t, role.APIRecordsGet().Interact(ctx, dns.APIRecordsGetInput{
		Zone: zone,
	}, &output))
	assert.NotNil(t, output)
}

func TestAPIRecordsPut(t *testing.T) {
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
	assert.NoError(t, role.APIRecordsPut().Interact(ctx, dns.APIRecordsPutInput{
		Zone:     name,
		Hostname: "test",
		Type:     "A",
		Data:     "192.0.2.1",
	}, &struct{}{}))

	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			name,
			"test",
			"A",
		),
		dns.Record{
			Data: "192.0.2.1",
		},
	)
}

func TestAPIRecordsDelete(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	role := dns.New(inst)

	zone := utils.EnsureTrailingPeriod(tests.RandomString())
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			zone,
		).String(),
		tests.MustJSON(dns.Zone{}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			zone,
			"test",
			"A",
		).String(),
		tests.MustJSON(dns.Record{
			Data: "192.0.2.1",
		}),
	))

	assert.NoError(t, role.APIRecordsDelete().Interact(ctx, dns.APIRecordsDeleteInput{
		Zone:     zone,
		Hostname: "test",
		Type:     "A",
	}, &struct{}{}))

	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			zone,
			"test",
			"A",
		),
	)
}
