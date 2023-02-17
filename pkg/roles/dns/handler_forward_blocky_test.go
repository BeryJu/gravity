package dns_test

import (
	"testing"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dns"
	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestRoleDNS_BlockyForwarder(t *testing.T) {
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	inst.KV().Delete(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
		).Prefix(true).String(),
		clientv3.WithPrefix(),
	)
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
					"type":       "forward_blocky",
					"blocklists": "http://127.0.0.1:9005/blocky_file.txt",
					"to":         "127.0.0.1:1053",
				},
			},
		}),
	)
	role := dns.New(inst)
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, RoleConfig()))
	defer role.Stop()

	tests.WaitForPort(extconfig.Get().Listen(1054))
	assert.Equal(t, []string{"0.0.0.0", "::"}, tests.DNSLookup("gravity.beryju.io.", extconfig.Get().Listen(1054)))
}
