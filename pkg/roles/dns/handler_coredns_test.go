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

const CoreDNSConfig = `.:1342 {
    hosts {
        10.0.0.1 example.org
        fallthrough
    }
}`

func TestRoleDNSHandlerCoreDNS(t *testing.T) {
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
					"type":   "coredns",
					"config": CoreDNSConfig,
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
			Name:   "example.org.",
			Qtype:  d.TypeA,
			Qclass: d.ClassINET,
		},
	}, "example.org.	IN	A	10.0.0.1")
}
