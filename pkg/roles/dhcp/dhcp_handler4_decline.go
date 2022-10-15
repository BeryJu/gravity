package dhcp

import (
	"github.com/insomniacslk/dhcp/dhcpv4"
	"go.uber.org/zap"
)

func (r *Role) HandleDHCPDecline4(req *Request4) *dhcpv4.DHCPv4 {
	match := r.FindLease(req)
	if match == nil {
		return nil
	}
	_, err := r.i.KV().Delete(req.Context, match.etcdKey)
	if err != nil {
		r.log.Warn("failed to delete lease", zap.Error(err))
	}

	dhcpRequests.WithLabelValues(req.MessageType().String(), match.scope.Name).Inc()
	return nil
}
