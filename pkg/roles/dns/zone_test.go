package dns_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dns"
	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
)

func TestRoleDNSZoneFind(t *testing.T) {
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
		tests.MustPB(&types.Zone{
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
	zone := role.FindZone("foo.bar.")
	assert.Equal(t, zone, role.FindZone("bar.baz."))
}
