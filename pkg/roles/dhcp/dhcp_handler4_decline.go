package dhcp

import (
	"github.com/insomniacslk/dhcp/dhcpv4"
	"go.uber.org/zap"
)

func (r *Role) HandleDHCPDecline4(req *Request4) *dhcpv4.DHCPv4 {
	match := r.FindLease(req)
	if match == nil {
		scope := r.findScopeForRequest(req)
		if scope == nil {
			req.log.Info("no scope found")
			return nil
		}
		req.log.Debug("found scope for new lease", zap.String("scope", scope.Name))
		match = scope.leaseFor(req)
		if match == nil {
			return nil
		}
		// because this happens for DHCP decline, the IP is assumed to be already taken
		// and we only get here if there's no lease so the device is assumed to be managed
		// externally, so create a leave with an "invalid" identifier which won't be picked
		match.Identifier = "invalid"
		return nil
	}
	// since there's no further requests to confirm this lease, save it directly with the TTL of the scope
	err := match.Put(req.Context, match.scope.TTL)
	if err != nil {
		r.log.Warn("failed to put lease", zap.Error(err))
	}

	dhcpRequests.WithLabelValues(req.MessageType().String(), match.scope.Name).Inc()
	return nil
}
