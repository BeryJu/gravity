package storage

import (
	"testing"

	"beryju.io/gravity/pkg/extconfig"
	"github.com/stretchr/testify/assert"
)

func TestKey(t *testing.T) {
	c := NewClient(extconfig.Get().Etcd.Prefix, extconfig.Get().Etcd.Endpoint)
	k := c.Key().Add("foo", "bar")
	assert.Equal(t, "/foo/bar", k.String())
	k = c.Key().Add("foo", "bar").Prefix(true)
	assert.Equal(t, "/foo/bar/", k.String())
}
