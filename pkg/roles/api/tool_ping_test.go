package api_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestAPIToolPing(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("api", ctx)
	role := api.New(inst)
	tests.PanicIfError(role.Start(ctx, []byte{}))
	defer role.Stop()

	var output api.APIToolPingOutput
	tests.PanicIfError(role.APIToolPing().Interact(ctx, api.APIToolPingInput{
		Host: "localhost",
	}, &output))
	assert.NotNil(t, output)
}
