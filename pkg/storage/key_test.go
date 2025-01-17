package storage_test

import (
	"testing"

	"beryju.io/gravity/pkg/storage"
	"github.com/stretchr/testify/assert"
)

func TestKey(t *testing.T) {
	c := storage.NewClient("/gravity", nil, false, "localhost:2379")
	k := c.Key().Add("foo", "bar")
	assert.Equal(t, "/foo/bar", k.String())
	k = c.Key().Add("foo", "bar").Prefix(true)
	assert.Equal(t, "/foo/bar/", k.String())
	k = c.Key().Add("foo", "bar").Prefix(true).Up()
	assert.Equal(t, "/foo/", k.String())
}

func TestKeyCopy(t *testing.T) {
	c := storage.NewClient("/gravity", nil, false, "localhost:2379")
	k := c.Key().Add("foo", "bar")
	assert.Equal(t, "/foo/bar", k.String())
	assert.Equal(t, "/foo/bar", k.Copy().String())
}

func TestKeyParse(t *testing.T) {
	assert.Equal(t, "/foo/bar", storage.KeyFromString("/foo/bar").String())
	assert.True(t, storage.KeyFromString("/foo/bar/").IsPrefix())
}
