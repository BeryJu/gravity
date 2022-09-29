package extconfig

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtconfig(t *testing.T) {
	assert.NotNil(t, Get())
	assert.NotNil(t, Get().Dirs())
	assert.NotNil(t, Get().EtcdClient())
	assert.True(t, strings.HasSuffix(Get().Listen(1234), "1234"))
}

func TestResolver(t *testing.T) {
	assert.NotNil(t, Resolver())
}

func TestTransport(t *testing.T) {
	assert.NotNil(t, Transport())
}

func TestGetIP(t *testing.T) {
	ip, err := GetIP()
	assert.NotNil(t, ip)
	assert.Nil(t, err)
}

func TestVersion(t *testing.T) {
	assert.Equal(t, Version+"-dev", FullVersion())
	BuildHash = "foo"
	assert.Equal(t, Version+"-foo", FullVersion())
	BuildHash = "foobqerqewrqwer"
	assert.Equal(t, Version+"-foobqerq", FullVersion())
}
