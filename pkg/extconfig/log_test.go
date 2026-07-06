package extconfig

import (
	"testing"

	"beryju.io/gravity/pkg/extconfig/log_iml"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func TestLogLevel(t *testing.T) {
	c := &ExtConfig{
		LogLevel: "debug",
	}
	assert.Equal(t, c.LogLevelFor("root"), zapcore.DebugLevel)
}

func TestLogLevel_Invalid(t *testing.T) {
	c := &ExtConfig{
		LogLevel: "something",
	}
	assert.Equal(t, c.LogLevelFor("root"), zapcore.InfoLevel)
}

func TestLogLevel_Role(t *testing.T) {
	c := &ExtConfig{
		LogLevel: "debug,test-role=warn",
	}
	assert.Equal(t, c.LogLevelFor("root"), zapcore.DebugLevel)
	assert.Equal(t, c.LogLevelFor("test-role"), zapcore.WarnLevel)
}

func TestLogLevel_Role_Full(t *testing.T) {
	c := &ExtConfig{
		LogLevel: "debug,test-role=warn",
	}
	c.logger = c.BuildLogger()

	roleId := "test-role"
	roleLogger := c.Logger().Named("role." + roleId).WithOptions(zap.IncreaseLevel(c.LogLevelFor(roleId)))

	assert.Equal(t, c.LogLevelFor("root"), zapcore.DebugLevel)
	assert.Equal(t, roleLogger.Level(), zapcore.WarnLevel)
}

// When LOG_LEVEL=error and no role-specific override exists, the role should
// inherit the root level (error), not fall back to info.
func TestLogLevel_Role_NoOverride(t *testing.T) {
	c := &ExtConfig{
		LogLevel: "error",
	}
	assert.Equal(t, zapcore.ErrorLevel, c.LogLevelFor("root"))
	assert.Equal(t, zapcore.ErrorLevel, c.LogLevelFor("etcd"))
	assert.Equal(t, zapcore.ErrorLevel, c.LogLevelFor("dns"))
}

func TestLogLevel_Role_Full_Decrease(t *testing.T) {
	c := &ExtConfig{
		LogLevel: "warn,test-role=debug",
	}
	c.logger = c.BuildLogger()

	roleId := "test-role"
	roleLogger := c.Logger().Named("role." + roleId).WithOptions(SetLevel(c.LogLevelFor(roleId))).With(zap.String("foo", "bar"))

	roleLogger.Debug("foo")

	found := false
	for _, msg := range log_iml.Get().Messages() {
		if msg.Entry.Message == "foo" && msg.Entry.Level == zap.DebugLevel {
			found = true
		}
	}
	assert.True(t, found)
	assert.Equal(t, c.LogLevelFor("root"), zapcore.WarnLevel)
	assert.Equal(t, roleLogger.Level(), zapcore.DebugLevel)
}
