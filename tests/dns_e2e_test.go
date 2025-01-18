package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
)

func TestDNS_Single(t *testing.T) {
	ctx := Context(t)
	RunGravity(t, nil)

	// DHCP tester
	tester, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			FromDockerfile: testcontainers.FromDockerfile{
				Context:    "../hack/e2e/",
				Dockerfile: "dns.Dockerfile",
			},
		},
		Started: true,
	})
	testcontainers.CleanupContainer(t, tester)
	assert.NoError(t, err)

	_, dig := ExecCommand(t, tester, []string{"dig", "+short", "10.0.0.1.nip.io"})
	assert.Equal(t, "10.0.0.1\n", string(dig))
}
