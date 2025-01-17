//go:build e2e
// +build e2e

package tests

import (
	"bytes"
	"testing"

	"github.com/efficientgo/e2e"
	"github.com/stretchr/testify/assert"
)

func TestDNS_Single(t *testing.T) {
	RunGravity(t)

	rb := env.Runnable("dns-tester").
		Init(e2e.StartOptions{
			Image: "gravity-testing:dns",
		})
	assert.NoError(t, e2e.StartAndWaitReady(rb))
	b := bytes.NewBuffer([]byte{})
	err := rb.Exec(e2e.NewCommand("dig", "+short", "10.0.0.1.nip.io"), e2e.WithExecOptionStdout(b))
	assert.NoError(t, err)
	assert.Equal(t, "10.0.0.1\n", b.String())
}
