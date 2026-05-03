package debug_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/debug"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestRoleStartNoConfig(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("debug", ctx)
	role := debug.New(inst)
	assert.NotNil(t, role)
	assert.NoError(t, role.Start(ctx, []byte{}))
	defer role.Stop()
}

func TestRoleStartEmptyConfig(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("debug", ctx)
	role := debug.New(inst)
	assert.NotNil(t, role)
	assert.NoError(t, role.Start(ctx, []byte("{}")))
	defer role.Stop()
}

func TestRoleIndex(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ctx := tests.Context()
	inst := rootInst.ForRole("debug", ctx)
	role := debug.New(inst)
	assert.NotNil(t, role)
	assert.NoError(t, role.Start(ctx, []byte("{}")))
	defer role.Stop()

	rr := httptest.NewRecorder()
	req, _ := http.NewRequest(http.MethodGet, "/", nil)
	role.Mux().ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
}
