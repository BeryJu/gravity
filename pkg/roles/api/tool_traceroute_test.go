package api_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestAPIToolTraceroute(t *testing.T) {
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("api", ctx)
	role := api.New(inst)
	role.Start(ctx, []byte{})
	defer role.Stop()

	var output api.APIToolTracerouteOutput
	role.APIToolTraceroute().Interact(ctx, api.APIToolTracerouteInput{
		Host: "localhost",
	}, &output)
	assert.NotNil(t, output)
}
