package storage

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKey(t *testing.T) {
	c := NewClient("/gravity", nil, false, "localhost:2379")
	k := c.Key().Add("foo", "bar")
	assert.Equal(t, "/foo/bar", k.String())
	k = c.Key().Add("foo", "bar").Prefix(true)
	assert.Equal(t, "/foo/bar/", k.String())
}
