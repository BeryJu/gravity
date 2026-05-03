package dns_test

import (
	"testing"
	"time"

	d "github.com/miekg/dns"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dns"
	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestRoleDNS_IPForwarder_v4(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			types.DNSRootZone,
		).String(),
		tests.MustJSON(dns.Zone{
			HandlerConfigs: []map[string]interface{}{
				{
					"type": "forward_ip",
					"to":   "127.0.0.1:1053",
				},
			},
		}),
	))
	role := dns.New(inst)
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, RoleConfig()))
	defer role.Stop()

	AssertDNS(t, role, []d.Question{
		{
			Name:   "gravity.beryju.io.",
			Qtype:  d.TypeA,
			Qclass: d.ClassINET,
		},
	}, "gravity.beryju.io.	3600	IN	A	10.0.0.1")
}

func TestRoleDNS_IPForwarder_v4_Cache(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			types.DNSRootZone,
		).String(),
		tests.MustJSON(dns.Zone{
			HandlerConfigs: []map[string]interface{}{
				{
					"type":      "forward_ip",
					"to":        "127.0.0.1:1053",
					"cache_ttl": "-2",
				},
			},
		}),
	))
	role := dns.New(inst)
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, RoleConfig()))
	defer role.Stop()

	AssertDNS(t, role, []d.Question{
		{
			Name:   "gravity.beryju.io.",
			Qtype:  d.TypeA,
			Qclass: d.ClassINET,
		},
	}, "gravity.beryju.io.	3600	IN	A	10.0.0.1")

	// We don't have a signal for when a record is persisted to the cache
	// so wait for things to settle
	time.Sleep(3 * time.Second)
	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			types.DNSRootZone,
			"gravity.beryju.io",
			types.DNSRecordTypeA,
			"0",
		),
		dns.Record{
			Data: "10.0.0.1",
		},
	)
}

// Test caching of a DNS Record in the apex of a zone
func TestRoleDNS_IPForwarder_v4_Cache_Zone(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			"gravity.beryju.io.",
		).String(),
		tests.MustJSON(dns.Zone{
			HandlerConfigs: []map[string]interface{}{
				{
					"type":      "forward_ip",
					"to":        []string{"127.0.0.1:1053"},
					"cache_ttl": -2,
				},
			},
		}),
	))
	role := dns.New(inst)
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, RoleConfig()))
	defer role.Stop()

	AssertDNS(t, role, []d.Question{
		{
			Name:   "gravity.beryju.io.",
			Qtype:  d.TypeA,
			Qclass: d.ClassINET,
		},
	}, "gravity.beryju.io.	3600	IN	A	10.0.0.1")

	// We don't have a signal for when a record is persisted to the cache
	// so wait for things to settle
	time.Sleep(3 * time.Second)
	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			"gravity.beryju.io.",
			"@",
			types.DNSRecordTypeA,
			"0",
		),
		dns.Record{
			Data: "10.0.0.1",
		},
	)
}

func TestRoleDNS_IPForwarder_v6(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			types.DNSRootZone,
		).String(),
		tests.MustJSON(dns.Zone{
			HandlerConfigs: []map[string]interface{}{
				{
					"type": "forward_ip",
					"to":   "127.0.0.1:1053",
				},
			},
		}),
	))
	role := dns.New(inst)
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, RoleConfig()))
	defer role.Stop()

	AssertDNS(t, role, []d.Question{
		{
			Name:   "ipv6.t.gravity.beryju.io.",
			Qtype:  d.TypeAAAA,
			Qclass: d.ClassINET,
		},
	}, "ipv6.t.gravity.beryju.io.	3600	IN	AAAA	fe80::1")
}
