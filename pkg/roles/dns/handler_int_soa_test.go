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

func TestRoleDNS_SOA(t *testing.T) {
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
			Authoritative: true,
			DefaultTTL:    3600,
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
			Name:   TestZone,
			Qtype:  d.TypeSOA,
			Qclass: d.ClassINET,
		},
	}, "example.com.	3600	IN	SOA	example.com. root.example.com. 1337 600 15 5 3600")
}

func TestRoleDNS_SOA_FromEtcd(t *testing.T) {
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
			Authoritative: true,
			DefaultTTL:    3600,
			HandlerConfigs: []map[string]interface{}{
				{
					"type": "forward_ip",
					"to":   "127.0.0.1:1053",
				},
			},
		}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			TestZone,
			"@",
			"SOA",
			"0",
		).String(),
		tests.MustJSON(dns.Record{
			Data:       "ns1.example.com.",
			SOAMbox:    "hostmaster.example.com.",
			SOASerial:  2024010101,
			SOARefresh: 3600,
			SOARetry:   900,
			SOAExpire:  604800,
			TTL:        3600,
		}),
	))
	role := dns.New(inst)
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, RoleConfig()))
	defer role.Stop()

	AssertDNS(t, role, []d.Question{
		{
			Name:   TestZone,
			Qtype:  d.TypeSOA,
			Qclass: d.ClassINET,
		},
	}, "example.com.	3600	IN	SOA	ns1.example.com. hostmaster.example.com. 2024010101 3600 900 604800 3600")
}
