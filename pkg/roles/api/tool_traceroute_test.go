package api_test

import (
	"syscall"
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestAPIToolTraceroute(t *testing.T) {
	tests.Setup(t)
	if syscall.Getuid() != 0 {
		t.Skip("Traceroute requires root permissions")
	}
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("api", ctx)
	role := api.New(inst)
	tests.PanicIfError(role.Start(ctx, []byte{}))
	defer role.Stop()

	var output api.APIToolTracerouteOutput
	tests.PanicIfError(role.APIToolTraceroute().Interact(ctx, api.APIToolTracerouteInput{
		Host: "localhost",
	}, &output))
	assert.NotNil(t, output)
}
