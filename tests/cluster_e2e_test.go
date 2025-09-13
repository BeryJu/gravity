//go:build e2e
// +build e2e

package tests

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"beryju.io/gravity/tests/gravity"
	"github.com/docker/docker/api/types/container"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
	"github.com/testcontainers/testcontainers-go/wait"
)

func TestCluster_Join(t *testing.T) {
	ctx := Context(t)

	net, err := network.New(ctx, network.WithAttachable())
	assert.NoError(t, err)
	testcontainers.CleanupNetwork(t, net)

	// Create initial gravity node
	gr := gravity.New(t, gravity.WithNet(net))

	cwd, err := os.Getwd()
	assert.NoError(t, err)

	// Create 2nd gravity node
	gravity2, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "gravity:e2e-test",
			ExposedPorts: []string{"8008", "8009"},
			WaitingFor:   wait.ForHTTP("/healthz/ready").WithPort("8009"),
			Hostname:     "gravity-2",
			Networks:     []string{net.Name},
			Env: map[string]string{
				"LOG_LEVEL":         "debug",
				"ETCD_JOIN_CLUSTER": fmt.Sprintf("%s,http://gravity-1:8008", gravity.Token()),
				"GOCOVERDIR":        "/coverage",
			},
			HostConfigModifier: func(hostConfig *container.HostConfig) {
				hostConfig.Binds = []string{
					fmt.Sprintf("%s:/coverage", filepath.Join(cwd, "/coverage")),
				}
			},
		},
		Started: true,
	})
	testcontainers.CleanupContainer(t, gravity2)
	assert.NoError(t, err)

	// // Create 3rd gravity node
	// gravity3, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
	// 	ContainerRequest: testcontainers.ContainerRequest{
	// 		Image:        "gravity:e2e-test",
	// 		ExposedPorts: []string{"8008", "8009"},
	// 		WaitingFor:   wait.ForHTTP("/healthz/ready").WithPort("8009"),
	// 		Hostname:     "gravity-3",
	// 		Networks:     []string{net.Name},
	// 		Env: map[string]string{
	// 			"LOG_LEVEL":         "debug",
	// 			"ETCD_JOIN_CLUSTER": fmt.Sprintf("%s,http://gravity-1:8008", GravityToken),
	// 			"GOCOVERDIR":        "/coverage",
	// 		},
	// 		HostConfigModifier: func(hostConfig *container.HostConfig) {
	// 			hostConfig.Binds = []string{
	// 				fmt.Sprintf("%s:/coverage", filepath.Join(cwd, "/coverage")),
	// 			}
	// 		},
	// 	},
	// 	Started: true,
	// })
	// testcontainers.CleanupContainer(t, gravity3)
	// assert.NoError(t, err)

	// Check that all nodes are in the cluster
	ac := gr.APIClient()
	c, _, err := ac.ClusterAPI.ClusterGetClusterInfo(ctx).Execute()
	assert.NoError(t, err)
	assert.Len(t, c.Instances, 2)
}
