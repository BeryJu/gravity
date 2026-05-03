package dns_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dns"
	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/tests"
	d "github.com/miekg/dns"
)

func BenchmarkRoleDNS_DefaultRootZone(b *testing.B) {
	tests.Setup(b)
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
				{
					"type": "etcd",
				},
				{
					"type":      "forward_ip",
					"to":        []string{"127.0.0.1:1053"},
					"cache_ttl": 3600,
				},
			},
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
					Name:   "gravity.beryju.io.",
					Qtype:  d.TypeA,
					Qclass: d.ClassINET,
				},
			},
		})
		_ = fw.Msg()
	}
}
