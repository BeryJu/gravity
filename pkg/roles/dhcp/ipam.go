package dhcp

import (
	"net"
	"net/netip"
)

type IPAM interface {
	NextFreeAddress() *netip.Addr
	IsIPFree(netip.Addr) bool
	GetSubnetMask() net.IPMask
}
