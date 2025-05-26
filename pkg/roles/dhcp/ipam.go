package dhcp

import (
	"math/big"
	"net"
	"net/netip"
)

type IPAM interface {
	// Return the next free IP address, could be sequential or random
	NextFreeAddress(identifier string) *netip.Addr
	// Check if IP address is in usage (should also check if IP is in specified range and subnet)
	// Can optionally also check if the IP address is pingable
	// `identifier` might be given as well for a device that could request an address
	// that it had already taken
	IsIPFree(addr netip.Addr, identifier *string) bool
	// Get the subnet mask of the scope
	GetSubnetMask() net.IPMask
	// Mark an IP as used
	UseIP(addr netip.Addr, identifier string)
	// Mark an IP as unused
	FreeIP(ip netip.Addr)
	// Update config when scope is updated
	UpdateConfig(s *Scope) error
	// Amount of usable IPs excluding any exclusions
	UsableSize() *big.Int
}
