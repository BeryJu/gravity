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
