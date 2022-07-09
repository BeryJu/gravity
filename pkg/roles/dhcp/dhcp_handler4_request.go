package dhcp

import (
	"net"

	"github.com/insomniacslk/dhcp/dhcpv4"
)

func (r *DHCPRole) handleDHCPRequest4(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	match := r.findLease(m)

	if match == nil {
		r.log.Debug("no lease found, creating new")
		scope := r.findScopeForRequest(conn, peer, m)
		if scope == nil {
			return
		}
		r.log.Debug("found scope for new lease")
		match = scope.createLeaseFor(conn, peer, m)
		r.log.Debug("creating new lease")
		match.put(scope.TTL)
	}

	// TODO: Update
	match.reply(conn, peer, m, func(d *dhcpv4.DHCPv4) *dhcpv4.DHCPv4 {
		d.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeAck))
		return d
	})
}

func (r *DHCPRole) findLease(m *dhcpv4.DHCPv4) *Lease {
	return nil
}

func (r *DHCPRole) createLease(lease *Lease, scope *Scope) {
	// err := r.Create(context.Background(), lease)
	// if err != nil {
	// 	r.l.Error(err, "failed to create lease")
	// }
	// r.i.DispatchEvent(roles.NewEvent[any](
	// 	[]string{KeyRole, KeyLeases, "create"},
	// 	// Leas
	// ))

	// TODO: Send DNS create event
	// dns, err := dns.GetDNSProviderForScope(*scope)
	// if err != nil {
	// 	r.l.Error(err, "failed to get DNS provider")
	// 	return
	// }
	// err = dns.CreateRecord(lease)
	// if err != nil {
	// 	r.l.Error(err, "failed to delete DNS record")
	// }
}
