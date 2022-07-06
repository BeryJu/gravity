package dhcp

import (
	"encoding/base64"
	"net"
	"time"

	"beryju.io/ddet/pkg/extconfig"
	"github.com/insomniacslk/dhcp/dhcpv4"
)

func (r *DHCPRole) replyWithLease(
	lease *Lease,
	conn net.PacketConn,
	peer net.Addr,
	m *dhcpv4.DHCPv4,
	modifyResponse func(*dhcpv4.DHCPv4) *dhcpv4.DHCPv4,
) {
	rep, err := dhcpv4.NewReplyFromRequest(m)
	if err != nil {
		r.log.WithError(err).Warning("failed to create reply")
		return
	}
	rep = modifyResponse(rep)

	ipLeaseDuration, err := time.ParseDuration(lease.AddressLeaseTime)
	if err != nil {
		r.log.WithField("default", "24h").WithError(err).Warning("failed to parse address lease duration, defaulting")
		ipLeaseDuration = time.Hour * 24
	}
	rep.UpdateOption(dhcpv4.OptDNS(net.IP(extconfig.Get().Instance.IP)))
	rep.UpdateOption(dhcpv4.OptIPAddressLeaseTime(ipLeaseDuration))
	rep.UpdateOption(dhcpv4.OptSubnetMask(lease.scope.ipam.GetSubnetMask()))

	rep.YourIPAddr = net.ParseIP(lease.Address)
	rep.UpdateOption(dhcpv4.OptHostName(lease.Hostname))

	for _, opt := range lease.scope.Options {
		finalVal := make([]byte, 0)
		r.log.Debug("applying options from optionset", "option", opt.Tag)
		if opt.Tag == nil {
			continue
		}

		// Values which are directly converted from string to byte
		if opt.Value != nil {
			i := net.ParseIP(*opt.Value)
			if i == nil {
				finalVal = []byte(*opt.Value)
			} else {
				finalVal = dhcpv4.IPs([]net.IP{i}).ToBytes()
			}
		}

		// For non-stringable values, get b64 decoded values
		if len(opt.Value64) > 0 {
			values64 := make([]byte, 0)
			for _, v := range opt.Value64 {
				va, err := base64.StdEncoding.DecodeString(v)
				if err != nil {
					r.log.WithError(err).Warning("failed to convert base64 value to byte")
				} else {
					values64 = append(values64, va...)
				}
			}
			finalVal = values64
		}
		dopt := dhcpv4.OptGeneric(dhcpv4.GenericOptionCode(*opt.Tag), finalVal)
		rep.UpdateOption(dopt)
	}

	r.log.Trace(rep.Summary(), "peer", peer.String())
	if _, err := conn.WriteTo(rep.ToBytes(), peer); err != nil {
		r.log.WithError(err).Warning("failed to write reply")
	}
}

func (r *DHCPRole) findScopeForRequest(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) *Scope {
	var match *Scope
	for _, scope := range r.scopes {
		// TODO: priority and order
		if scope.match(conn, peer, m) {
			r.log.WithField("name", scope.Name).Debug("selected scope based on match")
			match = scope
		}
		if match == nil && scope.Default {
			r.log.WithField("name", scope.Name).Debug("selected scope based on default state")
			match = scope
		}
	}
	if match != nil {
		r.log.Trace("found scope for request")
	}
	return match
}
