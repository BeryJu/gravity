//go:build e2e
// +build e2e

package tests

import (
	"fmt"
	"testing"

	"beryju.io/gravity/tests/gravity"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
)

func TestCluster_Join(t *testing.T) {
	ctx := Context(t)

	net, err := network.New(ctx, network.WithAttachable())
	assert.NoError(t, err)
	testcontainers.CleanupNetwork(t, net)

	// Create initial gravity node
	gr := gravity.New(t, gravity.WithNet(net))

	// Create 2nd gravity node
	gravity.New(
		t,
		gravity.WithEnv("ETCD_JOIN_CLUSTER", fmt.Sprintf("%s,http://gravity-1:8008", gravity.Token())),
		gravity.WithHostname("gravity-2"),
		gravity.WithNet(net),
	)

	// Create 3rd gravity node
	gravity.New(
		t,
		gravity.WithEnv("ETCD_JOIN_CLUSTER", fmt.Sprintf("%s,http://gravity-1:8008", gravity.Token())),
		gravity.WithHostname("gravity-3"),
		gravity.WithNet(net),
	)

	// Check that all nodes are in the cluster
	ac := gr.APIClient()
	c, _, err := ac.ClusterAPI.ClusterGetClusterInfo(ctx).Execute()
	assert.NoError(t, err)
	assert.Len(t, c.Instances, 3)
}
