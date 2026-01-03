package discovery_test

import (
	"net"
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dhcp"
	dhcpTypes "beryju.io/gravity/pkg/roles/dhcp/types"
	"beryju.io/gravity/pkg/roles/discovery"
	"beryju.io/gravity/pkg/roles/discovery/types"
	"beryju.io/gravity/pkg/roles/dns"
	dnsTypes "beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/gorilla/securecookie"
	"github.com/stretchr/testify/assert"
)

const (
	DockerNetworkCIDR = "10.200.0.0/28"
)

func TestDiscoveryConvert(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	role := discovery.New(rootInst.ForRole("discovery", ctx))
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, []byte(tests.MustJSON(discovery.RoleConfig{
		Enabled: true,
	}))))
	defer role.Stop()

	inst := rootInst.ForRole("test", ctx)

	// Create DNS Zone to register host in
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			dnsTypes.KeyRole,
			dnsTypes.KeyZones,
			"example.com.",
		).String(),
		tests.MustJSON(dns.Zone{
			HandlerConfigs: []map[string]interface{}{
				{
					"type": "etcd",
				},
			},
		}),
	))
	// Create DNS Zone for reverse
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			dnsTypes.KeyRole,
			dnsTypes.KeyZones,
			"200.10.in-addr.arpa.",
		).String(),
		tests.MustJSON(dns.Zone{
			HandlerConfigs: []map[string]interface{}{
				{
					"type": "etcd",
				},
			},
		}),
	))

	// Create DHCP Scope to register host in
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			dhcpTypes.KeyRole,
			dhcpTypes.KeyScopes,
			"test",
		).String(),
		tests.MustJSON(dhcp.Scope{
			SubnetCIDR: DockerNetworkCIDR,
			TTL:        86400,
			DNS: &dhcp.ScopeDNS{
				Zone: "example.com.",
			},
			IPAM: map[string]string{
				"type":        "internal",
				"range_start": "10.200.0.1",
				"range_end":   "10.200.0.250",
			},
		}),
	))

	// Start DNS & DHCP to register events
	dnsr := dns.New(rootInst.ForRole("dns", ctx))
	assert.NotNil(t, dnsr)
	assert.Nil(t, dnsr.Start(ctx, []byte(tests.MustJSON(dns.RoleConfig{
		Port: -1,
	}))))
	defer dnsr.Stop()

	dhcpr := dhcp.New(rootInst.ForRole("dhcp", ctx))
	assert.NotNil(t, dhcpr)
	assert.Nil(t, dhcpr.Start(ctx, []byte(tests.MustJSON(dhcp.RoleConfig{
		Port: -1,
	}))))
	defer dhcpr.Stop()

	// Manually create a device
	mac := net.HardwareAddr(securecookie.GenerateRandomKey(6)).String()
	tests.PanicIfError(inst.KV().Put(
		ctx,
		inst.KV().Key(
			types.KeyRole,
			types.KeyDevices,
			"test",
		).String(),
		tests.MustJSON(discovery.Device{
			IP:       "10.200.0.1",
			Hostname: "foo",
			MAC:      mac,
		}),
	))

	err := role.APIDevicesApply().Interact(ctx, discovery.APIDevicesApplyInput{
		Identifier: "test",
		To:         "dhcp",
		DHCPScope:  "test",
		DNSZone:    "example.com.",
	}, &struct{}{})
	assert.NoError(t, err)

	// Check DHCP lease
	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			dhcpTypes.KeyRole,
			dhcpTypes.KeyLeases,
			mac,
		),
		dhcp.Lease{
			ScopeKey:    "test",
			Address:     "10.200.0.1",
			Hostname:    "foo",
			Expiry:      0,
			Description: "",
		},
	)
	// Check forward DNS Record
	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			dnsTypes.KeyRole,
			dnsTypes.KeyZones,
			"example.com.",
			"foo",
			"A",
			mac,
		),
		dns.Record{
			Data: "10.200.0.1",
		},
	)
	// Check reverse DNS Record
	tests.AssertEtcd(
		t,
		inst.KV(),
		inst.KV().Key(
			dnsTypes.KeyRole,
			dnsTypes.KeyZones,
			"200.10.in-addr.arpa.",
			"1.0",
			"PTR",
			mac,
		),
		dns.Record{
			Data: "foo.example.com.",
		},
	)
}
