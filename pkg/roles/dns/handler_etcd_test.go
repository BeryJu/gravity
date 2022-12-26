package dns_test

import (
	"testing"
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dns"
	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestRoleDNSHandlerEtcd(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("dns")
	inst.KV().Put(
		tests.Context(),
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
		tests.Context(),
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
	ctx := tests.Context()
	assert.Nil(t, role.Start(ctx, RoleConfig()))
	time.Sleep(3 * time.Second)
	defer role.Stop()

	assert.Equal(t, []string{"10.1.2.3"}, tests.DNSLookup("foo.", extconfig.Get().Listen(53)))
}
