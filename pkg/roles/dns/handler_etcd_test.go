package dns_test

import (
	"net"
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
	defer tests.Setup(t)()
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

	fw := NewNullDNSWriter()
	role.Handler(fw, &d.Msg{
		Question: []d.Question{
			{
				Name:   "foo.example.com.",
				Qtype:  d.TypeA,
				Qclass: d.ClassINET,
			},
		},
	})
	ans := fw.Msg().Answer[0]
	assert.Equal(t, net.ParseIP("10.1.2.3").String(), ans.(*d.A).A.String())
}

// Test DNS Entry at root of zone
func TestRoleDNS_Etcd_Root(t *testing.T) {
	defer tests.Setup(t)()
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

	fw := NewNullDNSWriter()
	role.Handler(fw, &d.Msg{
		Question: []d.Question{
			{
				Name:   TestZone,
				Qtype:  d.TypeA,
				Qclass: d.ClassINET,
			},
		},
	})
	ans := fw.Msg().Answer[0]
	assert.Equal(t, net.ParseIP("10.1.2.3").String(), ans.(*d.A).A.String())
}

func TestRoleDNS_Etcd_Wildcard(t *testing.T) {
	defer tests.Setup(t)()
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

	fw := NewNullDNSWriter()
	role.Handler(fw, &d.Msg{
		Question: []d.Question{
			{
				Name:   "foo.example.com.",
				Qtype:  d.TypeA,
				Qclass: d.ClassINET,
			},
		},
	})
	ans := fw.Msg().Answer[0]
	assert.Equal(t, net.ParseIP("10.1.2.3").String(), ans.(*d.A).A.String())
}

func TestRoleDNS_Etcd_CNAME(t *testing.T) {
	defer tests.Setup(t)()
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
			types.DNSRecordTypeCNAME,
			"0",
		).String(),
		tests.MustJSON(dns.Record{
			Data: "bar.example.com.",
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
			Data: "10.2.3.4",
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
				Name:   "bar.example.com.",
				Qtype:  d.TypeA,
				Qclass: d.ClassINET,
			},
		},
	})
	ans := fw.Msg().Answer[0]
	assert.Equal(t, net.ParseIP("10.2.3.4").String(), ans.(*d.A).A.String())

	fw = NewNullDNSWriter()
	role.Handler(fw, &d.Msg{
		Question: []d.Question{
			{
				Name:   "foo.example.com.",
				Qtype:  d.TypeCNAME,
				Qclass: d.ClassINET,
			},
		},
	})
	ans = fw.Msg().Answer[0]
	assert.Equal(t, "bar.example.com.", ans.(*d.CNAME).Target)

	fw = NewNullDNSWriter()
	role.Handler(fw, &d.Msg{
		Question: []d.Question{
			{
				Name:   "foo.example.com.",
				Qclass: d.ClassINET,
			},
		},
	})
	assert.Len(t, fw.Msg().Answer, 2)
	assert.Equal(t, "bar.example.com.", fw.Msg().Answer[0].(*d.CNAME).Target)
	assert.Equal(t, net.ParseIP("10.2.3.4").String(), fw.Msg().Answer[1].(*d.A).A.String())
}

// Test DNS CNAME that points to another zone
func TestRoleDNS_Etcd_CNAME_MultiZone(t *testing.T) {
	defer tests.Setup(t)()
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

	fw := dns.NewNullDNSWriter()
	role.Handler(fw, &d.Msg{
		Question: []d.Question{
			{
				Name:   "foo.example.com.",
				Qtype:  d.TypeCNAME,
				Qclass: d.ClassINET,
			},
		},
	})
	assert.Len(t, fw.Msg().Answer, 2)
	assert.Equal(t, "bar.example.net.", fw.Msg().Answer[0].(*d.CNAME).Target)
	assert.Equal(t, net.ParseIP("10.2.3.4").String(), fw.Msg().Answer[1].(*d.A).A.String())
}

// Test DNS CNAME that points to another zone, recursive
func TestRoleDNS_Etcd_CNAME_Recursive(t *testing.T) {
	defer tests.Setup(t)()
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

	fw := dns.NewNullDNSWriter()
	role.Handler(fw, &d.Msg{
		Question: []d.Question{
			{
				Name:   "foo.example.com.",
				Qtype:  d.TypeCNAME,
				Qclass: d.ClassINET,
			},
		},
	})
	assert.Len(t, fw.Msg().Answer, 2)
}

func TestRoleDNS_Etcd_WildcardNested(t *testing.T) {
	defer tests.Setup(t)()
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

	fw := NewNullDNSWriter()
	role.Handler(fw, &d.Msg{
		Question: []d.Question{
			{
				Name:   "foo.bar.example.com.",
				Qtype:  d.TypeA,
				Qclass: d.ClassINET,
			},
		},
	})
	ans := fw.Msg().Answer[0]
	assert.Equal(t, net.ParseIP("10.1.2.3").String(), ans.(*d.A).A.String())
}

func TestRoleDNS_Etcd_MixedCase(t *testing.T) {
	defer tests.Setup(t)()
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

	fw := NewNullDNSWriter()
	role.Handler(fw, &d.Msg{
		Question: []d.Question{
			{
				Name:   "bar.example.com.",
				Qtype:  d.TypeA,
				Qclass: d.ClassINET,
			},
		},
	})
	ans := fw.Msg().Answer[0]
	assert.Equal(t, net.ParseIP("10.1.2.3").String(), ans.(*d.A).A.String())
}

func TestRoleDNS_Etcd_MixedCase_Reverse(t *testing.T) {
	defer tests.Setup(t)()
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

	fw := NewNullDNSWriter()
	role.Handler(fw, &d.Msg{
		Question: []d.Question{
			{
				Name:   "bar.eXaMpLe.CoM.",
				Qtype:  d.TypeA,
				Qclass: d.ClassINET,
			},
		},
	})
	ans := fw.Msg().Answer[0]
	assert.Equal(t, net.ParseIP("10.1.2.3").String(), ans.(*d.A).A.String())
}
