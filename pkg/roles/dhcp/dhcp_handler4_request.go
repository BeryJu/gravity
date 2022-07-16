package dhcp

import (
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

func (r *Role) handleDHCPRequest4(peer net.Addr, m *dhcpv4.DHCPv4) *dhcpv4.DHCPv4 {
	match := r.findLease(m)

	if match == nil {
		scope := r.findScopeForRequest(peer, m)
		if scope == nil {
			return nil
		}
		r.log.WithField("scope", scope.Name).Debug("found scope for new lease")
		match = scope.createLeaseFor(peer, m)
	}

	match.put(match.scope.TTL)

	dhcpRequests.WithLabelValues(m.MessageType().String(), match.scope.Name).Inc()

	rep := match.createReply(peer, m)
	rep.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeAck))
	return rep
}

func (r *Role) findLease(m *dhcpv4.DHCPv4) *Lease {
	lease, ok := r.leases[r.DeviceIdentifier(m)]
	if !ok {
		return nil
	}
	return lease
}
