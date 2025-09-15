package gravity

import (
	"context"
	"encoding/base64"
	"fmt"
	"os"
	"path/filepath"
	"testing"
	"time"

	"beryju.io/gravity/api"
	"beryju.io/gravity/pkg/extconfig"
	"github.com/docker/docker/api/types/container"
	"github.com/gorilla/securecookie"
	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

var (
	gravityPassword string
	gravityToken    string
)

func init() {
	gravityPassword = base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32))
	gravityToken = base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32))
}

func Password() string {
	return gravityPassword
}

func Token() string {
	return gravityToken
}

type Gravity struct {
	name      string
	container testcontainers.Container
	t         *testing.T
}

type GravityOption func(t *testing.T, req *testcontainers.ContainerRequest)

func WithEnv(key string, value string) GravityOption {
	return func(t *testing.T, req *testcontainers.ContainerRequest) {
		req.Env[key] = value
	}
}

func WithNet(net *testcontainers.DockerNetwork) GravityOption {
	return func(t *testing.T, req *testcontainers.ContainerRequest) {
		req.Networks = append(req.Networks, net.Name)
	}
}

func WithHostname(name string) GravityOption {
	return func(t *testing.T, req *testcontainers.ContainerRequest) {
		req.ConfigModifier = func(c *container.Config) {
			c.Hostname = name
		}
		req.LogConsumerCfg = &testcontainers.LogConsumerConfig{
			Consumers: []testcontainers.LogConsumer{
				&StdoutLogConsumer{T: t, Prefix: name},
			},
		}
		req.Name = name
	}
}

func WithoutWait() GravityOption {
	return func(t *testing.T, req *testcontainers.ContainerRequest) {
		req.WaitingFor = nil
	}
}

func New(t *testing.T, opts ...GravityOption) *Gravity {
	ctx := context.Background()
	cwd, err := os.Getwd()
	assert.NoError(t, err)

	req := testcontainers.ContainerRequest{
		Image:        "gravity:e2e-test",
		ExposedPorts: []string{"8008", "8009"},
		WaitingFor:   wait.ForHTTP("/healthz/ready").WithPort("8009").WithStartupTimeout(3 * time.Minute),
		Env: map[string]string{
			"LOG_LEVEL":      "debug",
			"ADMIN_PASSWORD": gravityPassword,
			"ADMIN_TOKEN":    gravityToken,
			"GOCOVERDIR":     "/coverage",
		},
		Networks: []string{},
		HostConfigModifier: func(hostConfig *container.HostConfig) {
			hostConfig.Binds = []string{
				fmt.Sprintf("%s:/coverage", filepath.Join(cwd, "/coverage")),
			}
		},
	}

	WithHostname("gravity-1")(t, &req)
	for _, opt := range opts {
		opt(t, &req)
	}
	name := req.Name
	req.Name = fmt.Sprintf("gravity-%s-%s", t.Name(), req.Name)

	gravityContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	assert.NoError(t, err)
	testcontainers.CleanupContainer(t, gravityContainer)

	return &Gravity{
		container: gravityContainer,
		t:         t,
		name:      name,
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
	config.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", gravityToken))
	config.UserAgent = fmt.Sprintf("gravity-e2e-testing/%s", extconfig.FullVersion())
	return api.NewAPIClient(config)
}

func (g *Gravity) EtcdID() string {
	nodes, _, err := g.APIClient().RolesEtcdAPI.EtcdGetMembers(context.Background()).Execute()
	assert.NoError(g.t, err)
	for _, node := range nodes.Members {
		if node.GetName() == g.name {
			return node.GetId()
		}
	}
	return ""
}

func (g *Gravity) Container() testcontainers.Container {
	return g.container
}
