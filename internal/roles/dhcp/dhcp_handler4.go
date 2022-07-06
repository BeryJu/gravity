package dhcp

import (
	"fmt"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

func (r *DHCPServerRole) handler4(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	r.log.WithField("dhcp-request", m.Summary()).Trace("DHCP Request")
	switch mt := m.MessageType(); mt {
	case dhcpv4.MessageTypeDiscover:
		r.handleDHCPDiscover4(conn, peer, m)
	case dhcpv4.MessageTypeRequest:
		r.handleDHCPRequest4(conn, peer, m)
	case dhcpv4.MessageTypeRelease:
		r.handleDHCPRelease4(conn, peer, m)
	default:
		r.log.WithField("type", mt).Info("unsupported message type")
		return
	}
}

func (r *DHCPServerRole) handleDHCPDiscover4(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	fmt.Println(m.String())
	// match := r.findLease(m)
	// found := match != nil

	// if !found {
	// 	r.l.V(1).Info("no lease found, creating new")
	// 	scope := r.findScopeForRequest(conn, peer, m)
	// 	if scope == nil {
	// 		return
	// 	}
	// 	r.l.V(1).Info("found scope for new lease")
	// 	match = r.createLeaseFor(scope, conn, peer, m)
	// 	// Don't save the new lease yet as the client hasn't committed to it yet
	// }
	// // If we already have a saved lease, update its status
	// if found {
	// 	match.Status.LastRequest = time.Now().Format(time.RFC3339)
	// 	err := r.Update(context.Background(), match)
	// 	if err != nil {
	// 		r.l.Error(err, "failed to update lease")
	// 	}
	// }
	// r.replyWithLease(match, conn, peer, m, func(d *dhcpv4.DHCPv4) *dhcpv4.DHCPv4 {
	// 	d.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeOffer))
	// 	return d
	// })

}

func (r *DHCPServerRole) handleDHCPRequest4(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {

}

func (r *DHCPServerRole) handleDHCPRelease4(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {

}
