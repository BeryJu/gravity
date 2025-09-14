//go:build e2e
// +build e2e

package tests

import (
	"testing"

	"beryju.io/gravity/tests/gravity"
	"github.com/stretchr/testify/assert"
)

func TestCLI_Health(t *testing.T) {
	gr := gravity.New(t)

	_, health := ExecCommand(t, gr.Container(), []string{"gravity", "cli", "health"})
	assert.Contains(t, string(health), "gravity-1")
}
