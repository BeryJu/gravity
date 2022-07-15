package dhcp

import (
	"net"
	"net/netip"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	s := Scope{
		cidr: netip.MustParsePrefix("192.0.2.0/24"),
	}
	assert.Equal(t, s.match(nil, net.ParseIP("192.0.2.2"), nil), 24)
	assert.Equal(t, s.match(nil, net.ParseIP("192.3.2.2"), nil), -1)
}
