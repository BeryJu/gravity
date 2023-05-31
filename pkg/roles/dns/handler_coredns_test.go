package dns_test

import (
	"testing"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dns"
	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
	"google.golang.org/protobuf/types/known/structpb"
)

const CoreDNSConfig = `.:1342 {
    hosts {
        10.0.0.1 example.org
        fallthrough
    }
}`

func TestRoleDNSHandlerCoreDNS(t *testing.T) {
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
						"type":   structpb.NewStringValue("coredns"),
						"config": structpb.NewStringValue(CoreDNSConfig),
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
	assert.Equal(t, []string{"10.0.0.1"}, tests.DNSLookup("example.org.", extconfig.Get().Listen(1054)))
}
