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

func TestRoleDNS_Memory(t *testing.T) {
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
					"type": "memory",
				},
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
	defer role.Stop()

	AssertDNS(t, role, []d.Question{
		{
			Name:   "foo.",
			Qtype:  d.TypeA,
			Qclass: d.ClassINET,
		},
	}, "foo.	0	IN	A	10.1.2.3")
}

func TestRoleDNS_Memory_Wildcard(t *testing.T) {
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
					"type": "memory",
				},
			},
		}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			types.DNSRootZone,
			"*",
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
	defer role.Stop()

	AssertDNS(t, role, []d.Question{
		{
			Name:   "foo.",
			Qtype:  d.TypeA,
			Qclass: d.ClassINET,
		},
	}, "foo.	0	IN	A	10.1.2.3")
}

func TestRoleDNS_Memory_CNAME(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			"test.",
		).String(),
		tests.MustJSON(dns.Zone{
			HandlerConfigs: []map[string]interface{}{
				{
					"type": "memory",
				},
			},
		}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			"test.",
			"foo",
			types.DNSRecordTypeCNAME,
			"0",
		).String(),
		tests.MustJSON(dns.Record{
			Data: "bar.test.",
		}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			"test.",
			"bar",
			types.DNSRecordTypeA,
			"0",
		).String(),
		tests.MustJSON(dns.Record{
			Data: "10.2.3.4",
		}),
	))

	role := dns.New(inst)
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, RoleConfig()))
	defer role.Stop()

	t.Run("A", func(t *testing.T) {
		AssertDNS(t, role, []d.Question{
			{
				Name:   "bar.test.",
				Qtype:  d.TypeA,
				Qclass: d.ClassINET,
			},
		}, "bar.test.	0	IN	A	10.2.3.4")
	})

	t.Run("CNAME", func(t *testing.T) {
		AssertDNS(t, role,
			[]d.Question{
				{
					Name:   "foo.test.",
					Qtype:  d.TypeCNAME,
					Qclass: d.ClassINET,
				},
			},
			"foo.test.	0	IN	CNAME	bar.test.",
		)
	})

	t.Run("Anything", func(t *testing.T) {
		AssertDNS(t, role,
			[]d.Question{
				{
					Name:   "foo.test.",
					Qtype:  d.TypeA,
					Qclass: d.ClassINET,
				},
			},
			"foo.test.	0	IN	CNAME	bar.test.",
			"bar.test.	0	IN	A	10.2.3.4",
		)
	})
}

func TestRoleDNS_Memory_WildcardNested(t *testing.T) {
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
					"type": "memory",
				},
			},
		}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			types.DNSRootZone,
			"*.*",
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
	defer role.Stop()

	AssertDNS(t, role, []d.Question{
		{
			Name:   "foo.bar.",
			Qtype:  d.TypeA,
			Qclass: d.ClassINET,
		},
	}, "foo.bar.	0	IN	A	10.1.2.3")
}

func TestRoleDNS_Memory_MixedCase(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			"TesT.",
		).String(),
		tests.MustJSON(dns.Zone{
			HandlerConfigs: []map[string]interface{}{
				{
					"type": "memory",
				},
			},
		}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			"TesT.",
			"bar",
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
	defer role.Stop()

	AssertDNS(t, role, []d.Question{
		{
			Name:   "bar.test.",
			Qtype:  d.TypeA,
			Qclass: d.ClassINET,
		},
	}, "bar.test.	0	IN	A	10.1.2.3")
}

func TestRoleDNS_Memory_MixedCase_Reverse(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			"test.",
		).String(),
		tests.MustJSON(dns.Zone{
			HandlerConfigs: []map[string]interface{}{
				{
					"type": "memory",
				},
			},
		}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			"test.",
			"bar",
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
	defer role.Stop()

	AssertDNS(t, role, []d.Question{
		{
			Name:   "bar.TesT.",
			Qtype:  d.TypeA,
			Qclass: d.ClassINET,
		},
	}, "bar.test.	0	IN	A	10.1.2.3")
}
