package dns_test

import (
	"testing"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dns"
	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestRoleDNSHandlerEtcd(t *testing.T) {
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	inst.KV().Put(
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
	)
	inst.KV().Put(
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
	)

	role := dns.New(inst)
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, RoleConfig()))
	defer role.Stop()

	tests.WaitForPort(1054)
	assert.Equal(t, []string{"10.1.2.3"}, tests.DNSLookup("foo.", extconfig.Get().Listen(1054)))
}
