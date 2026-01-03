package server_test

import (
	"encoding/json"
	"testing"

	"beryju.io/gravity/cmd/server"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestGenerateSchema(t *testing.T) {
	tests.Setup(t)
	called := false
	server.GenerateSchema(tests.Context(), "json", func(schema []byte) {
		assert.NotEqual(t, "", string(schema))
		var out interface{}
		err := json.Unmarshal(schema, &out)
		assert.NoError(t, err)
		assert.NotNil(t, out)
		called = true
	})
	assert.True(t, called)
}
