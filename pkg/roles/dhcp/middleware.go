package dhcp

import (
	"encoding/hex"
	"fmt"
	"net"
	"time"

	"github.com/getsentry/sentry-go"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
	log "github.com/sirupsen/logrus"
)

func (r *Role) recoverMiddleware4(inner server4.Handler) server4.Handler {
	return func(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
		defer func() {
			err := recover()
			if err == nil {
				return
			}
			if e, ok := err.(error); ok {
				r.log.WithError(e).Warning("recover in dhcp handler")
				sentry.CaptureException(e)
			} else {
				r.log.WithField("panic", err).Warning("recover in dhcp handler")
			}
		}()
		inner(conn, peer, m)
	}
}

func (r *Role) loggingMiddleware4(inner server4.Handler) server4.Handler {
	return func(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
		f := log.Fields{
			"msgType":          m.MessageType(),
			"client":           peer.String(),
			"localAddr":        conn.LocalAddr().String(),
			"opCode":           m.OpCode.String(),
			"hwType":           m.HWType.String(),
			"hopCount":         m.HopCount,
			"transactionID":    m.TransactionID.String(),
			"flagsToString":    m.FlagsToString(),
			"clientIPAddr":     m.ClientIPAddr.String(),
			"yourIPAddr":       m.YourIPAddr.String(),
			"serverIPAddr":     m.ServerIPAddr.String(),
			"gatewayIPAddr":    m.GatewayIPAddr.String(),
			"serverHostName":   m.ServerHostName,
			"clientIdentifier": hex.EncodeToString(m.Options.Get(dhcpv4.OptionClientIdentifier)),
		}
		start := time.Now()
		inner(conn, peer, m)
		duration := float64(time.Since(start)) / float64(time.Millisecond)
		f["runtimeMS"] = fmt.Sprintf("%0.3f", duration)
		r.log.WithFields(f).Info(m.ClientHWAddr.String())
	}
}
