package dhcp

import (
	"github.com/insomniacslk/dhcp/dhcpv4"
	"go.uber.org/zap"
)

func (r *Role) HandleDHCPDiscover4(req *Request4) *dhcpv4.DHCPv4 {
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
		err := match.Put(req.Context, int64(r.cfg.LeaseNegotiateTimeout))
		if err != nil {
			req.log.Warn("failed to update lease during discover creation", zap.Error(err))
		}
	} else {
		go func() {
			err := match.Put(req.Context, match.scope.TTL)
			if err != nil {
				req.log.Warn("failed to update lease during discover", zap.Error(err))
			}
		}()
	}

	dhcpRequests.WithLabelValues(req.MessageType().String(), match.scope.Name).Inc()

	match.scope.executeHook("onDHCPRequestBefore", req)
	rep := match.createReply(req)
	rep.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeOffer))
	match.scope.executeHook("onDHCPRequestAfter", req, rep)
	return rep
}
