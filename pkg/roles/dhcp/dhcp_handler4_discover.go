package dhcp

import (
	"context"
	"net"

	"beryju.io/ddet/pkg/roles/dhcp/types"
	"github.com/insomniacslk/dhcp/dhcpv4"
)

func (r *DHCPRole) handleDHCPDiscover4(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
	match, err := r.i.KV().KV.Get(context.TODO(), r.i.KV().Key(types.KeyRole, types.KeyLeases, m.ClientHWAddr.String()))
	var lease *Lease
	if err == nil && len(match.Kvs) > 0 {
		// TODO: Update lease of lease
		lease, err = r.leaseFromKV(match.Kvs[0])
		if err != nil {
			// TODO: remove invalid lease?
			r.log.WithError(err).Warning("failed to parse lease")
			return
		}
	} else {
		r.log.Debug("no lease found, creating new")
		scope := r.findScopeForRequest(conn, peer, m)
		if scope == nil {
			r.log.Warning("no scope found")
			return
		}
		r.log.Debug("found scope for new lease")
		lease = scope.createLeaseFor(conn, peer, m)
		lease.put(int64(r.cfg.LeaseNegotiateTimeout))
	}

	lease.reply(conn, peer, m, func(d *dhcpv4.DHCPv4) *dhcpv4.DHCPv4 {
		d.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeOffer))
		return d
	})
}
