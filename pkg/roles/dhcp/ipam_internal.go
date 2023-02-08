package dhcp

import (
	"net"
	"net/netip"
	"strconv"
	"sync"
	"time"

	"github.com/go-ping/ping"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

const InternalIPAMType = "internal"

type InternalIPAM struct {
	SubnetCIDR netip.Prefix
	Start      netip.Addr
	End        netip.Addr

	shouldPing bool
	ips        map[string]bool
	ipsy       sync.RWMutex

	log   *zap.Logger
	role  *Role
	scope *Scope
}

func NewInternalIPAM(role *Role, s *Scope) (*InternalIPAM, error) {
	sub, err := netip.ParsePrefix(s.SubnetCIDR)
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse scope cidr")
	}
	start, err := netip.ParseAddr(s.IPAM["range_start"])
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse 'range_start'")
	}
	end, err := netip.ParseAddr(s.IPAM["range_end"])
	if err != nil {
		return nil, errors.Wrap(err, "failed to parse 'range_end'")
	}
	ipam := &InternalIPAM{
		SubnetCIDR: sub,
		Start:      start,
		End:        end,
		log:        role.log.With(zap.String("ipam", "internal")),
		role:       role,
		scope:      s,
		ips:        make(map[string]bool),
		ipsy:       sync.RWMutex{},
	}
	sp := s.IPAM["should_ping"]
	if sp != "" {
		shouldPing, err := strconv.ParseBool(sp)
		if err != nil {
			return nil, err
		}
		ipam.shouldPing = shouldPing
	}
	return ipam, nil
}

func (i *InternalIPAM) NextFreeAddress() *netip.Addr {
	initialIp := i.Start
	for {
		// Since we start checking at the beginning of the range, check in the loop if we've
		// hit the end and just give up, as the range is full
		if i.End.Compare(initialIp) == -1 {
			break
		}
		i.log.Debug("checking for free ip", zap.String("ip", initialIp.String()))
		// Check if IP is in the correct subnet
		if !i.SubnetCIDR.Contains(initialIp) {
			return nil
		}
		if i.IsIPFree(initialIp) {
			return &initialIp
		} else {
			i.ipsy.Lock()
			i.ips[initialIp.String()] = true
			i.ipsy.Unlock()
		}
		initialIp = initialIp.Next()
	}
	return nil
}

func (i *InternalIPAM) IsIPFree(ip netip.Addr) bool {
	if i.ips[ip.String()] {
		i.log.Debug("discarding", zap.String("reason", "used (in memory)"))
		return false
	}
	// Ip is less than the start of the range
	if i.Start.Compare(ip) == 1 {
		i.log.Debug("discarding", zap.String("reason", "before started"))
		return false
	}
	// Ip is more than the end of the range
	if i.End.Compare(ip) == -1 {
		i.log.Debug("discarding", zap.String("reason", "after end"))
		return false
	}
	// check for existing leases
	for _, l := range i.role.leases {
		// Ignore leases from other scopes
		if l.ScopeKey != i.scope.Name {
			continue
		}
		if l.Address == ip.String() {
			i.log.Debug("discarding", zap.String("reason", "existing lease"))
			return false
		}
	}
	// By default, we dont try to ping the "free" IP
	// and at this point we're confident that it's free
	if !i.shouldPing {
		return true
	}
	pinger, err := ping.NewPinger(ip.String())
	if err != nil {
		i.log.Warn("failed to ping IP", zap.Error(err))
		return true
	}
	pinger.Count = 1
	pinger.Timeout = 1 * time.Second
	pings := false
	pinger.OnRecv = func(pkt *ping.Packet) {
		i.log.Debug("discarding", zap.String("reason", "pings"))
		pings = true
	}
	err = pinger.Run()
	if err != nil {
		i.log.Debug("failed to ping ip", zap.Error(err))
		return false
	}
	return !pings
}

func (i *InternalIPAM) GetSubnetMask() net.IPMask {
	_, cidr, err := net.ParseCIDR(i.SubnetCIDR.String())
	if err != nil {
		// This should never happen as the CIDR is validated in the constructor
		panic(err)
	}
	return cidr.Mask
}
