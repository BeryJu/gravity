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
	rootInst := instance.New()
	inst := rootInst.ForRole("discovery")
	role := discovery.New(inst)
	ctx := tests.Context()

	inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyDevices,
			tests.RandomString(),
		).String(),
		tests.MustJSON(discovery.Device{}),
	)

	var output discovery.APIDevicesGetOutput
	assert.NoError(t, role.APIDevicesGet().Interact(ctx, struct{}{}, &output))
	assert.NotNil(t, output)
}

func TestDeviceApplyDHCP(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("discovery")

	// Create DHCP role to register events
	dhcpRole := dhcp.New(rootInst.ForRole("dhcp"))
	dhcpRole.Start(tests.Context(), []byte{})
	defer dhcpRole.Stop()

	role := discovery.New(inst)

	name := tests.RandomString()
	inst.KV().Put(
		tests.Context(),
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
	)
	inst.KV().Put(
		tests.Context(),
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
	)

	assert.NoError(t, role.APIDevicesApply().Interact(tests.Context(), discovery.APIDevicesApplyInput{
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
	rootInst := instance.New()
	inst := rootInst.ForRole("discovery")

	// Create DHCP role to register events
	dhcpRole := dhcp.New(rootInst.ForRole("dhcp"))
	dhcpRole.Start(tests.Context(), []byte{})
	defer dhcpRole.Stop()
	// Create DNS role to register events
	dnsRole := dns.New(rootInst.ForRole("dns"))
	dnsRole.Start(tests.Context(), []byte{})
	defer dnsRole.Stop()

	role := discovery.New(inst)

	name := tests.RandomString()
	inst.KV().Put(
		tests.Context(),
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
	)
	inst.KV().Put(
		tests.Context(),
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
	)
	inst.KV().Put(
		tests.Context(),
		inst.KV().Key(
			dnsTypes.KeyRole,
			dnsTypes.KeyZones,
			"gravity.beryju.io.",
		).String(),
		tests.MustJSON(dns.Zone{}),
	)
	inst.KV().Put(
		tests.Context(),
		inst.KV().Key(
			dnsTypes.KeyRole,
			dnsTypes.KeyZones,
			"0.192.in-addr.arpa.",
		).String(),
		tests.MustJSON(dns.Zone{}),
	)

	assert.NoError(t, role.APIDevicesApply().Interact(tests.Context(), discovery.APIDevicesApplyInput{
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
		),
		dns.Record{
			Data: "test.gravity.beryju.io.",
		},
	)
}

func TestDeviceApplyDNSWithReverse(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("discovery")

	// Create DNS role to register events
	dnsRole := dns.New(rootInst.ForRole("dns"))
	dnsRole.Start(tests.Context(), []byte{})
	defer dnsRole.Stop()

	role := discovery.New(inst)

	name := tests.RandomString()
	inst.KV().Put(
		tests.Context(),
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
	)
	inst.KV().Put(
		tests.Context(),
		inst.KV().Key(
			dnsTypes.KeyRole,
			dnsTypes.KeyZones,
			"gravity.beryju.io.",
		).String(),
		tests.MustJSON(dns.Zone{}),
	)
	inst.KV().Put(
		tests.Context(),
		inst.KV().Key(
			dnsTypes.KeyRole,
			dnsTypes.KeyZones,
			"0.192.in-addr.arpa.",
		).String(),
		tests.MustJSON(dns.Zone{}),
	)

	assert.NoError(t, role.APIDevicesApply().Interact(tests.Context(), discovery.APIDevicesApplyInput{
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
		),
		dns.Record{
			Data: "test.gravity.beryju.io.",
		},
	)
}

func TestAPIDevicesDelete(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("discovery")
	role := discovery.New(inst)
	ctx := tests.Context()

	name := tests.RandomString()

	inst.KV().Put(
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
	)

	assert.NoError(t, role.APIDevicesDelete().Interact(tests.Context(), discovery.APIDevicesDeleteInput{
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
