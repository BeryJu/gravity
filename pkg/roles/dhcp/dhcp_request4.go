package dhcp

import (
	"context"
	"fmt"
	"net"

	"beryju.io/gravity/pkg/extconfig"
	"github.com/google/uuid"
	"github.com/insomniacslk/dhcp/dhcpv4"
	"go.uber.org/zap"
	"golang.org/x/net/ipv4"
)

type Request4 struct {
	*dhcpv4.DHCPv4
	peer      net.Addr
	log       *zap.Logger
	Context   context.Context
	oob       *ipv4.ControlMessage
	requestId string
}

type contextRequestID struct{}

func (r *Role) NewRequest4(dhcp *dhcpv4.DHCPv4) *Request4 {
	requestId := fmt.Sprintf("%s-%s", uuid.New().String(), dhcp.TransactionID.String())
	return &Request4{
		DHCPv4:    dhcp,
		Context:   context.WithValue(r.ctx, contextRequestID{}, requestId),
		peer:      &net.UDPAddr{},
		log:       r.log.With(zap.String("request", requestId)),
		requestId: requestId,
	}
}

// Use the instance ip unless the the interface is not bound
func (req *Request4) LocalIP() string {
	ip := extconfig.Get().Instance.IP
	if req.oob != nil {
		ief, err := net.InterfaceByIndex(req.oob.IfIndex)
		if err != nil {
			return ip
		}
		addrs, err := ief.Addrs()
		if err != nil {
			return ip
		}
		for _, addr := range addrs {
			if ipv4Addr := addr.(*net.IPNet).IP.To4(); ipv4Addr != nil {
				ip = ipv4Addr.String()
				req.log.Debug("Unbound interface found", zap.String("ifname", ief.Name), zap.String("ip", ip))
				return ip
			}
		}
	}
	return ip
}
