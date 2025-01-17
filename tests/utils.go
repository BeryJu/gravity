//go:build e2e
// +build e2e

package tests

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"

	"beryju.io/gravity/api"
	"beryju.io/gravity/pkg/extconfig"
	"github.com/efficientgo/e2e"
	"github.com/gorilla/securecookie"
	"github.com/stretchr/testify/assert"
)

var (
	GravityPassword string
	GravityToken    string
)

func init() {
	GravityPassword = base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32))
	GravityToken = base64.StdEncoding.EncodeToString(securecookie.GenerateRandomKey(32))
}

var gravity e2e.Runnable
var env *e2e.DockerEnvironment

func RunGravity(t *testing.T) {
	e, err := e2e.New()
	assert.NoError(t, err)
	env = e
	// Make sure resources (e.g docker containers, network, dir) are cleaned.
	t.Cleanup(e.Close)

	rb := e.Runnable("gravity").
		WithPorts(map[string]int{
			"http":    8008,
			"metrics": 8009,
		}).
		Init(e2e.StartOptions{
			Image: "gravity:e2e-test",
			EnvVars: map[string]string{
				"LOG_LEVEL":      "debug",
				"ADMIN_PASSWORD": GravityPassword,
				"ADMIN_TOKEN":    GravityToken,
				"GOCOVERDIR":     "/coverage",
			},
			Volumes: []string{
				"./coverage:/coverage",
			},
			Readiness: e2e.NewHTTPReadinessProbe(
				"metrics", "/healthz/ready", 200, 200,
			),
		})
	assert.NoError(t, e2e.StartAndWaitReady(rb))
	gravity = rb
}

func Context(t *testing.T) context.Context {
	ctx, cn := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cn()
	})
	return ctx
}

func APIClient() *api.APIClient {
	config := api.NewConfiguration()
	config.Debug = true
	config.Scheme = "http"
	config.Host = gravity.Endpoint("http")
	config.AddDefaultHeader("Authorization", fmt.Sprintf("Bearer %s", GravityToken))
	config.UserAgent = fmt.Sprintf("gravity-e2e-testing/%s", extconfig.FullVersion())
	return api.NewAPIClient(config)
}
