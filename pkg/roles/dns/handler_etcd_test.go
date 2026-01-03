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

const TestZone = "example.com."
const TestZone2 = "example.net."

func TestRoleDNS_Etcd(t *testing.T) {
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
			HandlerConfigs: []map[string]interface{}{
				{
					"type": "etcd",
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
			Name:   "foo.example.com.",
			Qtype:  d.TypeA,
			Qclass: d.ClassINET,
		},
	}, "foo.example.com.	0	IN	A	10.1.2.3")
}

// Test DNS Entry at root of zone
func TestRoleDNS_Etcd_Root(t *testing.T) {
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
			HandlerConfigs: []map[string]interface{}{
				{
					"type": "etcd",
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
			Name:   "example.com.",
			Qtype:  d.TypeA,
			Qclass: d.ClassINET,
		},
	}, "example.com.	0	IN	A	10.1.2.3")
}

func TestRoleDNS_Etcd_Wildcard(t *testing.T) {
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
			HandlerConfigs: []map[string]interface{}{
				{
					"type": "etcd",
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
			Name:   "foo.example.com.",
			Qtype:  d.TypeA,
			Qclass: d.ClassINET,
		},
	}, "foo.example.com.	0	IN	A	10.1.2.3")
}

func TestRoleDNS_Etcd_CNAME(t *testing.T) {
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
			HandlerConfigs: []map[string]interface{}{
				{
					"type": "etcd",
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
			"svc",
			types.DNSRecordTypeCNAME,
			"0",
		).String(),
		tests.MustJSON(dns.Record{
			Data: "host.example.com.",
		}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			TestZone,
			"host",
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
		// Test lookup of direct host
		AssertDNS(t, role, []d.Question{
			{
				Name:   "host.example.com.",
				Qtype:  d.TypeA,
				Qclass: d.ClassINET,
			},
		}, "host.example.com.	0	IN	A	10.2.3.4")
	})

	t.Run("CNAME", func(t *testing.T) {
		// Explicitly ask for CNAME
		AssertDNS(t, role, []d.Question{
			{
				Name:   "svc.example.com.",
				Qtype:  d.TypeCNAME,
				Qclass: d.ClassINET,
			},
		}, "svc.example.com.	0	IN	CNAME	host.example.com.")
	})

	t.Run("Anything", func(t *testing.T) {
		// Ask for anything
		AssertDNS(t, role,
			[]d.Question{
				{
					Name:   "svc.example.com.",
					Qtype:  d.TypeA,
					Qclass: d.ClassINET,
				},
			},
			"svc.example.com.	0	IN	CNAME	host.example.com.",
			"host.example.com.	0	IN	A	10.2.3.4",
		)
	})
}

// Test DNS CNAME that points to another zone
func TestRoleDNS_Etcd_CNAME_MultiZone(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	// Create both zones and both records
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			TestZone,
		).String(),
		tests.MustJSON(dns.Zone{
			HandlerConfigs: []map[string]interface{}{
				{
					"type": "etcd",
				},
			},
		}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			TestZone2,
		).String(),
		tests.MustJSON(dns.Zone{
			HandlerConfigs: []map[string]interface{}{
				{
					"type": "etcd",
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
			"foo",
			types.DNSRecordTypeCNAME,
			"0",
		).String(),
		tests.MustJSON(dns.Record{
			Data: "bar.example.net.",
		}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			TestZone2,
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

	// Ask for anything
	AssertDNS(t, role,
		[]d.Question{
			{
				Name:   "foo.example.com.",
				Qtype:  d.TypeA,
				Qclass: d.ClassINET,
			},
		},
		"foo.example.com.	0	IN	CNAME	bar.example.net.",
		"bar.example.net.	0	IN	A	10.2.3.4",
	)
}

// Test DNS CNAME that points to another zone, recursive
func TestRoleDNS_Etcd_CNAME_Recursive(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	// Create both zones and both records
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			TestZone,
		).String(),
		tests.MustJSON(dns.Zone{
			HandlerConfigs: []map[string]interface{}{
				{
					"type": "etcd",
				},
			},
		}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			TestZone2,
		).String(),
		tests.MustJSON(dns.Zone{
			HandlerConfigs: []map[string]interface{}{
				{
					"type": "etcd",
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
			"foo",
			types.DNSRecordTypeCNAME,
			"0",
		).String(),
		tests.MustJSON(dns.Record{
			Data: "bar.example.net.",
		}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			TestZone2,
			"bar",
			types.DNSRecordTypeCNAME,
			"0",
		).String(),
		tests.MustJSON(dns.Record{
			Data: "foo.example.com",
		}),
	))

	role := dns.New(inst)
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, RoleConfig()))
	defer role.Stop()

	AssertDNS(t, role,
		[]d.Question{
			{
				Name:   "foo.example.com.",
				Qtype:  d.TypeA,
				Qclass: d.ClassINET,
			},
		},
		"foo.example.com.	0	IN	CNAME	bar.example.net.",
		"bar.example.net.	0	IN	CNAME	foo.example.com.",
		"bar.example.net.	0	IN	CNAME	foo.example.com.",
		"foo.example.com.	0	IN	CNAME	bar.example.net.",
		"foo.example.com.	0	IN	CNAME	bar.example.net.",
		"bar.example.net.	0	IN	CNAME	foo.example.com.",
		"bar.example.net.	0	IN	CNAME	foo.example.com.",
		"foo.example.com.	0	IN	CNAME	bar.example.net.",
		"foo.example.com.	0	IN	CNAME	bar.example.net.",
		"bar.example.net.	0	IN	CNAME	foo.example.com.",
		"bar.example.net.	0	IN	CNAME	foo.example.com.",
		"foo.example.com.	0	IN	CNAME	bar.example.net.",
		"foo.example.com.	0	IN	CNAME	bar.example.net.",
		"bar.example.net.	0	IN	CNAME	foo.example.com.",
		"bar.example.net.	0	IN	CNAME	foo.example.com.",
		"foo.example.com.	0	IN	CNAME	bar.example.net.",
		"foo.example.com.	0	IN	CNAME	bar.example.net.",
		"bar.example.net.	0	IN	CNAME	foo.example.com.",
		"bar.example.net.	0	IN	CNAME	foo.example.com.",
		"foo.example.com.	0	IN	CNAME	bar.example.net.",
		"foo.example.com.	0	IN	CNAME	bar.example.net.",
	)
}

func TestRoleDNS_Etcd_WildcardNested(t *testing.T) {
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
			HandlerConfigs: []map[string]interface{}{
				{
					"type": "etcd",
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

	AssertDNS(t, role,
		[]d.Question{
			{
				Name:   "foo.bar.example.com.",
				Qtype:  d.TypeA,
				Qclass: d.ClassINET,
			},
		},
		"foo.bar.example.com.	0	IN	A	10.1.2.3",
	)
}

func TestRoleDNS_Etcd_MixedCase(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			"eXaMpLe.CoM.",
		).String(),
		tests.MustJSON(dns.Zone{
			HandlerConfigs: []map[string]interface{}{
				{
					"type": "etcd",
				},
			},
		}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			"eXaMpLe.CoM.",
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

	AssertDNS(t, role,
		[]d.Question{
			{
				Name:   "bar.example.com.",
				Qtype:  d.TypeA,
				Qclass: d.ClassINET,
			},
		},
		"bar.example.com.	0	IN	A	10.1.2.3",
	)
}

func TestRoleDNS_Etcd_MixedCase_Reverse(t *testing.T) {
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
			HandlerConfigs: []map[string]interface{}{
				{
					"type": "etcd",
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

	AssertDNS(t, role,
		[]d.Question{
			{
				Name:   "bar.eXaMpLe.CoM.",
				Qtype:  d.TypeA,
				Qclass: d.ClassINET,
			},
		},
		"bar.example.com.	0	IN	A	10.1.2.3",
	)
}
