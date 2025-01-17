package tests

import (
	"testing"

	"beryju.io/gravity/api"
	"github.com/efficientgo/e2e"
	"github.com/stretchr/testify/assert"
)

func TestDHCP_Single(t *testing.T) {
	t.Skip("Broken test")
	RunGravity(t)
	ctx := Context(t)

	ac := APIClient()
	// Create test network
	_, err := ac.RolesDhcpApi.DhcpPutScopes(ctx).DhcpAPIScopesPutInput(api.DhcpAPIScopesPutInput{
		SubnetCidr: "10.100.0.0/24",
		Ttl:        86400,
		Ipam: map[string]string{
			"range_end":   "10.100.0.200",
			"range_start": "10.100.0.100",
			"type":        "internal",
		},
		Options: []api.TypesDHCPOption{},
	}).Scope("network-A").Execute()
	assert.NoError(t, err)

	rb := env.Runnable("dhcp-tester").
		Init(e2e.StartOptions{
			Image: "gravity-testing:dhcp-client",
			Capabilities: []e2e.RunnableCapabilities{
				"NET_ADMIN",
			},
		})
	assert.NoError(t, e2e.StartAndWaitReady(rb))
	err = rb.Exec(e2e.NewCommand("dhclient", "-v"))
	assert.NoError(t, err)
}
