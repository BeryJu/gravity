package tests

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCLI_Health(t *testing.T) {
	gr := RunGravity(t, WithEnv("LOG_LEVEL", "warn"))

	_, health := ExecCommand(t, gr.container, []string{"gravity", "cli", "health"})
	assert.Contains(t, string(health), "gravity-1")
}
