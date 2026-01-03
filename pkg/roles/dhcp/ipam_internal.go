package dhcp

import (
	"math/big"
	"net"
	"net/netip"
	"strconv"
	"sync"
	"time"

	"github.com/netdata/go.d.plugin/pkg/iprange"
	"github.com/pkg/errors"
	probing "github.com/prometheus-community/pro-bing"
	"go.uber.org/zap"
)

const InternalIPAMType = "internal"

type InternalIPAM struct {
	Start netip.Addr
	End   netip.Addr

	ipf        sync.Mutex
	log        *zap.Logger
	role       *Role
	scope      *Scope
	SubnetCIDR netip.Prefix

	shouldPing bool
}

func NewInternalIPAM(role *Role, s *Scope) (*InternalIPAM, error) {
	ipam := &InternalIPAM{
		log:   role.log.With(zap.String("ipam", "internal")),
		role:  role,
		scope: s,
		ipf:   sync.Mutex{},
	}
	err := ipam.UpdateConfig(s)
	if err != nil {
		return nil, err
	}
	return ipam, nil
}

func (i *InternalIPAM) UpdateConfig(s *Scope) error {
	sub, err := netip.ParsePrefix(s.SubnetCIDR)
	if err != nil {
		return errors.Wrap(err, "failed to parse scope cidr")
	}
	start, err := netip.ParseAddr(s.IPAM["range_start"])
	if err != nil {
		return errors.Wrap(err, "failed to parse 'range_start'")
	}
	end, err := netip.ParseAddr(s.IPAM["range_end"])
	if err != nil {
		return errors.Wrap(err, "failed to parse 'range_end'")
	}
	i.SubnetCIDR = sub
	i.Start = start
	i.End = end
	sp := s.IPAM["should_ping"]
	if sp != "" {
		shouldPing, err := strconv.ParseBool(sp)
		if err != nil {
			return err
		}
		i.shouldPing = shouldPing
	}
	return nil
}

// Return the next free IP in the range defined by `.Start` and `.End`, inclusive.
// Any returned address is _not_ marked as used, as this is up to the caller of the function
// Might return `nil` if no more free IP Address is available
func (i *InternalIPAM) NextFreeAddress(identifier string) *netip.Addr {
	i.ipf.Lock()
	defer i.ipf.Unlock()
	currentIP := i.Start
	// Since we start checking at the beginning of the range, check in the loop if we've
	// hit the end and just give up, as the range is full
	for i.End.Compare(currentIP) != -1 {
		i.log.Debug("checking for free IP", zap.String("ip", currentIP.String()))
		// Check if IP is in the correct subnet
		if !i.SubnetCIDR.Contains(currentIP) {
			break
		}
		if i.IsIPFree(currentIP, &identifier) {
			// Free IP is returned, _not_ marked as used, this the responsibility of the caller
			// to mark the IP as used
			return &currentIP
		}
		currentIP = currentIP.Next()
	}
	i.log.Warn("no more empty IPs left", zap.String("lastIp", currentIP.String()))
	return nil
}

func (i *InternalIPAM) FreeIP(ip netip.Addr) {
}

func (i *InternalIPAM) UseIP(ip netip.Addr, identifier string) {
}

func (i *InternalIPAM) IsIPFree(ip netip.Addr, identifier *string) bool {
	if identifier != nil {
		l := i.role.leases.Get(*identifier)
		if l != nil && l.Address == ip.String() {
			i.log.Debug("allowing", zap.String("ip", ip.String()), zap.String("reason", "existing IP of lease"))
			return true
		}
	}
	for _, l := range i.role.leases.Iter() {
		if l.Address == ip.String() {
			i.log.Debug("discarding", zap.String("ip", ip.String()), zap.String("reason", "used (in memory)"))
			return false
		}
	}
	// IP is less than the start of the range
	if i.Start.Compare(ip) == 1 {
		i.log.Debug("discarding", zap.String("ip", ip.String()), zap.String("reason", "before started"))
		return false
	}
	// IP is more than the end of the range
	if i.End.Compare(ip) == -1 {
		i.log.Debug("discarding", zap.String("ip", ip.String()), zap.String("reason", "after end"))
		return false
	}
	// check for existing leases
	for _, l := range i.role.leases.Iter() {
		// Ignore leases from other scopes
		if l.ScopeKey != i.scope.Name {
			continue
		}
		if l.Address != ip.String() {
			continue
		}
		if identifier != nil && l.Identifier == *identifier {
			i.UseIP(ip, *identifier)
			i.log.Debug("allowing", zap.String("ip", ip.String()), zap.String("reason", "existing matching lease"))
			return true
		}
		i.log.Debug("discarding", zap.String("ip", ip.String()), zap.String("reason", "existing lease"))
		return false
	}
	// check for ping
	if i.shouldPing && i.ping(ip) {
		return false
	}
	i.log.Debug("allowing", zap.String("ip", ip.String()), zap.String("reason", "free"))
	return true
}

// Attempt to ping the IP. If we can ping the IP successfully, return true, and if we fail
// for any reason return false
func (i *InternalIPAM) ping(ip netip.Addr) bool {
	pinger, err := probing.NewPinger(ip.String())
	if err != nil {
		i.log.Warn("failed to ping IP", zap.Error(err))
		return false
	}
	pinger.Count = 1
	pinger.Timeout = 1 * time.Second
	pings := false
	pinger.OnRecv = func(pkt *probing.Packet) {
		i.log.Debug("discarding", zap.String("ip", ip.String()), zap.String("reason", "pings"))
		pings = false
	}
	err = pinger.Run()
	if err != nil {
		i.log.Debug("failed to ping ip", zap.Error(err))
		return false
	}
	return pings
}

func (i *InternalIPAM) GetSubnetMask() net.IPMask {
	_, cidr, err := net.ParseCIDR(i.SubnetCIDR.String())
	if err != nil {
		// This should never happen as the CIDR is validated in the constructor
		panic(err)
	}
	return cidr.Mask
}

func (i *InternalIPAM) UsableSize() *big.Int {
	ips := iprange.New(i.Start.AsSlice(), i.End.AsSlice())
	return ips.Size()
}
