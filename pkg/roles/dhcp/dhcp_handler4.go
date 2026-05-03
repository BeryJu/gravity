package dhcp

import (
	"context"
	"errors"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"

	"beryju.io/gravity/pkg/extconfig"
	"github.com/getsentry/sentry-go"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/net/ipv4"
)

func getIP(addr net.Addr) net.IP {
	clientIP := ""
	switch addr := addr.(type) {
	case *net.UDPAddr:
		clientIP = addr.IP.String()
	}
	return net.ParseIP(clientIP)
}

var ErrNilResponse = errors.New("no DHCP response")

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

func (h *handler4) Serve() error {
	for {
		b := *bufpool.Get().(*[]byte)
		b = b[:MaxDatagram] // Reslice to max capacity in case the buffer in pool was resliced smaller

		n, oob, peer, err := h.pc.ReadFrom(b)
		if err != nil {
			return err
		}
		go func(buf []byte, oob *ipv4.ControlMessage, peer net.Addr) {
			_ = h.Handle(buf, oob, peer)
		}(b[:n], oob, peer.(*net.UDPAddr))
	}
}

var debugDHCPGatewayReplyPeer bool

func init() {
	debugDHCPGatewayReplyPeer = os.Getenv("GRAVITY_DEBUG_DHCP_GATEWAY_REPLY_CIADDR") != ""
}

func (h *handler4) Handle(buf []byte, oob *ipv4.ControlMessage, peer net.Addr) error {
	if extconfig.Get().ListenOnlyMode {
		return nil
	}
	context, canc := context.WithCancel(h.role.ctx)
	defer canc()
	m, err := dhcpv4.FromBytes(buf)
	bufpool.Put(&buf)
	if err != nil {
		return fmt.Errorf("error parsing dhcpv4 request: %w", err)
	}

	r := h.role.NewRequest4(m)
	r.peer = peer
	r.Context = context
	r.oob = oob

	span := sentry.StartTransaction(r.Context, h.role.DeviceIdentifier(r.DHCPv4))
	span.Op = "gravity.dhcp.request"
	hub := sentry.GetHubFromContext(span.Context())
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	hub.Scope().SetUser(sentry.User{
		Username:  m.HostName(),
		IPAddress: strings.Split(peer.String(), ":")[0],
	})

	span.SetData("http.request.method", m.MessageType().String())
	defer span.Finish()
	resp := h.HandleRequest(r)

	if resp == nil {
		r.log.Debug("handler4: dropping request because response is nil")
		return ErrNilResponse
	}

	h.role.logDHCPMessage(r, resp, []zapcore.Field{
		zap.String("type", "response"),
	})
	useEthernet := false
	var p *net.UDPAddr
	if !r.GatewayIPAddr.IsUnspecified() {
		r.log.Debug("sending response to gateway")
		// giaddr should be set to the Relay's IP Address, however it is the IP of the subnet
		// the client should get an IP for. We might not always be able to directly reply to that IP
		// especially in environments where we can't adjust firewall/routing rules. (like e2e tests)
		// when this environment variable is set, reply directly to the IP we got the UDP request from,
		// which is not the RFC defined behaviour
		p = &net.UDPAddr{IP: r.GatewayIPAddr, Port: dhcpv4.ServerPort}
		if debugDHCPGatewayReplyPeer {
			p.IP = getIP(r.peer)
		}
	} else if resp.MessageType() == dhcpv4.MessageTypeNak {
		r.log.Debug("sending response to bcast (NAK)")
		p = &net.UDPAddr{IP: net.IPv4bcast, Port: dhcpv4.ClientPort}
	} else if !r.ClientIPAddr.IsUnspecified() {
		r.log.Debug("sending response to client")
		p = &net.UDPAddr{IP: r.ClientIPAddr, Port: dhcpv4.ClientPort}
	} else if r.IsBroadcast() {
		r.log.Debug("sending response to bcast")
		p = &net.UDPAddr{IP: net.IPv4bcast, Port: dhcpv4.ClientPort}
	} else {
		// sends a layer2 frame so that we can define the destination MAC address
		p = &net.UDPAddr{IP: resp.YourIPAddr, Port: dhcpv4.ClientPort}
		useEthernet = true
	}

	var woob *ipv4.ControlMessage
	if p.IP.Equal(net.IPv4bcast) || p.IP.IsLinkLocalUnicast() || useEthernet {
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
		r.log.Debug("sending via ethernet")
		intf, err := net.InterfaceByIndex(woob.IfIndex)
		if err != nil {
			r.log.Error("handler4: Can not get Interface for index", zap.Error(err), zap.Int("index", woob.IfIndex))
			return fmt.Errorf("handler4: Can not get Interface for index: %w", err)
		}
		err = h.sendEthernet(*intf, resp)
		if err != nil {
			r.log.Error("handler4: Cannot send Ethernet packet", zap.Error(err))
			return fmt.Errorf("handler4: Cannot send Ethernet packet: %w", err)
		}
	} else {
		b := resp.ToBytes()
		n, err := h.pc.WriteTo(b, woob, p)
		if err != nil {
			r.log.Error("handler4: conn.Write failed", zap.Error(err), zap.String("peer", p.String()))
			return fmt.Errorf("handler4: failed to write response: %w", err)
		}
		if len(b) != n {
			r.log.Warn("handler4: did not send all bytes", zap.Int("length", len(b)), zap.Int("sent", n))
		}
	}
	return nil
}

func (h *handler4) HandleRequest(r *Request4) *dhcpv4.DHCPv4 {
	if r.OpCode != dhcpv4.OpcodeBootRequest {
		h.role.log.Info("handler4: unsupported opcode", zap.String("opcode", r.OpCode.String()))
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
	case dhcpv4.MessageTypeRelease:
		handler = h.role.HandleDHCPRelease4
	default:
		r.log.Info("Unsupported message type", zap.String("dhcpMsg", mt.String()))
		return nil
	}

	return h.role.recoverMiddleware4(
		h.role.loggingMiddleware4(
			handler,
		),
	)(r)
}
