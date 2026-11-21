package dhcp_test

import (
	"net"
	"testing"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dhcp"
	"beryju.io/gravity/pkg/roles/dhcp/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/stretchr/testify/assert"
)

// TestDHCPDiscover_MultipleIPs_SecondaryIP verifies that findScopeForRequest selects
// the scope matching the handler's BindIP, not the primary instance IP, when multiple
// instance IPs are configured and a broadcast discover arrives on the secondary handler.
func TestDHCPDiscover_MultipleIPs_SecondaryIP(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)
	role := dhcp.New(inst)

	extconfig.Get().Instance.IPs = []string{"10.0.1.1", "10.0.2.1"}

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(types.KeyRole, types.KeyScopes, "subnet1").String(),
		tests.MustJSON(dhcp.Scope{
			SubnetCIDR: "10.0.1.0/24",
			TTL:        86400,
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "10.0.1.100",
				"range_end":   "10.0.1.200",
			},
		}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(types.KeyRole, types.KeyScopes, "subnet2").String(),
		tests.MustJSON(dhcp.Scope{
			SubnetCIDR: "10.0.2.0/24",
			TTL:        86400,
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "10.0.2.100",
				"range_end":   "10.0.2.200",
			},
		}),
	))
	tests.PanicIfError(role.Start(ctx, RoleConfig()))
	defer role.Stop()

	req, err := dhcpv4.FromBytes(DHCPDiscoverPayload)
	assert.NoError(t, err)

	req4 := role.NewRequest4(req)
	req4.BindIP = "10.0.2.1"

	res := role.Handle4(req4)
	assert.NotNil(t, res)

	subnet2 := net.IPNet{IP: net.ParseIP("10.0.2.0").To4(), Mask: net.CIDRMask(24, 32)}
	assert.True(t, subnet2.Contains(res.YourIPAddr),
		"expected YourIPAddr %s to be in 10.0.2.0/24", res.YourIPAddr)
}

// TestDHCPDiscover_MultipleIPs_PrimaryIP verifies that the primary IP handler still
// correctly selects its own scope when multiple instance IPs are configured.
func TestDHCPDiscover_MultipleIPs_PrimaryIP(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)
	role := dhcp.New(inst)

	extconfig.Get().Instance.IPs = []string{"10.0.1.1", "10.0.2.1"}

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(types.KeyRole, types.KeyScopes, "subnet1").String(),
		tests.MustJSON(dhcp.Scope{
			SubnetCIDR: "10.0.1.0/24",
			TTL:        86400,
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "10.0.1.100",
				"range_end":   "10.0.1.200",
			},
		}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(types.KeyRole, types.KeyScopes, "subnet2").String(),
		tests.MustJSON(dhcp.Scope{
			SubnetCIDR: "10.0.2.0/24",
			TTL:        86400,
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "10.0.2.100",
				"range_end":   "10.0.2.200",
			},
		}),
	))
	tests.PanicIfError(role.Start(ctx, RoleConfig()))
	defer role.Stop()

	req, err := dhcpv4.FromBytes(DHCPDiscoverPayload)
	assert.NoError(t, err)

	req4 := role.NewRequest4(req)
	req4.BindIP = "10.0.1.1"

	res := role.Handle4(req4)
	assert.NotNil(t, res)

	subnet1 := net.IPNet{IP: net.ParseIP("10.0.1.0").To4(), Mask: net.CIDRMask(24, 32)}
	assert.True(t, subnet1.Contains(res.YourIPAddr),
		"expected YourIPAddr %s to be in 10.0.1.0/24", res.YourIPAddr)
}

// TestDHCPDiscover_MultipleIPs_SameMACDifferentScopes verifies that a device can
// move between scopes selected by their listener IP while retaining its MAC address.
func TestDHCPDiscover_MultipleIPs_SameMACDifferentScopes(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)
	role := dhcp.New(inst)

	extconfig.Get().Instance.IPs = []string{"10.0.1.1", "10.0.2.1"}

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(types.KeyRole, types.KeyScopes, "subnet1").String(),
		tests.MustJSON(dhcp.Scope{
			SubnetCIDR: "10.0.1.0/24",
			TTL:        86400,
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "10.0.1.100",
				"range_end":   "10.0.1.200",
			},
		}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(types.KeyRole, types.KeyScopes, "subnet2").String(),
		tests.MustJSON(dhcp.Scope{
			SubnetCIDR: "10.0.2.0/24",
			TTL:        86400,
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "10.0.2.100",
				"range_end":   "10.0.2.200",
			},
		}),
	))
	tests.PanicIfError(role.Start(ctx, RoleConfig()))
	defer role.Stop()

	firstRequest, err := dhcpv4.FromBytes(DHCPDiscoverPayload)
	assert.NoError(t, err)
	firstRequest4 := role.NewRequest4(firstRequest)
	firstRequest4.BindIP = "10.0.1.1"
	firstResponse := role.Handle4(firstRequest4)
	assert.NotNil(t, firstResponse)
	assert.Equal(t, "10.0.1.100", firstResponse.YourIPAddr.String())

	secondRequest, err := dhcpv4.FromBytes(DHCPDiscoverPayload)
	assert.NoError(t, err)
	secondRequest4 := role.NewRequest4(secondRequest)
	secondRequest4.BindIP = "10.0.2.1"
	secondResponse := role.Handle4(secondRequest4)
	assert.NotNil(t, secondResponse)
	assert.Equal(t, firstRequest.ClientHWAddr.String(), secondRequest.ClientHWAddr.String())
	assert.Equal(t, "10.0.2.100", secondResponse.YourIPAddr.String())
}

// TestDHCPDiscover_SameInstanceIP_DifferentVLANs verifies that relay gateway
// addresses select the correct VLAN-backed scope when requests arrive on the same
// instance IP with the same MAC address.
func TestDHCPDiscover_SameInstanceIP_DifferentVLANs(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)
	role := dhcp.New(inst)

	extconfig.Get().Instance.IPs = []string{"10.0.0.1"}

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(types.KeyRole, types.KeyScopes, "vlan10").String(),
		tests.MustJSON(dhcp.Scope{
			SubnetCIDR: "10.0.10.0/24",
			TTL:        86400,
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "10.0.10.100",
				"range_end":   "10.0.10.200",
			},
		}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(types.KeyRole, types.KeyScopes, "vlan20").String(),
		tests.MustJSON(dhcp.Scope{
			SubnetCIDR: "10.0.20.0/24",
			TTL:        86400,
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "10.0.20.100",
				"range_end":   "10.0.20.200",
			},
		}),
	))
	tests.PanicIfError(role.Start(ctx, RoleConfig()))
	defer role.Stop()

	firstRequest, err := dhcpv4.FromBytes(DHCPDiscoverPayload)
	assert.NoError(t, err)
	firstRequest.GatewayIPAddr = net.ParseIP("10.0.10.1")
	firstRequest4 := role.NewRequest4(firstRequest)
	firstRequest4.BindIP = "10.0.0.1"
	firstResponse := role.Handle4(firstRequest4)
	assert.NotNil(t, firstResponse)
	assert.Equal(t, "10.0.10.100", firstResponse.YourIPAddr.String())

	secondRequest, err := dhcpv4.FromBytes(DHCPDiscoverPayload)
	assert.NoError(t, err)
	secondRequest.GatewayIPAddr = net.ParseIP("10.0.20.1")
	secondRequest4 := role.NewRequest4(secondRequest)
	secondRequest4.BindIP = "10.0.0.1"
	secondResponse := role.Handle4(secondRequest4)
	assert.NotNil(t, secondResponse)
	assert.Equal(t, firstRequest.ClientHWAddr.String(), secondRequest.ClientHWAddr.String())
	assert.Equal(t, "10.0.20.100", secondResponse.YourIPAddr.String())
}

// TestDHCPDiscover_SameInstanceIP_VLANReservationOnlyInFirstScope verifies that
// a reservation applies only in its scope. The same device receives a dynamic
// lease when requesting through another VLAN on the same instance IP.
func TestDHCPDiscover_SameInstanceIP_VLANReservationOnlyInFirstScope(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("dhcp", ctx)
	role := dhcp.New(inst)

	extconfig.Get().Instance.IPs = []string{"10.0.0.1"}

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(types.KeyRole, types.KeyScopes, "vlan10").String(),
		tests.MustJSON(dhcp.Scope{
			SubnetCIDR: "10.0.10.0/24",
			TTL:        86400,
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "10.0.10.100",
				"range_end":   "10.0.10.200",
			},
		}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(types.KeyRole, types.KeyScopes, "vlan20").String(),
		tests.MustJSON(dhcp.Scope{
			SubnetCIDR: "10.0.20.0/24",
			TTL:        86400,
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "10.0.20.100",
				"range_end":   "10.0.20.200",
			},
		}),
	))
	assert.NoError(t, role.APILeasesPut().Interact(ctx, dhcp.APILeasesPutInput{
		Identifier: "b2:b7:86:2c:d3:fa",
		Scope:      "vlan10",
		Address:    "10.0.10.150",
		Expiry:     -1,
	}, &struct{}{}))
	tests.PanicIfError(role.Start(ctx, RoleConfig()))
	defer role.Stop()

	reservedRequest, err := dhcpv4.FromBytes(DHCPDiscoverPayload)
	assert.NoError(t, err)
	reservedRequest.GatewayIPAddr = net.ParseIP("10.0.10.1")
	reservedRequest4 := role.NewRequest4(reservedRequest)
	reservedRequest4.BindIP = "10.0.0.1"
	reservedResponse := role.Handle4(reservedRequest4)
	assert.NotNil(t, reservedResponse)
	assert.Equal(t, "10.0.10.150", reservedResponse.YourIPAddr.String())

	dynamicRequest, err := dhcpv4.FromBytes(DHCPDiscoverPayload)
	assert.NoError(t, err)
	dynamicRequest.GatewayIPAddr = net.ParseIP("10.0.20.1")
	dynamicRequest4 := role.NewRequest4(dynamicRequest)
	dynamicRequest4.BindIP = "10.0.0.1"
	dynamicResponse := role.Handle4(dynamicRequest4)
	assert.NotNil(t, dynamicResponse)
	assert.Equal(t, reservedRequest.ClientHWAddr.String(), dynamicRequest.ClientHWAddr.String())
	assert.Equal(t, "10.0.20.100", dynamicResponse.YourIPAddr.String())

	reservedRequestAgain, err := dhcpv4.FromBytes(DHCPDiscoverPayload)
	assert.NoError(t, err)
	reservedRequestAgain.GatewayIPAddr = net.ParseIP("10.0.10.1")
	reservedRequestAgain4 := role.NewRequest4(reservedRequestAgain)
	reservedRequestAgain4.BindIP = "10.0.0.1"
	reservedResponseAgain := role.Handle4(reservedRequestAgain4)
	assert.NotNil(t, reservedResponseAgain)
	assert.Equal(t, reservedRequest.ClientHWAddr.String(), reservedRequestAgain.ClientHWAddr.String())
	assert.Equal(t, "10.0.10.150", reservedResponseAgain.YourIPAddr.String())
}
