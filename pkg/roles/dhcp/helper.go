package dhcp

import (
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

func (r *DHCPRole) findScopeForRequest(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) *Scope {
	var match *Scope
	for _, scope := range r.scopes {
		// TODO: priority and order
		if scope.match(conn, peer, m) {
			r.log.WithField("name", scope.Name).Debug("selected scope based on match")
			match = scope
		}
		if match == nil && scope.Default {
			r.log.WithField("name", scope.Name).Debug("selected scope based on default state")
			match = scope
		}
	}
	if match != nil {
		r.log.Trace("found scope for request")
	}
	return match
}
