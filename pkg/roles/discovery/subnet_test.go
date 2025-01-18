package discovery_test

import (
	"testing"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/dhcp"
	dhcpTypes "beryju.io/gravity/pkg/roles/dhcp/types"
	"beryju.io/gravity/pkg/roles/discovery"
	"beryju.io/gravity/pkg/roles/dns"
	dnsTypes "beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

const (
	DockerNetworkCIDR = "10.200.0.0/28"

	DockerIPCoreDNS = "10.200.0.4"
)

func TestDiscoveryConvert(t *testing.T) {
	defer tests.Setup(t)()
	if !tests.HasLocalDocker() {
		t.Skip("Local docker required")
		return
	}
	extconfig.Get().ListenOnlyMode = false
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

	// Run discovery
	sub := role.NewSubnet("docker-test")
	sub.CIDR = DockerNetworkCIDR
	sub.DNSResolver = DockerIPCoreDNS
	devices := sub.RunDiscovery(ctx)

	err := role.APIDevicesApply().Interact(ctx, discovery.APIDevicesApplyInput{
		Identifier: devices[0].Identifier,
		To:         "dhcp",
		DHCPScope:  "test",
		DNSZone:    "example.com.",
	}, &struct{}{})
	assert.NoError(t, err)
}
