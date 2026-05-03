package log_iml_test

import (
	"testing"
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/extconfig/log_iml"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestInMemoryLogger(t *testing.T) {
	iml := log_iml.Get()
	iml.Flush()
	extconfig.Get().Logger().Warn("test")
	msgs := iml.Messages()
	assert.Len(t, msgs, 1)
}

func TestInMemoryLogger_FieldsToMap(t *testing.T) {
	iml := log_iml.Get()
	iml.Flush()
	n := time.Now()
	extconfig.Get().Logger().With(zap.String("logger_level", "foo")).Warn(
		"test",
		zap.Bool("bool", true),
		zap.Int8("int8", 123),
		zap.Uint8("uint8", 123),
		zap.Float32("float32", 123),
		zap.Float64("float64", 123),
		zap.String("string", "string"),
		zap.Time("time", n),
		zap.Any("any", struct{}{}),
	)
	msgs := iml.Messages()
	assert.Len(t, msgs, 1)
	fm := msgs[0].FieldsToMap()
	assert.NotNil(t, fm["bool"])
	assert.NotNil(t, fm["int8"])
	assert.NotNil(t, fm["uint8"])
	assert.NotNil(t, fm["float32"])
	assert.NotNil(t, fm["float64"])
	assert.NotNil(t, fm["string"])
	assert.NotNil(t, fm["time"])
	assert.NotNil(t, fm["any"])
	assert.NotNil(t, fm["logger_level"])
}

func TestInMemoryLogger_Trunc(t *testing.T) {
	iml := log_iml.Get()
	iml.Flush()
	for i := 0; i <= iml.MaxSize(); i++ {
		extconfig.Get().Logger().Warn("test")
	}
	// Log one more message
	extconfig.Get().Logger().Warn("test")
	msgs := iml.Messages()
	assert.Len(t, msgs, iml.MaxSize())
}
