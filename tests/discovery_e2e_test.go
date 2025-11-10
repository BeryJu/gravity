//go:build e2e

package tests

import (
	"testing"

	"beryju.io/gravity/tests/gravity"
	dockernetwork "github.com/docker/docker/api/types/network"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
)

func TestDiscovery_Simple(t *testing.T) {
	ctx := Context(t)

	net, err := network.New(
		ctx,
		network.WithIPAM(&dockernetwork.IPAM{
			Driver: "default",
			Config: []dockernetwork.IPAMConfig{
				{
					Subnet: "10.100.0.0/29",
				},
			},
		}),
		network.WithAttachable(),
	)
	assert.NoError(t, err)
	testcontainers.CleanupNetwork(t, net)

	gr := gravity.New(t,
		gravity.WithNet(net),
		// Use Docker gateway as DNS server to lookup name of other containers
		gravity.WithEnv("FALLBACK_DNS", "10.100.0.1:53"))

	tester, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			FromDockerfile: testcontainers.FromDockerfile{
				Context:    "../hack/e2e/",
				Dockerfile: "dns.Dockerfile",
				Repo:       "gravity-dns-client",
				KeepImage:  true,
			},
			Hostname: "client",
			Networks: []string{net.Name},
		},
		Started: true,
	})
	testcontainers.CleanupContainer(t, tester)
	assert.NoError(t, err)

	_, err = gr.APIClient().RolesDiscoveryAPI.
		DiscoverySubnetStart(ctx).
		Identifier("instance-subnet-gravity-1").
		Wait(true).
		Execute()
	assert.NoError(t, err)

	d, _, err := gr.APIClient().RolesDiscoveryAPI.DiscoveryGetDevices(ctx).Execute()
	assert.NoError(t, err)

	assert.Len(t, d.Devices, 3)
}
