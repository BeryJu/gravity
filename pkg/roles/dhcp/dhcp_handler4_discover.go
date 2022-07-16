package dhcp

import (
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

func (r *Role) handleDHCPDiscover4(peer net.Addr, m *dhcpv4.DHCPv4) *dhcpv4.DHCPv4 {
	match := r.findLease(m)
	if match == nil {
		scope := r.findScopeForRequest(peer, m)
		if scope == nil {
			r.log.Info("no scope found")
			return nil
		}
		r.log.WithField("scope", scope.Name).Debug("found scope for new lease")
		match = scope.createLeaseFor(peer, m)
		match.put(int64(r.cfg.LeaseNegotiateTimeout))
	} else {
		match.put(match.scope.TTL)
	}

	dhcpRequests.WithLabelValues(m.MessageType().String(), match.scope.Name).Inc()

	rep := match.createReply(peer, m)
	rep.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeOffer))
	return rep
}
