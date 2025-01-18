package tests

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"beryju.io/gravity/api"
	"beryju.io/gravity/tests/gravity"
	"github.com/docker/docker/api/types/container"
	dockernetwork "github.com/docker/docker/api/types/network"
	"github.com/docker/docker/api/types/strslice"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/network"
)

func TestDHCP_Simple(t *testing.T) {
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

	g := gravity.New(t, gravity.WithNet(net))

	ac := g.APIClient()
	// Create test network
	_, err = ac.RolesDhcpApi.DhcpPutScopes(ctx).DhcpAPIScopesPutInput(api.DhcpAPIScopesPutInput{
		SubnetCidr: "10.100.0.0/24",
		Ttl:        86400,
		Ipam: map[string]string{
			"range_end":   "10.100.0.200",
			"range_start": "10.100.0.100",
			"type":        "internal",
			"should_ping": "true",
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
				Repo:       "gravity-dhcp-client",
				KeepImage:  true,
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

	_, body := ExecCommand(t, tester, []string{"dhclient", "-v"})
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

func TestDHCP_Parallel(t *testing.T) {
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

	g := gravity.New(t, gravity.WithNet(net))

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

	for i := 0; i < 50; i++ {
		t.Run(fmt.Sprintf("test %d", i), func(t *testing.T) {
			t.Parallel()
			tester, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
				ContainerRequest: testcontainers.ContainerRequest{
					FromDockerfile: testcontainers.FromDockerfile{
						Context:    "../hack/e2e/",
						Dockerfile: "dhcp-client.Dockerfile",
						Repo:       "gravity-dhcp-client",
						KeepImage:  true,
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

			_, body := ExecCommand(t, tester, []string{"dhclient", "-v"})
			assert.Contains(t, string(body), "DHCPOFFER of")
			assert.Contains(t, string(body), "DHCPREQUEST for")
			assert.Contains(t, string(body), "DHCPACK of")
			assert.Contains(t, string(body), "bound to")
		})
	}

	// // Check correct lease exists
	// defer func() {
	// 	sc, _, err := ac.RolesDhcpApi.DhcpGetLeases(ctx).Scope("network-A").Execute()
	// 	assert.NoError(t, err)
	// 	assert.Len(t, sc.Leases, 50)
	// }()
}

func TestDHCP_Relay(t *testing.T) {
	ctx := Context(t)

	netA, err := network.New(
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
	testcontainers.CleanupNetwork(t, netA)

	netB, err := network.New(
		ctx,
		network.WithIPAM(&dockernetwork.IPAM{
			Driver: "default",
			Config: []dockernetwork.IPAMConfig{
				{
					Subnet: "10.101.0.0/24",
				},
			},
		}),
		network.WithAttachable(),
	)
	assert.NoError(t, err)
	testcontainers.CleanupNetwork(t, netB)

	g := gravity.New(t,
		gravity.WithEnv("GRAVITY_DEBUG_DHCP_GATEWAY_REPLY_CIADDR", "true"),
		gravity.WithNet(netA))
	gip, err := g.Container().ContainerIP(ctx)
	assert.NoError(t, err)

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
	_, err = ac.RolesDhcpApi.DhcpPutScopes(ctx).DhcpAPIScopesPutInput(api.DhcpAPIScopesPutInput{
		SubnetCidr: "10.101.0.0/24",
		Ttl:        86400,
		Ipam: map[string]string{
			"range_end":   "10.101.0.200",
			"range_start": "10.101.0.100",
			"type":        "internal",
		},
		Options: []api.TypesDHCPOption{},
	}).Scope("network-B").Execute()
	assert.NoError(t, err)

	// DHCP relay
	relay, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			FromDockerfile: testcontainers.FromDockerfile{
				Context:    "../hack/e2e/",
				Dockerfile: "dhcp-relay.Dockerfile",
				Repo:       "gravity-dhcp-relay",
				KeepImage:  true,
			},
			Cmd:      []string{"-d", gip},
			Networks: []string{netA.Name, netB.Name},
			HostConfigModifier: func(hostConfig *container.HostConfig) {
				hostConfig.CapAdd = strslice.StrSlice{"NET_ADMIN"}
			},
		},
		Started: true,
	})
	testcontainers.CleanupContainer(t, relay)
	assert.NoError(t, err)

	// DHCP tester
	tester, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			FromDockerfile: testcontainers.FromDockerfile{
				Context:    "../hack/e2e/",
				Dockerfile: "dhcp-client.Dockerfile",
				Repo:       "gravity-dhcp-client",
				KeepImage:  true,
			},
			Networks: []string{netB.Name},
			HostConfigModifier: func(hostConfig *container.HostConfig) {
				hostConfig.CapAdd = strslice.StrSlice{"NET_ADMIN"}
			},
		},
		Started: true,
	})
	testcontainers.CleanupContainer(t, tester)
	assert.NoError(t, err)

	_, body := ExecCommand(t, tester, []string{"dhclient", "-v"})
	assert.Contains(t, string(body), "DHCPOFFER of")
	assert.Contains(t, string(body), "DHCPREQUEST for")
	assert.Contains(t, string(body), "DHCPACK of")
	assert.Contains(t, string(body), "bound to")

	// Check correct lease exists
	sc, _, err := ac.RolesDhcpApi.DhcpGetLeases(ctx).Scope("network-B").Execute()
	assert.NoError(t, err)
	assert.Len(t, sc.Leases, 1)
	assert.Equal(t, "10.101.0.100", sc.Leases[0].Address)
}

func TestDHCP_WOL(t *testing.T) {
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

	g := gravity.New(t, gravity.WithNet(net))

	ac := g.APIClient()
	// Create test network
	_, err = ac.RolesDhcpApi.DhcpPutScopes(ctx).DhcpAPIScopesPutInput(api.DhcpAPIScopesPutInput{
		SubnetCidr: "10.100.0.0/24",
		Ttl:        86400,
		Ipam: map[string]string{
			"range_end":   "10.100.0.200",
			"range_start": "10.100.0.100",
			"type":        "internal",
			"should_ping": "true",
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
				Repo:       "gravity-dhcp-client",
				KeepImage:  true,
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

	_, body := ExecCommand(t, tester, []string{"dhclient", "-v"})
	assert.Contains(t, string(body), "DHCPOFFER of")
	assert.Contains(t, string(body), "DHCPREQUEST for")
	assert.Contains(t, string(body), "DHCPACK of")
	assert.Contains(t, string(body), "bound to")

	// Check correct lease exists
	sc, _, err := ac.RolesDhcpApi.DhcpGetLeases(ctx).Scope("network-A").Execute()
	assert.NoError(t, err)
	assert.Len(t, sc.Leases, 1)
	assert.Equal(t, "10.100.0.100", sc.Leases[0].Address)

	// WOL
	wg := sync.WaitGroup{}
	wg.Add(2)
	go func() {
		defer wg.Done()
		cmd := "tcpdump -c 1 -UlnXi any ether proto 0x0842 or udp port 9 2>/dev/null "
		_, body = ExecCommand(t, tester, []string{"bash", "-c", cmd})
		assert.NotEqual(t, "", body)
	}()
	time.Sleep(3 * time.Second)
	go func() {
		defer wg.Done()
		_, err := ac.RolesDhcpApi.DhcpWolLeases(ctx).Scope("network-A").Identifier(sc.Leases[0].Identifier).Execute()
		assert.NoError(t, err)
	}()
	wg.Wait()
}
