package dhcp

import (
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

func (r *DHCPRole) handleDHCPDiscover4(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	match := r.findLease(m)
	if match == nil {
		r.log.Debug("no lease found, creating new")
		scope := r.findScopeForRequest(conn, peer, m)
		if scope == nil {
			r.log.Warning("no scope found")
			return
		}
		r.log.Debug("found scope for new lease")
		match = scope.createLeaseFor(conn, peer, m)
	}
	match.put(int64(r.cfg.LeaseNegotiateTimeout))

	match.reply(conn, peer, m, func(d *dhcpv4.DHCPv4) *dhcpv4.DHCPv4 {
		d.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeOffer))
		return d
	})
}
