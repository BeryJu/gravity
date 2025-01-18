package tests

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"testing"

	"beryju.io/gravity/api"
	"beryju.io/gravity/pkg/extconfig"
	"github.com/docker/docker/api/types/container"
	"github.com/gorilla/securecookie"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/exec"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	GravityPassword string
	GravityToken    string
)

func init() {
	GravityPassword = base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32))
	GravityToken = base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32))
}

func Context(t *testing.T) context.Context {
	ctx, cn := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cn()
	})
	return ctx
}

type Gravity struct {
	container testcontainers.Container
	t         *testing.T
}

func RunGravity(t *testing.T, net *testcontainers.DockerNetwork) *Gravity {
	ctx := Context(t)
	cwd, err := os.Getwd()
	assert.NoError(t, err)

	req := testcontainers.ContainerRequest{
		Image:        "gravity:e2e-test",
		ExposedPorts: []string{"8008", "8009"},
		WaitingFor:   wait.ForHTTP("/healthz/ready").WithPort("8009"),
		Env: map[string]string{
			"LOG_LEVEL":      "debug",
			"ADMIN_PASSWORD": GravityPassword,
			"ADMIN_TOKEN":    GravityToken,
			"GOCOVERDIR":     "/coverage",
		},
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.Binds = []string{
				fmt.Sprintf("%s/coverage:/coverage", cwd),
			}
		},
	}

	if net != nil {
		req.Networks = []string{net.Name}
	}

	gravityContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	testcontainers.CleanupContainer(t, gravityContainer)
	assert.NoError(t, err)

	return &Gravity{
		container: gravityContainer,
		t:         t,
	}
}

func (g *Gravity) APIClient() *api.APIClient {
	ctx := context.Background()
	addr, err := g.container.Host(ctx)
	assert.NoError(g.t, err)
	port, err := g.container.MappedPort(ctx, "8008")
	assert.NoError(g.t, err)

	config := api.NewConfiguration()
	config.Debug = true
	config.Scheme = "http"
	config.Host = fmt.Sprintf("%s:%s", addr, port.Port())
	config.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", GravityToken))
	config.UserAgent = fmt.Sprintf("gravity-e2e-testing/%s", extconfig.FullVersion())
	return api.NewAPIClient(config)
}

func ExecCommand(t *testing.T, co testcontainers.Container, cmd []string, options ...exec.ProcessOption) (int, string) {
	ctx := Context(t)
	options = append(options, exec.Multiplexed())
	c, out, err := co.Exec(ctx, cmd, options...)
	assert.NoError(t, err)
	body, err := io.ReadAll(out)
	assert.NoError(t, err)
	return c, string(body)
}
