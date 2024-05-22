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
			".",
		).String(),
		tests.MustJSON(dns.Zone{
			HandlerConfigs: []map[string]string{
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
			".",
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
				Name:   "foo.",
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
			".",
		).String(),
		tests.MustJSON(dns.Zone{
			HandlerConfigs: []map[string]string{
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
			".",
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
				Name:   "foo.",
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
			"test.",
		).String(),
		tests.MustJSON(dns.Zone{
			HandlerConfigs: []map[string]string{
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

	fw := NewNullDNSWriter()
	role.Handler(fw, &d.Msg{
		Question: []d.Question{
			{
				Name:   "bar.test.",
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
				Name:   "foo.test.",
				Qtype:  d.TypeA,
				Qclass: d.ClassINET,
			},
		},
	})
	ans = fw.Msg().Answer[0]
	assert.Equal(t, "bar.test.", ans.(*d.CNAME).Target)
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
			".",
		).String(),
		tests.MustJSON(dns.Zone{
			HandlerConfigs: []map[string]string{
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
			".",
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
				Name:   "foo.bar.",
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
			"TesT.",
		).String(),
		tests.MustJSON(dns.Zone{
			HandlerConfigs: []map[string]string{
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

	fw := NewNullDNSWriter()
	role.Handler(fw, &d.Msg{
		Question: []d.Question{
			{
				Name:   "bar.test.",
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
			"test.",
		).String(),
		tests.MustJSON(dns.Zone{
			HandlerConfigs: []map[string]string{
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

	fw := NewNullDNSWriter()
	role.Handler(fw, &d.Msg{
		Question: []d.Question{
			{
				Name:   "bar.TesT.",
				Qtype:  d.TypeA,
				Qclass: d.ClassINET,
			},
		},
	})
	ans := fw.Msg().Answer[0]
	assert.Equal(t, net.ParseIP("10.1.2.3").String(), ans.(*d.A).A.String())
}
