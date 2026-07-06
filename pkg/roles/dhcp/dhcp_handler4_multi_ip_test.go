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
