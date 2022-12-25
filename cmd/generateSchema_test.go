package cmd_test

import (
	"encoding/json"
	"testing"

	"beryju.io/gravity/cmd"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
)

func TestGenerateSchema(t *testing.T) {
	called := false
	cmd.GenerateSchema(tests.Context(), "json", func(schema []byte) {
		assert.NotEqual(t, "", string(schema))
		var out interface{}
		json.Unmarshal(schema, &out)
		assert.NotNil(t, out)
		called = true
	})
	assert.True(t, called)
}
