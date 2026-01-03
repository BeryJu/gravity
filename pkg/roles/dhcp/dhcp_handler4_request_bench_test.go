package dhcp_test

import (
	"net"
	"net/netip"
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dhcp"
	"beryju.io/gravity/pkg/roles/dhcp/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/insomniacslk/dhcp/dhcpv4"
)

func BenchmarkRoleDHCP_Request(b *testing.B) {
	tests.Setup(b)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)
	role := dhcp.New(inst)

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyScopes,
			"test1",
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
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyScopes,
			"test2",
		).String(),
		tests.MustJSON(dhcp.Scope{
			SubnetCIDR: "10.100.1.0/24",
			Default:    true,
			TTL:        86400,
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "10.100.1.100",
				"range_end":   "10.100.1.250",
			},
		}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyScopes,
			"test3",
		).String(),
		tests.MustJSON(dhcp.Scope{
			SubnetCIDR: "10.100.2.0/24",
			Default:    true,
			TTL:        86400,
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "10.100.2.100",
				"range_end":   "10.100.2.250",
			},
		}),
	))

	tests.PanicIfError(role.Start(ctx, RoleConfig()))
	defer role.Stop()
	c, err := netip.ParsePrefix("10.100.0.0/24")
	if err != nil {
		panic(err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		rr := &dhcpv4.DHCPv4{
			GatewayIPAddr: net.ParseIP(c.Addr().Next().String()),
			ClientHWAddr:  tests.RandomMAC(),
			OpCode:        dhcpv4.OpcodeBootRequest,
		}
		rr.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeRequest))
		req4 := role.NewRequest4(rr)
		_ = role.Handle4(req4)
	}
}
