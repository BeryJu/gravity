package dhcp_test

import (
	"net"
	"testing"
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dhcp"
	"beryju.io/gravity/pkg/roles/dhcp/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/stretchr/testify/assert"
)

var (
	DHCPDiscoverPayload        = []byte{1, 1, 6, 0, 136, 9, 170, 251, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 178, 183, 134, 44, 211, 250, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 99, 130, 83, 99, 53, 1, 1, 55, 9, 1, 121, 3, 6, 15, 108, 114, 119, 252, 57, 2, 5, 220, 61, 7, 1, 178, 183, 134, 44, 211, 250, 51, 4, 0, 118, 167, 0, 255, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0}
	DHCPDiscoverGatewayPayload = []byte{1, 1, 6, 1, 233, 21, 60, 223, 0, 0, 128, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 10, 82, 2, 1, 0, 80, 86, 155, 45, 115, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 99, 130, 83, 99, 53, 1, 1, 61, 7, 1, 0, 80, 86, 155, 45, 115, 50, 4, 10, 82, 7, 200, 12, 15, 68, 69, 83, 75, 84, 79, 80, 45, 49, 67, 67, 53, 50, 79, 75, 60, 8, 77, 83, 70, 84, 32, 53, 46, 48, 55, 14, 1, 3, 6, 15, 31, 33, 43, 44, 46, 47, 119, 121, 249, 252, 255}
)

func TestDHCPDiscover(t *testing.T) {
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

	req, err := dhcpv4.FromBytes(DHCPDiscoverPayload)
	assert.NoError(t, err)
	req4 := role.NewRequest4(req)
	res := role.Handle4(req4)
	assert.NotNil(t, res)
	assert.Equal(t, "10.100.0.100", res.YourIPAddr.String())
	ones, bits := res.SubnetMask().Size()
	assert.Equal(t, 24, ones)
	assert.Equal(t, 32, bits)
	assert.Equal(t, "b2:b7:86:2c:d3:fa", res.ClientHWAddr.String())
	assert.Equal(t, 86400*time.Second, res.IPAddressLeaseTime(1*time.Second))
}

func TestDHCPDiscover_Gateway(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)
	role := dhcp.New(inst)

	extconfig.Get().Instance.IP = "10.82.7.103"

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyScopes,
			"test",
		).String(),
		tests.MustJSON(dhcp.Scope{
			SubnetCIDR: "10.82.7.0/24",
			TTL:        86400,
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "10.82.7.100",
				"range_end":   "10.82.7.250",
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
			SubnetCIDR: "10.82.2.0/24",
			TTL:        86400,
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "10.82.2.100",
				"range_end":   "10.82.2.250",
			},
		}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyScopes,
			"test_larger",
		).String(),
		tests.MustJSON(dhcp.Scope{
			SubnetCIDR: "10.82.0.0/16",
			TTL:        86400,
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "10.82.0.100",
				"range_end":   "10.82.0.250",
			},
		}),
	))
	tests.PanicIfError(role.Start(ctx, RoleConfig()))
	defer role.Stop()

	req, err := dhcpv4.FromBytes(DHCPDiscoverGatewayPayload)
	assert.NoError(t, err)
	req4 := role.NewRequest4(req)
	res := role.Handle4(req4)
	assert.NotNil(t, res)
	assert.Equal(t, "10.82.2.100", res.YourIPAddr.String())
	ones, bits := res.SubnetMask().Size()
	assert.Equal(t, 24, ones)
	assert.Equal(t, 32, bits)
	assert.Equal(t, "00:50:56:9b:2d:73", res.ClientHWAddr.String())
	assert.Equal(t, 86400*time.Second, res.IPAddressLeaseTime(1*time.Second))
}

func TestDHCPDiscover_RequestedIP_WrongSubnet(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)
	role := dhcp.New(inst)

	extconfig.Get().Instance.IP = "10.82.7.103"

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyScopes,
			"test",
		).String(),
		tests.MustJSON(dhcp.Scope{
			SubnetCIDR: "10.82.7.0/24",
			TTL:        86400,
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "10.82.7.100",
				"range_end":   "10.82.7.250",
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
			SubnetCIDR: "10.82.2.0/24",
			TTL:        86400,
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "10.82.2.100",
				"range_end":   "10.82.2.250",
			},
		}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyScopes,
			"test_larger",
		).String(),
		tests.MustJSON(dhcp.Scope{
			SubnetCIDR: "10.82.0.0/16",
			TTL:        86400,
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "10.82.0.100",
				"range_end":   "10.82.0.250",
			},
		}),
	))
	tests.PanicIfError(role.Start(ctx, RoleConfig()))
	defer role.Stop()

	req, err := dhcpv4.FromBytes(DHCPDiscoverGatewayPayload)
	req.UpdateOption(dhcpv4.OptRequestedIPAddress(net.ParseIP("10.82.7.200")))
	assert.NoError(t, err)
	req4 := role.NewRequest4(req)
	res := role.Handle4(req4)
	assert.NotNil(t, res)
	assert.Equal(t, "10.82.2.100", res.YourIPAddr.String())
	ones, bits := res.SubnetMask().Size()
	assert.Equal(t, 24, ones)
	assert.Equal(t, 32, bits)
	assert.Equal(t, "00:50:56:9b:2d:73", res.ClientHWAddr.String())
	assert.Equal(t, 86400*time.Second, res.IPAddressLeaseTime(1*time.Second))
}

func TestDHCPDiscoverDNS(t *testing.T) {
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

	tests.PanicIfError(role.Start(ctx, RoleConfig()))
	defer role.Stop()

	req, err := dhcpv4.FromBytes(DHCPDiscoverPayload)
	assert.NoError(t, err)
	req4 := role.NewRequest4(req)
	res := role.Handle4(req4)
	assert.NotNil(t, res)
	assert.Equal(t, "10.100.0.100", res.YourIPAddr.String())
	ones, bits := res.SubnetMask().Size()
	assert.Equal(t, 24, ones)
	assert.Equal(t, 32, bits)
	assert.Equal(t, "b2:b7:86:2c:d3:fa", res.ClientHWAddr.String())
	assert.Equal(t, 86400*time.Second, res.IPAddressLeaseTime(1*time.Second))
	assert.Equal(t, "test.gravity.beryju.io", res.DomainName())
}
