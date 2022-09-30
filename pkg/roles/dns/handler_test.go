package dns_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dns"
	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func CreateSomeRecords(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("dns")
	zone := "gravity.beryju.io."
	inst.KV().Put(
		tests.Context(),
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			zone,
		).String(),
		tests.MustJSON(dns.Zone{
			Authoritative: true,
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
			zone,
			"foo",
			"A",
		).String(),
		tests.MustJSON(dns.Record{
			// See https://en.wikipedia.org/wiki/Reserved_IP_addresses
			Data: "192.0.2.1",
		}),
	)
}

func Role() (*instance.Instance, *dns.Role) {
	rootInst := instance.New()
	inst := rootInst.ForRole("dns")
	role := dns.New(inst)
	return rootInst, role
}

func TestDNS_FindZone(t *testing.T) {
	CreateSomeRecords(t)
	_, role := Role()
	z := role.FindZone("foo.gravity.beryju.io.")
	assert.NotNil(t, z)
}
