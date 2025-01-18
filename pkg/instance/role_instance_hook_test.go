package instance_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestHook(t *testing.T) {
	defer tests.Setup(t)()
	rootInst := instance.New()
	ri := rootInst.ForRole("test", tests.Context())
	v := ri.ExecuteHook(roles.HookOptions{
		Source: `function test() {
			gravity.log("test");
			return true;
		}`,
		Method: "test",
	})
	assert.True(t, v.(bool))
}
