package dhcp

import (
	"github.com/insomniacslk/dhcp/dhcpv4"
	"go.uber.org/zap"
)

func (r *Role) HandleDHCPRelease4(req *Request4) *dhcpv4.DHCPv4 {
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
		return nil
	}
	if !match.IsReservation() {
		err := match.Delete(req.Context)
		if err != nil {
			req.log.Warn("failed to put lease", zap.Error(err))
		} else {
			req.log.Info("deleted lease from release")
		}
	}
	dhcpRequests.WithLabelValues(req.MessageType().String(), match.scope.Name).Inc()
	return nil
}
