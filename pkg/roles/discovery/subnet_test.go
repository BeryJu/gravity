package discovery_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/discovery"
	"beryju.io/gravity/pkg/tests"
	dockertypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/filters"
	"github.com/docker/docker/client"
	"github.com/stretchr/testify/assert"
)

func getDockerTestCIDR() dockertypes.NetworkResource {
	cli, err := client.NewClientWithOpts(
		client.FromEnv,
		client.WithAPIVersionNegotiation(),
	)
	if err != nil {
		panic(err)
	}
	networks, err := cli.NetworkList(tests.Context(), dockertypes.NetworkListOptions{
		Filters: filters.NewArgs(filters.KeyValuePair{
			Key:   "label",
			Value: "io.beryju.gravity/testing=true",
		}),
	})
	if err != nil {
		panic(err)
	}
	return networks[0]
}

func Test_Discovery_Docker(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("discovery")
	role := discovery.New(inst)
	assert.NotNil(t, role)
	ctx := tests.Context()
	assert.Nil(t, role.Start(ctx, []byte(tests.MustJSON(discovery.RoleConfig{
		Enabled: true,
	}))))
	defer role.Stop()

	sub := role.NewSubnet("docker-test")
	sub.CIDR = getDockerTestCIDR().IPAM.Config[0].Subnet
	devices := sub.RunDiscovery()
	assert.Equal(t, []discovery.Device{}, devices)
}
