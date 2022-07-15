package dhcp

import (
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

func (r *Role) handleDHCPRequest4(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	match := r.findLease(m)

	if match == nil {
		scope := r.findScopeForRequest(conn, peer, m)
		if scope == nil {
			return
		}
		r.log.Debug("found scope for new lease")
		match = scope.createLeaseFor(conn, peer, m)
	}

	// Run the update in a go-routine since etcd might not be reachable and
	// we don't want to timeout
	go match.put(match.scope.TTL)

	dhcpRequests.WithLabelValues(m.MessageType().String(), match.scope.Name).Inc()

	match.reply(conn, peer, m, func(d *dhcpv4.DHCPv4) *dhcpv4.DHCPv4 {
		d.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeAck))
		return d
	})
}

func (r *Role) findLease(m *dhcpv4.DHCPv4) *Lease {
	lease, ok := r.leases[r.DeviceIdentifier(m)]
	if !ok {
		return nil
	}
	return lease
}
