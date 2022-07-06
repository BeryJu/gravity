package ipam

import (
	"net"
	"net/netip"

	log "github.com/sirupsen/logrus"
)

type InternalIPAM struct {
	SubnetCIDR netip.Prefix

	Start netip.Addr
	End   netip.Addr

	log *log.Entry
}

func NewInternalIPAM(cidr string, rangeStart string, rangeEnd string) (*InternalIPAM, error) {
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
		log:        log.WithField("component", "ipam"),
	}, nil
}

func (i *InternalIPAM) NextFreeAddress() *netip.Addr {
	initialIp := i.SubnetCIDR.Addr()
	for {
		initialIp = initialIp.Next()
		i.log.WithField("ip", initialIp.String()).Debug("checking for free ip")
		// Check if IP is in the correct subnet
		if !i.SubnetCIDR.Contains(initialIp) {
			return nil
		}
		if i.IsIPFree(initialIp) {
			return &initialIp
		}
	}
}

func (i *InternalIPAM) IsIPFree(ip netip.Addr) bool {
	// free := false
	// get all leases to check

	// Ip is less than the start of the range
	if i.Start.Compare(ip) == 1 {
		i.log.Debug("discarding because before start")
		return false
	}
	// Ip is more than the end of the range
	if i.End.Compare(ip) == -1 {
		i.log.Debug("discarding because after end")
		return false
	}
	// check for existing leases
	// for _, l := range leasei.Items {
	// 	if l.Spec.Address == ip.String() {
	// 		r.l.V(1).Info("discarding because existing lease")
	// 		free = false
	// 		break
	// 	}
	// }
	// free = true
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
