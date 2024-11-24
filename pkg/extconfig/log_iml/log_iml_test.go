package log_iml_test

import (
	"testing"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/extconfig/log_iml"
	"github.com/stretchr/testify/assert"
)

func TestInMemoryLogger(t *testing.T) {
	iml := log_iml.Get()
	iml.Flush()
	extconfig.Get().Logger().Debug("test")
	msgs := iml.Messages()
	assert.Len(t, msgs, 1)
}

func TestInMemoryLogger_Trunc(t *testing.T) {
	iml := log_iml.Get()
	iml.Flush()
	for i := 0; i <= iml.MaxSize(); i++ {
		extconfig.Get().Logger().Debug("test")
	}
	// Log one more message
	extconfig.Get().Logger().Debug("test")
	msgs := iml.Messages()
	assert.Len(t, msgs, iml.MaxSize())
}
