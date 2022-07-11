package dhcp

import (
	"fmt"
	"net"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
	log "github.com/sirupsen/logrus"
)

func (ro *DHCPRole) loggingHandler4(inner server4.Handler) server4.Handler {
	return func(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
		f := log.Fields{
			"client":         peer.String(),
			"opCode":         m.OpCode,
			"hwType":         m.HWType,
			"hopCount":       m.HopCount,
			"transactionID":  m.TransactionID,
			"numSeconds":     m.NumSeconds,
			"flagsToString":  m.FlagsToString(),
			"flags":          m.Flags,
			"clientIPAddr":   m.ClientIPAddr,
			"yourIPAddr":     m.YourIPAddr,
			"serverIPAddr":   m.ServerIPAddr,
			"gatewayIPAddr":  m.GatewayIPAddr,
			"clientHWAddr":   m.ClientHWAddr,
			"serverHostName": m.ServerHostName,
			"bootFileName":   m.BootFileName,
		}
		start := time.Now()
		inner(conn, peer, m)
		duration := float64(time.Since(start)) / float64(time.Millisecond)
		f["runtimeMS"] = fmt.Sprintf("%0.3f", duration)
		ro.log.WithFields(f).Info("DHCPv4 request")
	}
}
