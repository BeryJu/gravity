package dhcp_test

import (
	"net"
	"testing"
	"time"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dhcp"
	"beryju.io/gravity/pkg/roles/dhcp/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/stretchr/testify/assert"
)

var DHCPDiscoverPayload = []byte{1, 1, 6, 0, 136, 9, 170, 251, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 178, 183, 134, 44, 211, 250, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 99, 130, 83, 99, 53, 1, 1, 55, 9, 1, 121, 3, 6, 15, 108, 114, 119, 252, 57, 2, 5, 220, 61, 7, 1, 178, 183, 134, 44, 211, 250, 51, 4, 0, 118, 167, 0, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}

func TestDHCPDiscover(t *testing.T) {
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)
	role := dhcp.New(inst)
	Cleanup()

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
	tests.PanicIfError(role.Start(ctx, []byte{}))
	defer role.Stop()

	req, err := dhcpv4.FromBytes(DHCPDiscoverPayload)
	assert.NoError(t, err)
	req4 := role.NewRequest4(req)
	res := role.Handler4(req4)
	assert.NotNil(t, res)
	assert.Equal(t, "10.100.0.100", res.YourIPAddr.String())
	ones, bits := res.SubnetMask().Size()
	assert.Equal(t, 24, ones)
	assert.Equal(t, 32, bits)
	assert.Equal(t, "b2:b7:86:2c:d3:fa", res.ClientHWAddr.String())
	assert.Equal(t, 86400*time.Second, res.IPAddressLeaseTime(1*time.Second))
}

func TestDHCPDiscoverGateway(t *testing.T) {
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)
	role := dhcp.New(inst)
	Cleanup()

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyScopes,
			"test",
		).String(),
		tests.MustJSON(dhcp.Scope{
			SubnetCIDR: "10.100.0.0/24",
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
			SubnetCIDR: "10.100.32.0/24",
			TTL:        86400,
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "10.100.32.100",
				"range_end":   "10.100.32.250",
			},
		}),
	))
	tests.PanicIfError(role.Start(ctx, []byte{}))
	defer role.Stop()

	req, err := dhcpv4.FromBytes(DHCPDiscoverPayload)
	req.GatewayIPAddr = net.ParseIP("10.100.32.1")
	assert.NoError(t, err)
	req4 := role.NewRequest4(req)
	res := role.Handler4(req4)
	assert.NotNil(t, res)
	assert.Equal(t, "10.100.32.100", res.YourIPAddr.String())
	ones, bits := res.SubnetMask().Size()
	assert.Equal(t, 24, ones)
	assert.Equal(t, 32, bits)
	assert.Equal(t, "b2:b7:86:2c:d3:fa", res.ClientHWAddr.String())
	assert.Equal(t, 86400*time.Second, res.IPAddressLeaseTime(1*time.Second))
}

func TestDHCPDiscoverDNS(t *testing.T) {
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)
	role := dhcp.New(inst)
	Cleanup()

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
			DNS: &dhcp.ScopeDNS{
				Zone:              "test.gravity.beryju.io",
				AddZoneInHostname: true,
			},
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "10.100.0.100",
				"range_end":   "10.100.0.250",
			},
		}),
	))

	tests.PanicIfError(role.Start(ctx, []byte{}))
	defer role.Stop()

	req, err := dhcpv4.FromBytes(DHCPDiscoverPayload)
	assert.NoError(t, err)
	req4 := role.NewRequest4(req)
	res := role.Handler4(req4)
	assert.NotNil(t, res)
	assert.Equal(t, "10.100.0.100", res.YourIPAddr.String())
	ones, bits := res.SubnetMask().Size()
	assert.Equal(t, 24, ones)
	assert.Equal(t, 32, bits)
	assert.Equal(t, "b2:b7:86:2c:d3:fa", res.ClientHWAddr.String())
	assert.Equal(t, 86400*time.Second, res.IPAddressLeaseTime(1*time.Second))
	assert.Equal(t, "test.gravity.beryju.io", res.DomainName())
}
