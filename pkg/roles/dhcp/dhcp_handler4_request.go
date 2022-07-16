package dhcp

import (
	"github.com/insomniacslk/dhcp/dhcpv4"
)

func (r *Role) handleDHCPRequest4(req *Request) *dhcpv4.DHCPv4 {
	match := r.findLease(req)

	if match == nil {
		scope := r.findScopeForRequest(req)
		if scope == nil {
			return nil
		}
		req.log.WithField("scope", scope.Name).Debug("found scope for new lease")
		match = scope.createLeaseFor(req)
		if match == nil {
			return nil
		}
	}

	match.put(match.scope.TTL)

	dhcpRequests.WithLabelValues(req.MessageType().String(), match.scope.Name).Inc()

	rep := match.createReply(req)
	rep.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeAck))
	return rep
}

func (r *Role) findLease(req *Request) *Lease {
	lease, ok := r.leases[r.DeviceIdentifier(req.DHCPv4)]
	if !ok {
		return nil
	}
	return lease
}
