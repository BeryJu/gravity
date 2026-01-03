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
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestBindImport(t *testing.T) {
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
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", types.DNSRootRecord, "MX", "bind-import-0f7de7b0"),
					values: []interface{}{`{"data":"mail.example.net.","mxPreference":20}`},
				},
				{
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", types.DNSRootRecord, "MX", "bind-import-74506170"),
					values: []interface{}{`{"data":"mail.example.com.","mxPreference":10}`},
				},
				{
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", "ftp.", "CNAME", "bind-import-99f5b921"),
					values: []interface{}{`{"data":"ftp.example.net."}`},
				},
				{
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", "joe.", "A", "bind-import-74f397b8"),
					values: []interface{}{`{"data":"192.168.254.6"}`},
				},
				{
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", "mail.", "A", "bind-import-7b40edd6"),
					values: []interface{}{`{"data":"192.168.254.4"}`},
				},
				{
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", "ns1.", "A", "bind-import-0c61a839"),
					values: []interface{}{`{"data":"192.168.254.2"}`},
				},
				{
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", "www.", "A", "bind-import-b49788f4"),
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
				{
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", types.DNSRootRecord, "MX", "bind-import-41b8b28c"),
					values: []interface{}{`{"data":"mail2.example.com.","mxPreference":20}`},
				},
				{
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", types.DNSRootRecord, "MX", "bind-import-7f4cd2b5"),
					values: []interface{}{`{"data":"mail.example.com.","mxPreference":10}`},
				},
				{
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", "dns1.", "A", "bind-import-c59601e0"),
					values: []interface{}{`{"data":"10.0.1.1"}`},
				},
				{
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", "dns2.", "A", "bind-import-c3945d32"),
					values: []interface{}{`{"data":"10.0.1.2"}`},
				},
				{
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", "ftp.", "A", "bind-import-cfb6f7a3"),
					values: []interface{}{`{"data":"10.0.1.4"}`},
				},
				{
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", "ftp.", "A", "bind-import-ebe4d20c"),
					values: []interface{}{`{"data":"10.0.1.3"}`},
				},
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
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", "mail.", "A", "bind-import-7b40edd6"),
					values: []interface{}{`{"data":"192.168.254.4"}`},
				},
				{
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", "mail.", "AAAA", "bind-import-f5f89df9"),
					values: []interface{}{`{"data":"2001:db8:3333:4444:5555:6666:7777:8888"}`},
				},
				{
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", types.DNSRootRecord, "TXT", "bind-import-97c7d3ac"),
					values: []interface{}{`{"data":"foo"}`},
				},
				{
					key:    ri.KV().Key(types.KeyRole, types.KeyZones, "example.com.", "_sip._udp.", "SRV", "bind-import-4361fc86"),
					values: []interface{}{`{"data":"fs1.example.com.","srvPort":5060,"srvPriority":10,"srvWeight":100}`},
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

			x, err := os.Open(file.file)
			assert.NoError(t, err)

			c, err := bind.New(api, x)
			assert.NoError(t, err)
			assert.NoError(t, c.Run(ctx))
			for _, kv := range file.kv {
				tests.AssertEtcd(t, ri.KV(), kv.key, kv.values...)
			}
			assert.NoError(t, x.Close())
		})
	}
}
