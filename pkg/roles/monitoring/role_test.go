package monitoring_test

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles/api"
	"beryju.io/gravity/pkg/roles/monitoring"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestRoleStartNoConfig(t *testing.T) {
	rootInst := instance.New()
	api := api.New(rootInst.ForRole("api"))
	api.Start(tests.Context(), []byte{})
	defer api.Stop()
	role := monitoring.New(rootInst.ForRole("monitoring"))
	assert.NotNil(t, role)
	ctx := tests.Context()
	assert.Nil(t, role.Start(ctx, []byte{}))
	defer role.Stop()
}

func TestRoleStartEmptyConfig(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("monitoring")
	role := monitoring.New(inst)
	assert.NotNil(t, role)
	ctx := tests.Context()
	assert.Nil(t, role.Start(ctx, []byte("{}")))
	defer role.Stop()
}

func TestRoleStartInvalidListen(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("monitoring")
	role := monitoring.New(inst)
	assert.NotNil(t, role)
	ctx := tests.Context()
	assert.Nil(t, role.Start(ctx, []byte(tests.MustJSON(monitoring.RoleConfig{
		Port: -1,
	}))))
	time.Sleep(1 * time.Second)
	assert.False(t, role.IsRunning())
	defer role.Stop()
}

func TestRoleHealth(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("monitoring")
	role := monitoring.New(inst)
	assert.NotNil(t, role)
	ctx := tests.Context()
	assert.Nil(t, role.Start(ctx, []byte("{}")))
	defer role.Stop()
	rr := httptest.NewRecorder()
	role.HandleHealthLive(rr, httptest.NewRequest("GET", "/", nil))
	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
}

func TestMetricsScrape(t *testing.T) {
	rootInst := instance.New()
	inst := rootInst.ForRole("monitoring")
	role := monitoring.New(inst)
	assert.NotNil(t, role)
	ctx := tests.Context()
	assert.Nil(t, role.Start(ctx, []byte("{}")))
	defer role.Stop()
	rr := httptest.NewRecorder()
	role.HandleMetrics(rr, httptest.NewRequest("GET", "/", nil))
	assert.Equal(t, http.StatusOK, rr.Result().StatusCode)
}
