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
	clientv3 "go.etcd.io/etcd/client/v3"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestRoleDNS_IPForwarder_v4(t *testing.T) {
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
		tests.MustJSON(types.Zone{
			HandlerConfigs: []*structpb.Struct{
				{
					Fields: map[string]*structpb.Value{
						"type": structpb.NewStringValue("forward_ip"),
						"to":   structpb.NewStringValue("127.0.0.1:1053"),
					},
				},
			},
		}),
	)
	role := dns.New(inst)
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, RoleConfig()))
	defer role.Stop()

	tests.WaitForPort(1054)
	assert.Equal(t, []string{"10.0.0.1"}, tests.DNSLookup("gravity.beryju.io.", extconfig.Get().Listen(1054)))
}

func TestRoleDNS_IPForwarder_v4_Cache(t *testing.T) {
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dns", ctx)
	inst.KV().Delete(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			".",
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
		tests.MustJSON(types.Zone{
			HandlerConfigs: []*structpb.Struct{
				{
					Fields: map[string]*structpb.Value{
						"type":      structpb.NewStringValue("forward_ip"),
						"to":        structpb.NewStringValue("127.0.0.1:1053"),
						"cache_ttl": structpb.NewStringValue("-2"),
					},
				},
			},
		}),
	)
	role := dns.New(inst)
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, RoleConfig()))
	defer role.Stop()
	tests.WaitForPort(1054)
	assert.Equal(t, []string{"10.0.0.1"}, tests.DNSLookup("gravity.beryju.io.", extconfig.Get().Listen(1054)))
	time.Sleep(3 * time.Second)
	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			".",
			"gravity.beryju.io",
			types.DNSRecordTypeA,
			"0",
		),
		types.Record{
			Data: "10.0.0.1",
		},
	)
}

func TestRoleDNS_IPForwarder_v6(t *testing.T) {
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
		tests.MustJSON(types.Zone{
			HandlerConfigs: []*structpb.Struct{
				{
					Fields: map[string]*structpb.Value{
						"type": structpb.NewStringValue("forward_ip"),
						"to":   structpb.NewStringValue("127.0.0.1:1053"),
					},
				},
			},
		}),
	)
	role := dns.New(inst)
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, RoleConfig()))
	defer role.Stop()

	tests.WaitForPort(1054)
	assert.Equal(t, []string{"fe80::1"}, tests.DNSLookup("ipv6.t.gravity.beryju.io.", extconfig.Get().Listen(1054)))
}
