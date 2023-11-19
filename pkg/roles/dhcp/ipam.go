package dhcp

import (
	"net"
	"net/netip"
)

type IPAM interface {
	// Return the next free IP address, could be sequential or random
	NextFreeAddress() *netip.Addr
	// Check if IP address is in usage (should also check if IP is in specified range and subnet)
	// Can optionally also check if the IP address is pingable
	IsIPFree(netip.Addr) bool
	// Get the subnet mask of the scope
	GetSubnetMask() net.IPMask
}
