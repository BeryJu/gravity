package dhcp

import (
	"beryju.io/gravity/pkg/roles"
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

	err := match.Put(req.Context, match.scope.TTL)
	if err != nil {
		r.log.Warn("failed to put dhcp lease", zap.Error(err))
	}

	dhcpRequests.WithLabelValues(req.MessageType().String(), match.scope.Name).Inc()

	r.i.ExecuteHook(roles.HookOptions{
		Source: match.scope.Hook,
		Method: "onDHCPRequestBefore",
	}, req)
	rep := match.createReply(req)
	rep.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeAck))
	r.i.ExecuteHook(roles.HookOptions{
		Source: match.scope.Hook,
		Method: "onDHCPRequestAfter",
	}, req, rep)
	return rep
}
