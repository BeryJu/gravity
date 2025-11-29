//go:build e2e

package tests

import (
	"context"
	"io"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/exec"
)

func Context(t *testing.T) context.Context {
	ctx, cn := context.WithCancel(context.Background())
	t.Cleanup(func() {
		cn()
	})
	return ctx
}

func ExecCommand(t *testing.T, co testcontainers.Container, cmd []string, options ...exec.ProcessOption) (int, string) {
	ctx := Context(t)
	options = append(options, exec.Multiplexed())
	t.Logf("Running command '%s'...", strings.Join(cmd, " "))
	c, out, err := co.Exec(ctx, cmd, options...)
	assert.NoError(t, err)
	t.Logf("Error code: %d", c)
	body, err := io.ReadAll(out)
	assert.NoError(t, err)
	t.Logf("Command output: '%s'", string(body))
	return c, string(body)
}

func MustExec(t *testing.T, co testcontainers.Container, cmd string) string {
	rc, b := ExecCommand(t, co, []string{"bash", "-c", cmd})
	assert.Equal(t, 0, rc, b)
	return b
}
