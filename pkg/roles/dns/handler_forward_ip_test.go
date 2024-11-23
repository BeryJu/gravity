package dns_test

import (
	"net"
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
	defer tests.Setup(t)()
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			".",
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

	fw := NewNullDNSWriter()
	role.Handler(fw, &d.Msg{
		Question: []d.Question{
			{
				Name:   "gravity.beryju.io.",
				Qtype:  d.TypeA,
				Qclass: d.ClassINET,
			},
		},
	})
	ans := fw.Msg().Answer[0]
	assert.Equal(t, net.ParseIP("10.0.0.1").String(), ans.(*d.A).A.String())
}

func TestRoleDNS_IPForwarder_v4_Cache(t *testing.T) {
	defer tests.Setup(t)()
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			".",
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

	fw := NewNullDNSWriter()
	role.Handler(fw, &d.Msg{
		Question: []d.Question{
			{
				Name:   "gravity.beryju.io.",
				Qtype:  d.TypeA,
				Qclass: d.ClassINET,
			},
		},
	})
	ans := fw.Msg().Answer[0]
	assert.Equal(t, net.ParseIP("10.0.0.1").String(), ans.(*d.A).A.String())

	// We don't have a signal for when a record is persisted to the cache
	// so wait for things to settle
	time.Sleep(3 * time.Second)
	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			".",
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
	defer tests.Setup(t)()
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

	fw := NewNullDNSWriter()
	role.Handler(fw, &d.Msg{
		Question: []d.Question{
			{
				Name:   "gravity.beryju.io.",
				Qtype:  d.TypeA,
				Qclass: d.ClassINET,
			},
		},
	})
	ans := fw.Msg().Answer[0]
	assert.Equal(t, net.ParseIP("10.0.0.1").String(), ans.(*d.A).A.String())

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
	defer tests.Setup(t)()
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			".",
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

	fw := NewNullDNSWriter()
	role.Handler(fw, &d.Msg{
		Question: []d.Question{
			{
				Name:   "ipv6.t.gravity.beryju.io.",
				Qtype:  d.TypeAAAA,
				Qclass: d.ClassINET,
			},
		},
	})
	ans := fw.Msg().Answer[0]
	assert.Equal(t, net.ParseIP("fe80::1").String(), ans.(*d.AAAA).AAAA.String())
}
