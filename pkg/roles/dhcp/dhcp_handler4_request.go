package dhcp

import (
	"github.com/insomniacslk/dhcp/dhcpv4"
	"go.uber.org/zap"
)

func (r *Role) HandleDHCPRequest4(req *Request4) *dhcpv4.DHCPv4 {
	match := r.FindLease(req)

	if match == nil {
		scope := r.findScopeForRequest(req)
		if scope == nil {
			return nil
		}
		req.log.Debug("found scope for new lease", zap.String("scope", scope.Name))
		match = scope.createLeaseFor(req)
		if match == nil {
			return nil
		}
	}

	match.Put(req.Context, match.scope.TTL)

	dhcpRequests.WithLabelValues(req.MessageType().String(), match.scope.Name).Inc()

	rep := match.createReply(req)
	rep.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeAck))
	return rep
}

func (r *Role) FindLease(req *Request4) *Lease {
	r.leasesM.RLock()
	defer r.leasesM.RUnlock()
	lease, ok := r.leases[r.DeviceIdentifier(req.DHCPv4)]
	if !ok {
		return nil
	}
	return lease
}
