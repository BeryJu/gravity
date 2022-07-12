package dhcp

import (
	"encoding/hex"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
	log "github.com/sirupsen/logrus"
)

func (ro *DHCPRole) loggingHandler4(inner server4.Handler) server4.Handler {
	return func(conn net.PacketConn, peer net.Addr, m *dhcpv4.DHCPv4) {
		f := log.Fields{
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
			"clientHWAddr":     m.ClientHWAddr.String(),
			"serverHostName":   m.ServerHostName,
			"clientIdentifier": hex.EncodeToString(m.Options.Get(dhcpv4.OptionClientIdentifier)),
		}
		start := time.Now()
		inner(conn, peer, m)
		duration := float64(time.Since(start)) / float64(time.Millisecond)
		f["runtimeMS"] = fmt.Sprintf("%0.3f", duration)
		ro.log.WithFields(f).Infof("DHCPv4 %s", strings.ToTitle(m.MessageType().String()))
	}
}
