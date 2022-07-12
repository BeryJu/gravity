package dhcp

import (
	"context"
	"net"

	"beryju.io/gravity/pkg/roles/dhcp/types"
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
	}

	// Run the update in a go-routine since etcd might not be reachable and
	// we don't want to timeout
	go match.put(match.scope.TTL)
	match.reply(conn, peer, m, func(d *dhcpv4.DHCPv4) *dhcpv4.DHCPv4 {
		d.UpdateOption(dhcpv4.OptMessageType(dhcpv4.MessageTypeAck))
		return d
	})
}

func (r *DHCPRole) findLease(m *dhcpv4.DHCPv4) *Lease {
	match, err := r.i.KV().KV.Get(context.TODO(), r.i.KV().Key(types.KeyRole, types.KeyLeases, m.ClientHWAddr.String()))
	var lease *Lease
	if err != nil {
		r.log.WithError(err).Trace("no lease")
		return nil
	}
	if len(match.Kvs) < 1 {
		return nil
	}
	lease, err = r.leaseFromKV(match.Kvs[0])
	if err != nil {
		r.log.WithError(err).Warning("failed to get lease fromKV")
		return nil
	}
	return lease
}
