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

type IPUse struct {
	identifier string
	unknown    bool
}

type InternalIPAM struct {
	Start netip.Addr
	End   netip.Addr

	ips  map[string]*IPUse
	ipsm sync.RWMutex
	ipf  sync.Mutex

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
		ips:   make(map[string]*IPUse),
		ipsm:  sync.RWMutex{},
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

func (i *InternalIPAM) NextFreeAddress(identifier string) *netip.Addr {
	currentIP := i.Start
	for {
		// Since we start checking at the beginning of the range, check in the loop if we've
		// hit the end and just give up, as the range is full
		if i.End.Compare(currentIP) == -1 {
			break
		}
		i.log.Debug("checking for free IP", zap.String("ip", currentIP.String()))
		// Check if IP is in the correct subnet
		if !i.SubnetCIDR.Contains(currentIP) {
			break
		}
		if i.IsIPFree(currentIP, &identifier) {
			// Actually mark IP as used
			// i.UseIP(currentIP, identifier)
			return &currentIP
		}
		// As the IP is not free we're marking it as used, however this is a fallback if
		// `IsIPFree` didn't mark the IP as used by a specific lease
		i.useIP(currentIP, IPUse{
			identifier: identifier,
			unknown:    true,
		}, false)
		currentIP = currentIP.Next()
	}
	i.log.Warn("no more empty IPs left", zap.String("lastIp", currentIP.String()))
	return nil
}

func (i *InternalIPAM) UseIP(ip netip.Addr, identifier string) {
	i.useIP(ip, IPUse{
		identifier: identifier,
		unknown:    false,
	}, true)
}

func (i *InternalIPAM) useIP(ip netip.Addr, ipu IPUse, overwrite bool) {
	i.ipsm.Lock()
	defer i.ipsm.Unlock()
	if i.ips[ip.String()] != nil && !overwrite {
		return
	}
	i.ips[ip.String()] = &ipu
}

func (i *InternalIPAM) IsIPFree(ip netip.Addr, identifier *string) bool {
	i.ipf.Lock()
	defer i.ipf.Unlock()
	i.ipsm.RLock()
	mem := i.ips[ip.String()]
	i.ipsm.RUnlock()
	if mem != nil {
		i.log.Debug("discarding", zap.String("ip", ip.String()), zap.String("reason", "used (in memory)"))
		return false
	}
	// Ip is less than the start of the range
	if i.Start.Compare(ip) == 1 {
		i.log.Debug("discarding", zap.String("ip", ip.String()), zap.String("reason", "before started"))
		return false
	}
	// Ip is more than the end of the range
	if i.End.Compare(ip) == -1 {
		i.log.Debug("discarding", zap.String("ip", ip.String()), zap.String("reason", "after end"))
		return false
	}
	// check for existing leases
	for _, l := range i.role.leases {
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
	if identifier != nil {
		i.useIP(ip, IPUse{
			identifier: *identifier,
			unknown:    true,
		}, false)
	}
	i.log.Debug("allowing", zap.String("ip", ip.String()), zap.String("reason", "free"))
	return true
}

// Attempt to ping the IP. If we can ping the IP successfully, return true, and if we fail
// for any reason return false
func (i *InternalIPAM) ping(ip netip.Addr) bool {
	pinger, err := ping.NewPinger(ip.String())
	if err != nil {
		i.log.Warn("failed to ping IP", zap.Error(err))
		return false
	}
	pinger.Count = 1
	pinger.Timeout = 1 * time.Second
	pings := false
	pinger.OnRecv = func(pkt *ping.Packet) {
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
