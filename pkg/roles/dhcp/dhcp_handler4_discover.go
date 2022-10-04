package dhcp

import (
	"github.com/insomniacslk/dhcp/dhcpv4"
)

func (r *Role) HandleDHCPDiscover4(req *Request) *dhcpv4.DHCPv4 {
	match := r.FindLease(req)
	if match == nil {
		scope := r.findScopeForRequest(req)
		if scope == nil {
			req.log.Info("no scope found")
			return nil
		}
		req.log.WithField("scope", scope.Name).Debug("found scope for new lease")
		match = scope.createLeaseFor(req)
		if match == nil {
			return nil
		}
		match.put(req.Context, int64(r.cfg.LeaseNegotiateTimeout))
	} else {
		go match.put(req.Context, match.scope.TTL)
	}

	dhcpRequests.WithLabelValues(req.MessageType().String(), match.scope.Name).Inc()

	rep := match.createReply(req)
	rep.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeOffer))
	return rep
}
