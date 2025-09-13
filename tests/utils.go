//go:build e2e
// +build e2e

package tests

import (
	"context"
	"io"
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
	c, out, err := co.Exec(ctx, cmd, options...)
	assert.NoError(t, err)
	body, err := io.ReadAll(out)
	assert.NoError(t, err)
	return c, string(body)
}
