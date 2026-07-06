package dns_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dns"
	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/tests"
	d "github.com/miekg/dns"
	"github.com/stretchr/testify/assert"
)

func setupMemoryZone(t *testing.T) *dns.Role {
	t.Helper()
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
				{"type": "memory"},
			},
		}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			types.DNSRootZone,
			"foo",
			types.DNSRecordTypeA,
			"0",
		).String(),
		tests.MustJSON(dns.Record{
			Data: "10.1.2.3",
		}),
	))
	role := dns.New(inst)
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, RoleConfig()))
	t.Cleanup(role.Stop)
	return role
}

func TestRoleDNS_EDNS0_ClientWithEDNS0(t *testing.T) {
	tests.Setup(t)
	role := setupMemoryZone(t)

	fw := NewNullDNSWriter()
	msg := &d.Msg{
		Question: []d.Question{{Name: "foo.", Qtype: d.TypeA, Qclass: d.ClassINET}},
	}
	msg.SetEdns0(1232, false)

	role.Handler(fw, msg)

	resp := fw.Msg()
	assert.NotNil(t, resp)
	opt := resp.IsEdns0()
	assert.NotNil(t, opt, "response should include EDNS0 when client advertises it")
	assert.Equal(t, uint16(1232), opt.UDPSize(), "response EDNS0 size should match client's advertised size")

	// Must have exactly one OPT record
	optCount := 0
	for _, rr := range resp.Extra {
		if rr.Header().Rrtype == d.TypeOPT {
			optCount++
		}
	}
	assert.Equal(t, 1, optCount, "response must contain exactly one OPT record")
}

func TestRoleDNS_EDNS0_ClientWithoutEDNS0(t *testing.T) {
	tests.Setup(t)
	role := setupMemoryZone(t)

	fw := NewNullDNSWriter()
	msg := &d.Msg{
		Question: []d.Question{{Name: "foo.", Qtype: d.TypeA, Qclass: d.ClassINET}},
	}
	// No EDNS0 in query

	role.Handler(fw, msg)

	resp := fw.Msg()
	assert.NotNil(t, resp)
	assert.Nil(t, resp.IsEdns0(), "response must not include EDNS0 when client did not advertise it")
}

func TestRoleDNS_ZoneFind(t *testing.T) {
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
	zone := role.FindZone("foo.bar.")
	assert.Equal(t, zone, role.FindZone("bar.baz."))
}

func TestRoleDNS_ZoneNoHandler(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			TestZone,
		).String(),
		tests.MustJSON(dns.Zone{
			HandlerConfigs: []map[string]interface{}{},
		}),
	))
	role := dns.New(inst)
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, RoleConfig()))
	defer role.Stop()

	msg := AssertDNS(t, role, []d.Question{
		{
			Name:   "foo.example.com.",
			Qtype:  d.TypeA,
			Qclass: d.ClassINET,
		},
	})
	assert.Equal(t, d.RcodeSuccess, msg.Rcode)
}

func TestRoleDNS_NoZone(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			TestZone,
		).String(),
		tests.MustJSON(dns.Zone{
			HandlerConfigs: []map[string]interface{}{},
		}),
	))
	role := dns.New(inst)
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, RoleConfig()))
	defer role.Stop()

	msg := AssertDNS(t, role, []d.Question{
		{
			Name:   TestZone2,
			Qtype:  d.TypeA,
			Qclass: d.ClassINET,
		},
	})
	assert.Equal(t, d.RcodeNameError, msg.Rcode)
}
