package instance_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func Test_Events(t *testing.T) {
	rootInst := instance.New()
	called := false
	rootInst.ForRole("test").AddEventListener("test-topic", func(ev *roles.Event) {
		called = true
	})
	rootInst.DispatchEvent("test-topic", roles.NewEvent(tests.Context(), map[string]interface{}{}))
	assert.True(t, called)
}
