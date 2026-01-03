package api_test

import (
	"runtime"
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestAPIToolPortmap(t *testing.T) {
	if runtime.GOOS != "linux" {
		t.Skip("Test only supported on linux")
	}
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("api", ctx)
	role := api.New(inst)
	tests.PanicIfError(role.Start(ctx, []byte{}))
	defer role.Stop()

	var output api.APIToolPortmapOutput
	tests.PanicIfError(role.APIToolPortmap().Interact(ctx, api.APIToolPortmapInput{
		Host: "localhost",
	}, &output))
	assert.NotNil(t, output)
}
