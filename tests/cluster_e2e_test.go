//go:build e2e
// +build e2e

package tests

import (
	"errors"
	"fmt"
	"testing"

	"beryju.io/gravity/tests/gravity"
	"github.com/avast/retry-go/v4"
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
	c, _, err := ac.RolesEtcdAPI.EtcdGetMembers(ctx).Execute()
	assert.NoError(t, err)
	assert.Len(t, c.Members, 3)
}

func TestCluster_Join_NoWait(t *testing.T) {
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
		gravity.WithoutWait(),
	)

	// Create 3rd gravity node
	gravity.New(
		t,
		gravity.WithEnv("ETCD_JOIN_CLUSTER", fmt.Sprintf("%s,http://gravity-1:8008", gravity.Token())),
		gravity.WithHostname("gravity-3"),
		gravity.WithNet(net),
		gravity.WithoutWait(),
	)

	c, err := retry.DoWithData(
		func() (int, error) {
			// Check that all nodes are in the cluster
			ac := gr.APIClient()
			c, _, err := ac.RolesEtcdAPI.EtcdGetMembers(ctx).Execute()
			if err != nil {
				return 0, err
			}
			if len(c.Members) < 3 {
				return 0, errors.New("No enough members")
			}
			return len(c.Members), nil
		},
		retry.Attempts(50),
		retry.OnRetry(func(attempt uint, err error) {
			t.Logf("Checking cluster member count %d, %v", attempt, err)
		}),
	)
	assert.NoError(t, err)
	assert.Equal(t, 3, c)
}
