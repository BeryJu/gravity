//go:build e2e
// +build e2e

package tests

import (
	"io"
	"testing"

	"beryju.io/gravity/api"
	"github.com/docker/docker/api/types/container"
	dockernetwork "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
)

func TestDHCP_Single(t *testing.T) {
	ctx := Context(t)

	net, err := network.New(
		ctx,
		network.WithIPAM(&dockernetwork.IPAM{
			Driver: "default",
			Config: []dockernetwork.IPAMConfig{
				{
					Subnet: "10.100.0.0/24",
				},
			},
		}),
		network.WithAttachable(),
	)
	assert.NoError(t, err)
	testcontainers.CleanupNetwork(t, net)

	g := RunGravity(t, net)

	ac := g.APIClient()
	// Create test network
	_, err = ac.RolesDhcpApi.DhcpPutScopes(ctx).DhcpAPIScopesPutInput(api.DhcpAPIScopesPutInput{
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

	// DHCP tester
	tester, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			FromDockerfile: testcontainers.FromDockerfile{
				Context:    "../hack/e2e/",
				Dockerfile: "dhcp-client.Dockerfile",
			},
			Networks: []string{net.Name},
			HostConfigModifier: func(hostConfig *container.HostConfig) {
				hostConfig.CapAdd = strslice.StrSlice{"NET_ADMIN"}
			},
		},
		Started: true,
	})
	testcontainers.CleanupContainer(t, tester)
	assert.NoError(t, err)

	_, out, err := tester.Exec(ctx, []string{"dhclient", "-v"})
	assert.NoError(t, err)
	body, err := io.ReadAll(out)
	assert.NoError(t, err)
	assert.Contains(t, string(body), "DHCPOFFER of")
	assert.Contains(t, string(body), "DHCPREQUEST for")
	assert.Contains(t, string(body), "DHCPACK of")
	assert.Contains(t, string(body), "bound to")

	// Check correct lease exists
	sc, _, err := ac.RolesDhcpApi.DhcpGetLeases(ctx).Scope("network-A").Execute()
	assert.NoError(t, err)
	assert.Len(t, sc.Leases, 1)
	assert.Equal(t, "10.100.0.100", sc.Leases[0].Address)
}
