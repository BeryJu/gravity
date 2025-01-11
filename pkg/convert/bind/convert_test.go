package bind_test

import (
	"os"
	"testing"

	"beryju.io/gravity/pkg/convert/bind"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dns"
	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/storage"
	"beryju.io/gravity/pkg/tests"
	"beryju.io/gravity/pkg/tests/api"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestBindImport(t *testing.T) {
	defer tests.Setup(t)()
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
			file: "./fixtures/example.zone",
			kv: []kv{
				{
					key: ri.KV().Key(types.KeyRole, types.KeyZones, "example.com."),
					values: []interface{}{
						dns.Zone{DefaultTTL: 1814400, Authoritative: true, HandlerConfigs: []map[string]interface{}{
							{
								"type": "etcd",
							},
						}},
					},
				},
				{
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", types.DNSRootRecord, "MX", "bind-import"),
					values: []interface{}{`{"data":"mail.example.net.","mxPreference":20}`},
				},
				{
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", "ftp.", "CNAME", "bind-import"),
					values: []interface{}{`{"data":"ftp.example.net."}`},
				},
				{
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", "joe.", "A", "bind-import"),
					values: []interface{}{`{"data":"192.168.254.6"}`},
				},
				{
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", "mail.", "A", "bind-import"),
					values: []interface{}{`{"data":"192.168.254.4"}`},
				},
				{
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", "ns1.", "A", "bind-import"),
					values: []interface{}{`{"data":"192.168.254.2"}`},
				},
				{
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", "www.", "A", "bind-import"),
					values: []interface{}{`{"data":"192.168.254.7"}`},
				},
			},
		},
		{
			file: "./fixtures/example2.zone",
			kv: []kv{
				{
					key: ri.KV().Key(types.KeyRole, types.KeyZones, "example.com."),
					values: []interface{}{
						dns.Zone{DefaultTTL: 604800, Authoritative: true, HandlerConfigs: []map[string]interface{}{
							{
								"type": "etcd",
							},
						}},
					},
				},
				// {
				// 	key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", types.DNSRootRecord, "MX", "bind-import"),
				// 	values: []interface{}{`{"data":"mail.example.net.","mxPreference":20}`},
				// },
				// {
				// 	key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", "ftp.", "CNAME", "bind-import"),
				// 	values: []interface{}{`{"data":"ftp.example.net."}`},
				// },
				// {
				// 	key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", "joe.", "A", "bind-import"),
				// 	values: []interface{}{`{"data":"192.168.254.6"}`},
				// },
				// {
				// 	key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", "mail.", "A", "bind-import"),
				// 	values: []interface{}{`{"data":"192.168.254.4"}`},
				// },
				// {
				// 	key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", "ns1.", "A", "bind-import"),
				// 	values: []interface{}{`{"data":"192.168.254.2"}`},
				// },
				// {
				// 	key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", "www.", "A", "bind-import"),
				// 	values: []interface{}{`{"data":"192.168.254.7"}`},
				// },
			},
		},
		{
			file: "./fixtures/example-all.zone",
			kv: []kv{
				{
					key: ri.KV().Key(types.KeyRole, types.KeyZones, "example.com."),
					values: []interface{}{
						dns.Zone{DefaultTTL: 1814400, Authoritative: true, HandlerConfigs: []map[string]interface{}{
							{
								"type": "etcd",
							},
						}},
					},
				},
				{
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", "mail.", "A", "bind-import"),
					values: []interface{}{`{"data":"192.168.254.4"}`},
				},
				{
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", "mail.", "AAAA", "bind-import"),
					values: []interface{}{`{"data":"2001:db8:3333:4444:5555:6666:7777:8888"}`},
				},
				{
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", types.DNSRootRecord, "TXT", "bind-import"),
					values: []interface{}{`{"data":"foo"}`},
				},
				{
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", "_sip._udp.", "SRV", "bind-import"),
					values: []interface{}{`{"data":"fs1.example.com.","srvPort":5060,"srvPriority":10,"srvWeight":100}`},
				},
			},
		},
	}

	api, stop := api.APIClient(rootInst)
	defer stop()

	for _, file := range cases {
		t.Run(file.file, func(t *testing.T) {
			ri.KV().Delete(ctx, ri.KV().Key(
				types.KeyRole,
				types.KeyZones,
			).Prefix(true).String(), clientv3.WithPrefix())
			x, err := os.Open(file.file)
			assert.NoError(t, err)
			defer x.Close()
			c, err := bind.New(api, x)
			assert.NoError(t, err)
			assert.NoError(t, c.Run(ctx))
			for _, kv := range file.kv {
				tests.AssertEtcd(t, ri.KV(), kv.key, kv.values...)
			}
		})
	}
}
