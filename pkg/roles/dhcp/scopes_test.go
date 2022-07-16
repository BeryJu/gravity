package dhcp

import (
	"net"
	"net/netip"
	"testing"

	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestMatch(t *testing.T) {
	s := Scope{
		cidr: netip.MustParsePrefix("192.0.2.0/24"),
	}
	req := &Request{
		log: log.WithField("foo", "bar"),
	}
	assert.Equal(t, s.match(net.ParseIP("192.0.2.2"), req), 24)
	assert.Equal(t, s.match(net.ParseIP("192.3.2.2"), req), -1)
}
