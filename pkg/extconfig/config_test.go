package extconfig

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExtconfig(t *testing.T) {
	assert.NotNil(t, Get())
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
