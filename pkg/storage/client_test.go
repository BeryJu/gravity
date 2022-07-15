package storage

import (
	"testing"

	"beryju.io/gravity/pkg/extconfig"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	c := NewClient(extconfig.Get().Etcd.Prefix, extconfig.Get().Etcd.Endpoint)
	assert.NotNil(t, c)
	assert.Panics(t, func() {
		NewClient(extconfig.Get().Etcd.Prefix)
	})
}
