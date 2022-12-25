package dns_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dns"
	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestRoleDNSZoneFind(t *testing.T) {
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
					"type": "forward_ip",
					"to":   "127.0.0.1:1053",
				},
			},
		}),
	)
	role := dns.New(inst)
	assert.NotNil(t, role)
	ctx := tests.Context()
	assert.Nil(t, role.Start(ctx, RoleConfig()))
	defer role.Stop()
	zone := role.FindZone("foo.bar.")
	assert.Equal(t, zone, role.FindZone("bar.baz."))
}
