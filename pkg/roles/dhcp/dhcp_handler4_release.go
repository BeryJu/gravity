package dhcp

import (
	"net/netip"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"go.uber.org/zap"
)

func (r *Role) HandleDHCPRelease4(req *Request4) *dhcpv4.DHCPv4 {
	match := r.FindLease(req)
	if match == nil || match.IsReservation() {
		return nil
	}
	ip, err := netip.ParseAddr(match.Address)
	if err != nil {
		req.log.Warn("failed to parse address from lease", zap.Error(err))
	} else {
		match.scope.ipam.FreeIP(ip)
	}
	err = match.Delete(req.Context)
	if err != nil {
		req.log.Warn("failed to delete lease", zap.Error(err))
	} else {
		req.log.Info("deleted lease from release")
	}
	dhcpRequests.WithLabelValues(req.MessageType().String(), match.scope.Name).Inc()
	return nil
}
