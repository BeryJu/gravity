package dhcp_test

import (
	"fmt"
	"net"
	"net/netip"
	"testing"
	"time"

	"beryju.io/gravity/api"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dhcp"
	"beryju.io/gravity/pkg/roles/dhcp/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/stretchr/testify/assert"
)

var DHCPRequestPayload = []byte{1, 1, 6, 0, 136, 9, 170, 249, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 68, 144, 187, 102, 50, 4, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 99, 130, 83, 99, 53, 1, 3, 55, 9, 1, 121, 3, 6, 15, 108, 114, 119, 252, 57, 2, 5, 220, 61, 7, 1, 68, 144, 187, 102, 50, 4, 50, 4, 10, 120, 20, 64, 51, 4, 0, 118, 167, 0, 12, 14, 106, 101, 110, 115, 45, 105, 112, 104, 111, 110, 101, 45, 49, 50, 255, 0, 0, 0, 0}

func TestDHCPRequest(t *testing.T) {
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

	req, err := dhcpv4.FromBytes(DHCPRequestPayload)
	assert.NoError(t, err)
	req4 := role.NewRequest4(req)
	res := role.Handle4(req4)
	assert.NotNil(t, res)
	assert.Equal(t, "10.100.0.100", res.YourIPAddr.String())
	ones, bits := res.SubnetMask().Size()
	assert.Equal(t, 24, ones)
	assert.Equal(t, 32, bits)
	assert.Equal(t, "44:90:bb:66:32:04", res.ClientHWAddr.String())
	assert.Equal(t, 86400*time.Second, res.IPAddressLeaseTime(1*time.Second))
}

func TestDHCPRequest_Hook(t *testing.T) {
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
			Hook: `function onDHCPRequestAfter(req, res) {
				res.BootFileName = "foo"
			}`,
		}),
	))
	tests.PanicIfError(role.Start(ctx, RoleConfig()))
	defer role.Stop()

	req, err := dhcpv4.FromBytes(DHCPRequestPayload)
	assert.NoError(t, err)
	req4 := role.NewRequest4(req)
	res := role.Handle4(req4)
	assert.NotNil(t, res)
	assert.Equal(t, "10.100.0.100", res.YourIPAddr.String())
	ones, bits := res.SubnetMask().Size()
	assert.Equal(t, 24, ones)
	assert.Equal(t, 32, bits)
	assert.Equal(t, "44:90:bb:66:32:04", res.ClientHWAddr.String())
	assert.Equal(t, 86400*time.Second, res.IPAddressLeaseTime(1*time.Second))
	assert.Equal(t, "foo", res.BootFileName)
}

func TestDHCPRequest_Hook_UniFi(t *testing.T) {
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
			Hook: `const UniFiPrefix = [0x01, 0x04];
			const UniFiIP = net.parseIP("192.168.1.100", "v4");
			function onDHCPRequestAfter(req, res) {
				res.UpdateOption(dhcp.Opt(43, [...UniFiPrefix, ...UniFiIP]));
			}`,
		}),
	))
	tests.PanicIfError(role.Start(ctx, RoleConfig()))
	defer role.Stop()

	req, err := dhcpv4.FromBytes(DHCPRequestPayload)
	assert.NoError(t, err)
	req4 := role.NewRequest4(req)
	res := role.Handle4(req4)
	assert.NotNil(t, res)
	assert.Equal(t, "10.100.0.100", res.YourIPAddr.String())
	ones, bits := res.SubnetMask().Size()
	assert.Equal(t, 24, ones)
	assert.Equal(t, 32, bits)
	assert.Equal(t, "44:90:bb:66:32:04", res.ClientHWAddr.String())
	assert.Equal(t, 86400*time.Second, res.IPAddressLeaseTime(1*time.Second))
	vnd := res.GetOneOption(dhcpv4.OptionVendorSpecificInformation)
	assert.Equal(t, []byte{1, 4, 0xc0, 0xa8, 0x1, 0x64}, vnd)
}

func TestDHCPRequest_Hook_iPXE(t *testing.T) {
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
			Hook: `function onDHCPRequestAfter(req, res) {
				if (dhcp.GetString(77, req.Options) == "iPXE") {
					res.BootFileName = "foo"
					res.UpdateOption(dhcp.Opt(67, strconv.toBytes("foo")))
				}
			}`,
		}),
	))
	tests.PanicIfError(role.Start(ctx, RoleConfig()))
	defer role.Stop()

	req, err := dhcpv4.FromBytes(DHCPRequestPayload)
	req.UpdateOption(dhcpv4.OptUserClass("iPXE"))
	assert.NoError(t, err)
	req4 := role.NewRequest4(req)
	res := role.Handle4(req4)
	assert.NotNil(t, res)
	assert.Equal(t, "10.100.0.100", res.YourIPAddr.String())
	ones, bits := res.SubnetMask().Size()
	assert.Equal(t, 24, ones)
	assert.Equal(t, 32, bits)
	assert.Equal(t, "44:90:bb:66:32:04", res.ClientHWAddr.String())
	assert.Equal(t, 86400*time.Second, res.IPAddressLeaseTime(1*time.Second))
	assert.Equal(t, "foo", res.BootFileName)
	assert.Equal(t, "foo", res.BootFileNameOption())
}

func TestDHCPRequestDNS(t *testing.T) {
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

	req, err := dhcpv4.FromBytes(DHCPRequestPayload)
	assert.NoError(t, err)
	req4 := role.NewRequest4(req)
	res := role.Handle4(req4)
	assert.NotNil(t, res)
	assert.Equal(t, "10.100.0.100", res.YourIPAddr.String())
	ones, bits := res.SubnetMask().Size()
	assert.Equal(t, 24, ones)
	assert.Equal(t, 32, bits)
	assert.Equal(t, "44:90:bb:66:32:04", res.ClientHWAddr.String())
	assert.Equal(t, 86400*time.Second, res.IPAddressLeaseTime(1*time.Second))
	assert.Equal(t, "test.gravity.beryju.io", res.DomainName())
}

func TestDHCPRequestDNS_ChangedScope(t *testing.T) {
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
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyScopes,
			"test2",
		).String(),
		tests.MustJSON(dhcp.Scope{
			SubnetCIDR: "10.200.0.0/24",
			TTL:        86400,
			DNS: &dhcp.ScopeDNS{
				Zone:              "test2.gravity.beryju.io",
				AddZoneInHostname: true,
			},
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "10.200.0.100",
				"range_end":   "10.200.0.250",
			},
		}),
	))

	tests.PanicIfError(role.Start(ctx, RoleConfig()))
	defer role.Stop()

	// First ensure the lease is created as we expect
	req, err := dhcpv4.FromBytes(DHCPRequestPayload)
	assert.NoError(t, err)
	req4 := role.NewRequest4(req)
	res := role.Handle4(req4)
	assert.NotNil(t, res)
	assert.Equal(t, "10.100.0.100", res.YourIPAddr.String())
	ones, bits := res.SubnetMask().Size()
	assert.Equal(t, 24, ones)
	assert.Equal(t, 32, bits)
	assert.Equal(t, "44:90:bb:66:32:04", res.ClientHWAddr.String())
	assert.Equal(t, 86400*time.Second, res.IPAddressLeaseTime(1*time.Second))
	assert.Equal(t, "test.gravity.beryju.io", res.DomainName())

	// Now we're requesting an IP from another subnet, so the lease should move
	req.GatewayIPAddr = net.ParseIP("10.200.0.1")
	req4 = role.NewRequest4(req)
	res = role.Handle4(req4)
	assert.NotNil(t, res)
	assert.Equal(t, "10.200.0.100", res.YourIPAddr.String())
	ones, bits = res.SubnetMask().Size()
	assert.Equal(t, 24, ones)
	assert.Equal(t, 32, bits)
	assert.Equal(t, "44:90:bb:66:32:04", res.ClientHWAddr.String())
	assert.Equal(t, 86400*time.Second, res.IPAddressLeaseTime(1*time.Second))
	assert.Equal(t, "test2.gravity.beryju.io", res.DomainName())
}

func TestDHCP_Parallel(t *testing.T) {
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

	for i := 0; i < 10; i++ {
		t.Run(fmt.Sprintf("test %v", i), func(t *testing.T) {
			t.Parallel()
			time.Sleep(2 * time.Second)
			t.Logf("p done %v\n", time.Now())
		})
	}

	tester := func(cidr string) {
		c, err := netip.ParsePrefix(cidr)
		if err != nil {
			panic(err)
		}
		t.Run(cidr, func(t *testing.T) {
			t.Parallel()
			for i := 1; i < 100; i++ {
				rr := &dhcpv4.DHCPv4{
					GatewayIPAddr: net.ParseIP(c.Addr().Next().String()),
					ClientHWAddr:  tests.RandomMAC(),
					OpCode:        dhcpv4.OpcodeBootRequest,
				}
				rr.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeRequest))
				req4 := role.NewRequest4(rr)
				res := role.Handle4(req4)
				assert.NotNil(t, res)
				a, err := netip.ParseAddr(res.YourIPAddr.String())
				if err != nil {
					panic(err)
				}
				assert.True(t, c.Contains(a))
			}
		})
	}
	tester("10.100.0.0/24")
	tester("10.100.1.0/24")
	tester("10.100.2.0/24")
}

func TestDHCPRequest_Options(t *testing.T) {
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
			Options: []*types.DHCPOption{
				{
					TagName: types.TagNameRouter,
					Value:   api.PtrString("1.2.3.4"),
				},
				{
					TagName: types.TagNameBootfile,
					Value:   api.PtrString("foo"),
				},
				{
					TagName: "",
					Value:   api.PtrString("1.2.3.4"),
				},
				{
					TagName: "foo",
					Value:   api.PtrString("1.2.3.4"),
				},
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

	req, err := dhcpv4.FromBytes(DHCPRequestPayload)
	assert.NoError(t, err)
	req4 := role.NewRequest4(req)
	res := role.Handle4(req4)
	assert.NotNil(t, res)
	assert.Equal(t, "10.100.0.100", res.YourIPAddr.String())
	assert.Equal(t, "ffffff00", res.SubnetMask().String())
	assert.Equal(t, "44:90:bb:66:32:04", res.ClientHWAddr.String())
	assert.Equal(t, "1.2.3.4", res.Router()[0].String())
	assert.Equal(t, "foo", res.BootFileName)
	assert.Equal(t, "foo", res.BootFileNameOption())
	assert.Equal(t, 86400*time.Second, res.IPAddressLeaseTime(1*time.Second))
}
