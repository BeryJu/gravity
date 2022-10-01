package discovery_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/discovery"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

const (
	DockerNetworkCIDR = "10.200.0.0/28"

	DockerIPCoreDNS = "10.200.0.4"
)

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
	sub.CIDR = DockerNetworkCIDR
	sub.DNSResolver = DockerIPCoreDNS
	devices := sub.RunDiscovery()
	assert.Equal(t, []discovery.Device{}, devices)
}
