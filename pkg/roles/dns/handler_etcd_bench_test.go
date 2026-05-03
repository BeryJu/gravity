package dns_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dns"
	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/tests"
	d "github.com/miekg/dns"
)

func BenchmarkRoleDNS_Etcd(b *testing.B) {
	tests.Setup(b)
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
	_ = role.Start(ctx, RoleConfig())
	defer role.Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
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
		_ = fw.Msg()
	}
}
