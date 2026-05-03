package dhcp

import (
	"encoding/hex"

	"github.com/getsentry/sentry-go"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"go.uber.org/zap"
)

func (r *Role) recoverMiddleware4(inner Handler4) Handler4 {
	return func(req *Request4) *dhcpv4.DHCPv4 {
		defer func() {
			err := recover()
			if err == nil {
				return
			}
			if e, ok := err.(error); ok {
				r.log.Error("recover in dhcp handler", zap.Error(e))
				sentry.CaptureException(e)
			} else {
				r.log.Error("recover in dhcp handler", zap.Any("panic", err))
			}
		}()
		return inner(req)
	}
}

func (r *Role) logDHCPMessage(req *Request4, m *dhcpv4.DHCPv4, fields []zap.Field) {
	f := []zap.Field{
		zap.String("deviceIdentifier", r.DeviceIdentifier(m)),
		zap.String("opCode", m.OpCode.String()),
		zap.Uint8("hopCount", m.HopCount),
		zap.String("transactionID", m.TransactionID.String()),
		zap.String("flagsToString", m.FlagsToString()),
		zap.String("clientIPAddr", m.ClientIPAddr.String()),
		zap.String("yourIPAddr", m.YourIPAddr.String()),
		zap.String("serverIPAddr", m.ServerIPAddr.String()),
		zap.String("gatewayIPAddr", m.GatewayIPAddr.String()),
		zap.String("hostname", m.HostName()),
		zap.String("clientIdentifier", hex.EncodeToString(m.Options.Get(dhcpv4.OptionClientIdentifier))),
		zap.String("messageType", m.MessageType().String()),
	}
	req.log.With(f...).With(fields...).Info("DHCP packet")
}

func (r *Role) loggingMiddleware4(inner Handler4) Handler4 {
	return func(req *Request4) *dhcpv4.DHCPv4 {
		f := []zap.Field{
			zap.String("client", req.peer.String()),
			zap.String("type", "request"),
		}
		r.logDHCPMessage(req, req.DHCPv4, f)
		return inner(req)
	}
}
