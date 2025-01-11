//go:build !linux
// +build !linux

package dhcp

import (
	"errors"
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

func (h *handler4) sendEthernet(_ net.Interface, _ *dhcpv4.DHCPv4) error {
	return errors.New("sendEthernet not supported on current platform")
}
