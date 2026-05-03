package dhcp_test

import (
	"net"
	"runtime"
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dhcp"
	"beryju.io/gravity/pkg/roles/dhcp/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/ipv4"
)

func TestDHCP4_handle(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Test only supported on linux")
	}
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)
	role := dhcp.New(inst)

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyScopes,
			"test",
		).String(),
		tests.MustJSON(dhcp.Scope{
			SubnetCIDR: "10.100.0.0/24",
			Default:    true,
			TTL:        86400,
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "10.100.0.100",
				"range_end":   "10.100.0.250",
			},
		}),
	))
	tests.PanicIfError(role.Start(ctx, RoleConfig()))
	defer role.Stop()

	err := role.Handler4().Handle(DHCPRequestPayload, &ipv4.ControlMessage{
		IfIndex: 1,
	}, &net.UDPAddr{})
	if runtime.GOOS != "linux" {
		assert.Errorf(t, err, "sendEthernet not supported on current platform")
	} else {
		assert.NoError(t, err)
	}
}
