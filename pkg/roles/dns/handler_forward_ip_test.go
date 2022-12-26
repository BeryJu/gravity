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

func TestRoleDNSIPForwarder(t *testing.T) {
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
	time.Sleep(3 * time.Second)
	defer role.Stop()
	assert.Equal(t, []string{"10.0.0.1"}, tests.DNSLookup("gravity.beryju.io.", extconfig.Get().Listen(1054)))
}
