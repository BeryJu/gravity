package dhcp

import (
	"net"
	"net/netip"
	"time"

	"github.com/go-ping/ping"
	log "github.com/sirupsen/logrus"
)

type InternalIPAM struct {
	SubnetCIDR netip.Prefix

	Start netip.Addr
	End   netip.Addr

	log  *log.Entry
	role *Role
}

func NewInternalIPAM(role *Role, cidr string, rangeStart string, rangeEnd string) (*InternalIPAM, error) {
	sub, err := netip.ParsePrefix(cidr)
	if err != nil {
		return nil, err
	}
	start, err := netip.ParseAddr(rangeStart)
	if err != nil {
		return nil, err
	}
	end, err := netip.ParseAddr(rangeEnd)
	if err != nil {
		return nil, err
	}
	return &InternalIPAM{
		SubnetCIDR: sub,
		Start:      start,
		End:        end,
		log:        role.log.WithField("ipam", "internal"),
		role:       role,
	}, nil
}

func (i *InternalIPAM) NextFreeAddress() *netip.Addr {
	initialIp := i.Start
	for {
		i.log.WithField("ip", initialIp.String()).Debug("checking for free ip")
		// Check if IP is in the correct subnet
		if !i.SubnetCIDR.Contains(initialIp) {
			return nil
		}
		if i.IsIPFree(initialIp) {
			return &initialIp
		}
		initialIp = initialIp.Next()
	}
}

func (i *InternalIPAM) IsIPFree(ip netip.Addr) bool {
	// Ip is less than the start of the range
	if i.Start.Compare(ip) == 1 {
		i.log.Trace("discarding because before start")
		return false
	}
	// Ip is more than the end of the range
	if i.End.Compare(ip) == -1 {
		i.log.Trace("discarding because after end")
		return false
	}
	// check for existing leases
	for _, l := range i.role.leases {
		if l.Address == ip.String() {
			i.log.Debug("discarding because existing lease")
			return false
		}
	}
	pinger, err := ping.NewPinger(ip.String())
	if err != nil {
		i.log.WithError(err).Warning("failed to ping IP")
		return true
	}
	pinger.Count = 1
	pinger.Timeout = 1 * time.Second
	err = pinger.Run()
	if err == nil {
		// IP pings, so it's not usable
		return false
	}
	return true
}

func (i *InternalIPAM) GetSubnetMask() net.IPMask {
	_, cidr, err := net.ParseCIDR(i.SubnetCIDR.String())
	if err != nil {
		// This should never happen as the CIDR is validated in the constructor
		panic(err)
	}
	return cidr.Mask
}
