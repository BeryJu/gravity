package cmd_test

import (
	"encoding/json"
	"testing"

	"beryju.io/gravity/cmd"
	"github.com/stretchr/testify/assert"
)

func TestGenerateSchema(t *testing.T) {
	called := false
	cmd.GenerateSchema("json", func(schema []byte) {
		assert.NotEqual(t, "", string(schema))
		var out interface{}
		json.Unmarshal(schema, &out)
		assert.NotNil(t, out)
		called = true
	})
	assert.True(t, called)
}
