package coredns_test

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

const CoreDNSConfig = `.:1342 {
    hosts {
        10.0.0.1 example.org
        fallthrough
    }
}`

func RoleConfig() []byte {
	return []byte(tests.MustJSON(dns.RoleConfig{
		Port: 1054,
	}))
}

func TestRoleDNSHandlerCoreDNS(t *testing.T) {
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

	fw := dns.NewNullDNSWriter()
	role.Handler(fw, &d.Msg{
		Question: []d.Question{
			{
				Name:   "example.org.",
				Qtype:  d.TypeA,
				Qclass: d.ClassINET,
			},
		},
	})
	ans := fw.Msg().Answer[0]
	assert.Equal(t, net.ParseIP("10.0.0.1").String(), ans.(*d.A).A.String())
}
