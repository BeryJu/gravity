package dhcp

import (
	"net"

	"beryju.io/gravity/pkg/extconfig"
	"github.com/insomniacslk/dhcp/dhcpv4"
)

func (r *DHCPRole) handler4(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	if r.cfg.ListenOnly || extconfig.Get().ListenOnlyMode {
		return
	}

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
