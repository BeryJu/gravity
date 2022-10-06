package dhcp

import (
	"context"
	"fmt"
	"net"
	"sync"

	"beryju.io/gravity/pkg/extconfig"
	"github.com/getsentry/sentry-go"
	"github.com/google/uuid"
	"github.com/insomniacslk/dhcp/dhcpv4"
	log "github.com/sirupsen/logrus"
	"golang.org/x/net/ipv4"
)

// Credit to CoreDHCP
// https://github.com/coredhcp/coredhcp/blob/master/server/handle.go

type handler4 struct {
	role  *Role
	pc    *ipv4.PacketConn
	iface net.Interface
}

// XXX: performance-wise, Pool may or may not be good (see https://github.com/golang/go/issues/23199)
// Interface is good for what we want. Maybe "just" trust the GC and we'll be fine ?
var bufpool = sync.Pool{New: func() interface{} { r := make([]byte, MaxDatagram); return &r }}

// MaxDatagram is the maximum length of message that can be received.
const MaxDatagram = 1 << 16

type Handler4 func(req *Request4) *dhcpv4.DHCPv4

type Request4 struct {
	*dhcpv4.DHCPv4
	peer    net.Addr
	log     *log.Entry
	Context context.Context
	oob     *ipv4.ControlMessage
}

func (r *Role) NewRequest4(dhcp *dhcpv4.DHCPv4) *Request4 {
	return &Request4{
		DHCPv4:  dhcp,
		Context: r.ctx,
		peer:    &net.UDPAddr{},
		log:     r.log.WithField("request", fmt.Sprintf("%s-%s", uuid.New().String(), dhcp.TransactionID.String())),
	}
}

func (h *handler4) Serve() error {
	for {
		b := *bufpool.Get().(*[]byte)
		b = b[:MaxDatagram] //Reslice to max capacity in case the buffer in pool was resliced smaller

		n, oob, peer, err := h.pc.ReadFrom(b)
		if err != nil {
			h.role.log.WithError(err).Debug("Error reading from connection")
			return err
		}
		go h.handle(b[:n], oob, peer.(*net.UDPAddr))
	}
}

func (h *handler4) handle(buf []byte, oob *ipv4.ControlMessage, _peer net.Addr) {
	if extconfig.Get().ListenOnlyMode {
		return
	}
	context, canc := context.WithCancel(context.Background())
	defer canc()
	m, err := dhcpv4.FromBytes(buf)
	bufpool.Put(&buf)
	if err != nil {
		h.role.log.WithError(err).Info("Error parsing DHCPv4 request")
		return
	}

	r := h.role.NewRequest4(m)
	r.peer = _peer
	r.Context = context
	r.oob = oob

	span := sentry.StartSpan(
		r.Context,
		"gravity.roles.dhcp.request",
		sentry.TransactionName("gravity.roles.dhcp"),
	)
	span.Description = m.MessageType().String()
	defer span.Finish()
	resp := h.HandleRequest(r)

	if resp != nil {
		h.role.logDHCPMessage(r, resp, log.Fields{})
		useEthernet := false
		var peer *net.UDPAddr
		if !r.GatewayIPAddr.IsUnspecified() {
			// TODO: make RFC8357 compliant
			peer = &net.UDPAddr{IP: r.GatewayIPAddr, Port: dhcpv4.ServerPort}
		} else if resp.MessageType() == dhcpv4.MessageTypeNak {
			peer = &net.UDPAddr{IP: net.IPv4bcast, Port: dhcpv4.ClientPort}
		} else if !r.ClientIPAddr.IsUnspecified() {
			peer = &net.UDPAddr{IP: r.ClientIPAddr, Port: dhcpv4.ClientPort}
		} else if r.IsBroadcast() {
			peer = &net.UDPAddr{IP: net.IPv4bcast, Port: dhcpv4.ClientPort}
		} else {
			//sends a layer2 frame so that we can define the destination MAC address
			peer = &net.UDPAddr{IP: resp.YourIPAddr, Port: dhcpv4.ClientPort}
			useEthernet = true
		}

		var woob *ipv4.ControlMessage
		if peer.IP.Equal(net.IPv4bcast) || peer.IP.IsLinkLocalUnicast() || useEthernet {
			// Direct broadcasts, link-local and layer2 unicasts to the interface the ruest was
			// received on. Other packets should use the normal routing table in
			// case of asymetric routing
			switch {
			case h.iface.Index != 0:
				woob = &ipv4.ControlMessage{IfIndex: h.iface.Index}
			case r.oob != nil && r.oob.IfIndex != 0:
				woob = &ipv4.ControlMessage{IfIndex: r.oob.IfIndex}
			default:
				r.log.Error("HandleMsg4: Did not receive interface information")
			}
		}

		if useEthernet {
			r.log.Trace("sending via ethernet")
			intf, err := net.InterfaceByIndex(woob.IfIndex)
			if err != nil {
				r.log.WithError(err).WithField("index", woob.IfIndex).Error("handler4: Can not get Interface for index")
				return
			}
			err = sendEthernet(*intf, resp)
			if err != nil {
				r.log.WithError(err).Error("handler4: Cannot send Ethernet packet")
			}
		} else {
			if _, err := h.pc.WriteTo(resp.ToBytes(), woob, peer); err != nil {
				r.log.WithField("peer", peer).WithError(err).Error("handler4: conn.Write failed")
			}
		}
	} else {
		r.log.Debug("handler4: dropping request because response is nil")
	}
}

func (h *handler4) HandleRequest(r *Request4) *dhcpv4.DHCPv4 {
	if r.OpCode != dhcpv4.OpcodeBootRequest {
		h.role.log.WithField("opcode", r.OpCode.String()).Info("handler4: unsupported opcode")
		return nil
	}
	var handler Handler4
	switch mt := r.MessageType(); mt {
	case dhcpv4.MessageTypeDiscover:
		handler = h.role.HandleDHCPDiscover4
	case dhcpv4.MessageTypeRequest:
		handler = h.role.HandleDHCPRequest4
	case dhcpv4.MessageTypeDecline:
		handler = h.role.HandleDHCPDecline4
	default:
		r.log.WithField("msg", mt.String()).Info("Unsupported message type")
		return nil
	}

	return h.role.recoverMiddleware4(
		h.role.loggingMiddleware4(
			handler,
		),
	)(r)
}
