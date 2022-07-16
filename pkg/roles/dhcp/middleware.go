package dhcp

import (
	"encoding/hex"
	"net"

	"github.com/getsentry/sentry-go"
	"github.com/insomniacslk/dhcp/dhcpv4"
	log "github.com/sirupsen/logrus"
)

func (r *Role) recoverMiddleware4(inner Handler4) Handler4 {
	return func(peer net.Addr, m *dhcpv4.DHCPv4) *dhcpv4.DHCPv4 {
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
		return inner(peer, m)
	}
}

func (r *Role) logDHCPMessage(m *dhcpv4.DHCPv4, fields log.Fields) {
	f := log.Fields{
		"deviceIdentifier": r.DeviceIdentifier(m),
		"opCode":           m.OpCode.String(),
		"hopCount":         m.HopCount,
		"transactionID":    m.TransactionID.String(),
		"flagsToString":    m.FlagsToString(),
		"clientIPAddr":     m.ClientIPAddr.String(),
		"yourIPAddr":       m.YourIPAddr.String(),
		"serverIPAddr":     m.ServerIPAddr.String(),
		"gatewayIPAddr":    m.GatewayIPAddr.String(),
		"hostname":         m.HostName(),
		"clientIdentifier": hex.EncodeToString(m.Options.Get(dhcpv4.OptionClientIdentifier)),
	}
	r.log.WithFields(f).WithFields(fields).Info(m.MessageType().String())
}

func (r *Role) loggingMiddleware4(inner Handler4) Handler4 {
	return func(peer net.Addr, m *dhcpv4.DHCPv4) *dhcpv4.DHCPv4 {
		f := log.Fields{
			"client": peer.String(),
		}
		r.logDHCPMessage(m, f)
		return inner(peer, m)
	}
}
