package discovery_test

import (
	"testing"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/discovery"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

const (
	DockerNetworkCIDR = "10.200.0.0/28"

	DockerIPCoreDNS = "10.200.0.4"
)

func TestDiscoveryDocker(t *testing.T) {
	if !tests.HasLocalDocker() {
		return
	}
	extconfig.Get().ListenOnlyMode = false
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("discovery", ctx)
	role := discovery.New(inst)
	assert.NotNil(t, role)
	assert.Nil(t, role.Start(ctx, []byte(tests.MustJSON(discovery.RoleConfig{
		Enabled: true,
	}))))
	defer role.Stop()

	sub := role.NewSubnet("docker-test")
	sub.CIDR = DockerNetworkCIDR
	sub.DNSResolver = DockerIPCoreDNS
	devices := sub.RunDiscovery()
	assert.Equal(t, "10.200.0.1", devices[0].IP)

	assert.Equal(t, "etcd.t.gravity.beryju.io", devices[1].Hostname)
	assert.Equal(t, "10.200.0.2", devices[1].IP)
	assert.Equal(t, "", devices[1].MAC)

	assert.Equal(t, "minio.t.gravity.beryju.io", devices[2].Hostname)
	assert.Equal(t, "10.200.0.3", devices[2].IP)
	assert.Equal(t, "", devices[2].MAC)

	assert.Equal(t, "coredns.t.gravity.beryju.io", devices[3].Hostname)
	assert.Equal(t, "10.200.0.4", devices[3].IP)
	assert.Equal(t, "", devices[3].MAC)
}
