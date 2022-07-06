package dhcp

import (
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

func (r *DHCPRole) handleDHCPRelease4(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {

}
