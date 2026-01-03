package instance_test

import (
	"testing"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestHook(t *testing.T) {
	tests.Setup(t)
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

func TestHook_ParseIP(t *testing.T) {
	tests.Setup(t)
	rootInst := instance.New()
	ri := rootInst.ForRole("test", tests.Context())
	v := ri.ExecuteHook(roles.HookOptions{
		Source: `function test() {
			return net.parseIP("192.168.1.100", "v4");
		}`,
		Method: "test",
	})
	assert.Equal(t, []byte{0xc0, 0xa8, 0x1, 0x64}, v.([]byte))
}
