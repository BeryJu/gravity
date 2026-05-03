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

func TestRoleDNS_BlockyForwarder(t *testing.T) {
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
					"type":       "forward_blocky",
					"blocklists": "http://127.0.0.1:9005/blocky_file.txt",
					"to":         "127.0.0.1:1053",
				},
			},
		}),
	))
	role := dns.New(inst)
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, RoleConfig()))
	defer role.Stop()

	AssertDNS(t, role,
		[]d.Question{
			{
				Name:   "gravity.beryju.io.",
				Qtype:  d.TypeA,
				Qclass: d.ClassINET,
			},
		},
		"gravity.beryju.io.	0	IN	A	0.0.0.0",
	)
}

func TestRoleDNS_BlockyForwarder_Allow(t *testing.T) {
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
					"type":       "forward_blocky",
					"blocklists": "http://127.0.0.1:9005/blocky_file.txt",
					"allowlists": "gravity.beryju.io",
					"to":         "127.0.0.1:1053",
				},
			},
		}),
	))
	role := dns.New(inst)
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, RoleConfig()))
	defer role.Stop()

	AssertDNS(t, role,
		[]d.Question{
			{
				Name:   "gravity.beryju.io.",
				Qtype:  d.TypeA,
				Qclass: d.ClassINET,
			},
		},
		"gravity.beryju.io.	3600	IN	A	10.0.0.1",
	)
}
