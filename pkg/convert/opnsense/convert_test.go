package opnsense_test

import (
	"os"
	"testing"

	"beryju.io/gravity/pkg/convert/opnsense"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dns"
	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/storage"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestOpnsenseImport(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	// Create DNS role to register API routes
	ri := rootInst.ForRole("dns", ctx)
	dns.New(ri)

	type kv struct {
		key    *storage.Key
		values []interface{}
	}

	cases := []struct {
		file string
		kv   []kv
	}{
		{
			file: "./fixtures/opnsense_export.xml",
			kv: []kv{
				{
					key: ri.KV().Key(types.KeyRole, types.KeyZones, "mgt.lan."),
					values: []interface{}{
						dns.Zone{DefaultTTL: 86400, Authoritative: true, HandlerConfigs: []map[string]interface{}{
							{
								"type": "etcd",
							},
						}},
					},
				},
			},
		},
	}

	api, stop := tests.APIClient(rootInst)
	defer stop()

	for _, file := range cases {
		t.Run(file.file, func(t *testing.T) {
			// Clean etcd before testing so we can debug easier
			_, err := ri.KV().Delete(ctx, ri.KV().Key(
				types.KeyRole,
				types.KeyZones,
			).Prefix(true).String(), clientv3.WithPrefix())
			assert.NoError(t, err)

			x, err := os.ReadFile(file.file)
			assert.NoError(t, err)

			c, err := opnsense.New(api, x)
			assert.NoError(t, err)
			assert.NoError(t, c.Run(ctx))
			for _, kv := range file.kv {
				tests.AssertEtcd(t, ri.KV(), kv.key, kv.values...)
			}
		})
	}
}
