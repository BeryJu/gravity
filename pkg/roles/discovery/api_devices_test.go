package discovery_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dhcp"
	dhcpTypes "beryju.io/gravity/pkg/roles/dhcp/types"
	"beryju.io/gravity/pkg/roles/discovery"
	"beryju.io/gravity/pkg/roles/discovery/types"
	"beryju.io/gravity/pkg/roles/dns"
	dnsTypes "beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestAPIDevicesGet(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("discovery", ctx)
	role := discovery.New(inst)

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyDevices,
			tests.RandomString(),
		).String(),
		tests.MustJSON(discovery.Device{}),
	))

	var output discovery.APIDevicesGetOutput
	assert.NoError(t, role.APIDevicesGet().Interact(ctx, discovery.APIDevicesGetInput{}, &output))
	assert.NotNil(t, output)
}

func TestDeviceApplyDHCP(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("discovery", ctx)
	name := tests.RandomString()

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			dhcpTypes.KeyRole,
			dhcpTypes.KeyScopes,
			name,
		).String(),
		tests.MustJSON(dhcp.Scope{
			SubnetCIDR: "192.0.2.0/24",
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "192.0.2.1",
				"range_end":   "192.0.2.250",
			},
			DNS: &dhcp.ScopeDNS{
				Zone: "gravity.beryju.io.",
			},
		}),
	))

	// Create DHCP role to register events
	dhcpRole := dhcp.New(rootInst.ForRole("dhcp", ctx))
	tests.PanicIfError(dhcpRole.Start(ctx, []byte(tests.MustJSON(dhcp.RoleConfig{
		Port: 0,
	}))))
	defer dhcpRole.Stop()

	role := discovery.New(inst)

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyDevices,
			name,
		).String(),
		tests.MustJSON(discovery.Device{
			Hostname: "test.gravity.beryju.io",
			IP:       "192.0.2.1",
			MAC:      "aa:bb:cc",
		}),
	))

	assert.NoError(t, role.APIDevicesApply().Interact(ctx, discovery.APIDevicesApplyInput{
		Identifier: name,
		To:         "dhcp",
		DHCPScope:  name,
	}, &struct{}{}))

	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			dhcpTypes.KeyRole,
			dhcpTypes.KeyLeases,
			"aa:bb:cc",
		),
		dhcp.Lease{
			Hostname: "test",
			Address:  "192.0.2.1",
			ScopeKey: name,
		},
	)
}

func TestDeviceApplyDHCPWithDNS(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("discovery", ctx)
	name := tests.RandomString()

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			dhcpTypes.KeyRole,
			dhcpTypes.KeyScopes,
			name,
		).String(),
		tests.MustJSON(dhcp.Scope{
			SubnetCIDR: "192.0.2.0/24",
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "192.0.2.1",
				"range_end":   "192.0.2.250",
			},
			DNS: &dhcp.ScopeDNS{
				Zone: "gravity.beryju.io.",
			},
		}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			dnsTypes.KeyRole,
			dnsTypes.KeyZones,
			"gravity.beryju.io.",
		).String(),
		tests.MustJSON(dns.Zone{}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			dnsTypes.KeyRole,
			dnsTypes.KeyZones,
			"0.192.in-addr.arpa.",
		).String(),
		tests.MustJSON(dns.Zone{}),
	))

	// Create DHCP role to register events
	dhcpRole := dhcp.New(rootInst.ForRole("dhcp", ctx))
	tests.PanicIfError(dhcpRole.Start(ctx, []byte(tests.MustJSON(dhcp.RoleConfig{
		Port: 0,
	}))))
	defer dhcpRole.Stop()
	// Create DNS role to register events
	dnsRole := dns.New(rootInst.ForRole("dns", ctx))
	tests.PanicIfError(dnsRole.Start(ctx, []byte(tests.MustJSON(dns.RoleConfig{
		Port: 12123,
	}))))
	defer dnsRole.Stop()

	role := discovery.New(inst)

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyDevices,
			name,
		).String(),
		tests.MustJSON(discovery.Device{
			Hostname: "test.gravity.beryju.io",
			IP:       "192.0.2.1",
			MAC:      "aa:bb:cc",
		}),
	))

	assert.NoError(t, role.APIDevicesApply().Interact(ctx, discovery.APIDevicesApplyInput{
		Identifier: name,
		To:         "dhcp",
		DHCPScope:  name,
	}, &struct{}{}))

	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			dhcpTypes.KeyRole,
			dhcpTypes.KeyLeases,
			"aa:bb:cc",
		),
		dhcp.Lease{
			Hostname: "test",
			Address:  "192.0.2.1",
			ScopeKey: name,
		},
	)
	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			dnsTypes.KeyRole,
			dnsTypes.KeyZones,
			"gravity.beryju.io.",
			"test",
			"A",
			"aa:bb:cc",
		),
		dns.Record{
			Data: "192.0.2.1",
		},
	)
	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			dnsTypes.KeyRole,
			dnsTypes.KeyZones,
			"0.192.in-addr.arpa.",
			"1.2",
			"PTR",
			"aa:bb:cc",
		),
		dns.Record{
			Data: "test.gravity.beryju.io.",
		},
	)
}

func TestDeviceApplyDNSWithReverse(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("discovery", ctx)

	name := tests.RandomString()
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyDevices,
			name,
		).String(),
		tests.MustJSON(discovery.Device{
			Hostname: "test.gravity.beryju.io.",
			IP:       "192.0.2.1",
			MAC:      "aa:bb:cc",
		}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			dnsTypes.KeyRole,
			dnsTypes.KeyZones,
			"gravity.beryju.io.",
		).String(),
		tests.MustJSON(dns.Zone{}),
	))
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			dnsTypes.KeyRole,
			dnsTypes.KeyZones,
			"0.192.in-addr.arpa.",
		).String(),
		tests.MustJSON(dns.Zone{}),
	))

	// Create DNS role to register events
	dnsRole := dns.New(rootInst.ForRole("dns", ctx))
	tests.PanicIfError(dnsRole.Start(ctx, []byte(tests.MustJSON(dns.RoleConfig{
		Port: 0,
	}))))
	defer dnsRole.Stop()

	role := discovery.New(inst)
	assert.NoError(t, role.APIDevicesApply().Interact(ctx, discovery.APIDevicesApplyInput{
		Identifier: name,
		To:         "dns",
		DNSZone:    "gravity.beryju.io.",
	}, &struct{}{}))

	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			dnsTypes.KeyRole,
			dnsTypes.KeyZones,
			"gravity.beryju.io.",
			"test",
			"A",
			name,
		),
		dns.Record{
			Data: "192.0.2.1",
		},
	)
	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			dnsTypes.KeyRole,
			dnsTypes.KeyZones,
			"0.192.in-addr.arpa.",
			"1.2",
			"PTR",
			name,
		),
		dns.Record{
			Data: "test.gravity.beryju.io.",
		},
	)
}

func TestAPIDevicesDelete(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("discovery", ctx)
	role := discovery.New(inst)

	name := tests.RandomString()

	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyDevices,
			name,
		).String(),
		tests.MustJSON(discovery.Device{
			Hostname: "test",
			IP:       "192.0.2.1",
			MAC:      "aa:bb:cc",
		}),
	))

	assert.NoError(t, role.APIDevicesDelete().Interact(ctx, discovery.APIDevicesDeleteInput{
		Name: name,
	}, &struct{}{}))

	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			types.KeyRole,
			types.KeyDevices,
			name,
		),
	)
}
