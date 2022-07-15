package dhcp

import (
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

func (r *Role) handleDHCPDiscover4(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	match := r.findLease(m)
	if match == nil {
		scope := r.findScopeForRequest(conn, peer, m)
		if scope == nil {
			r.log.Info("no scope found")
			return
		}
		r.log.Debug("found scope for new lease")
		match = scope.createLeaseFor(conn, peer, m)
		match.put(int64(r.cfg.LeaseNegotiateTimeout))
	} else {
		go match.put(match.scope.TTL)
	}

	dhcpRequests.WithLabelValues(m.MessageType().String(), match.scope.Name).Inc()

	match.reply(conn, peer, m, func(d *dhcpv4.DHCPv4) *dhcpv4.DHCPv4 {
		d.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeOffer))
		return d
	})
}
