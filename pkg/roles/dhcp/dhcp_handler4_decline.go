package dhcp

import (
	"github.com/insomniacslk/dhcp/dhcpv4"
)

func (r *Role) handleDHCPDecline4(req *Request) *dhcpv4.DHCPv4 {
	match := r.findLease(req)
	if match == nil {
		return nil
	}
	_, err := r.i.KV().Delete(req.Context, match.etcdKey)
	if err != nil {
		r.log.WithError(err).Warning("failed to delete lease")
	}

	dhcpRequests.WithLabelValues(req.MessageType().String(), match.scope.Name).Inc()
	return nil
}
