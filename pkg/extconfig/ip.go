package extconfig

import (
	"errors"
	"net"

	"go.uber.org/zap"
)

func mustParseCIDR(cidr string) *net.IPNet {
	_, net, err := net.ParseCIDR(cidr)
	if err != nil {
		panic(err)
	}
	return net
}

var ulaPrefix = mustParseCIDR("fc00::/7")

func (e *ExtConfig) GetIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range ifaces {
		if i.Flags&net.FlagLoopback != 0 {
			continue
		}
		if i.Flags&net.FlagUp == 0 {
			continue
		}
		addrs, err := i.Addrs()
		if err != nil {
			e.Logger().Debug("failed to get IPs from interface", zap.Error(err), zap.String("if", i.Name))
			continue
		}
		for _, addr := range addrs {
			ip, _, err := net.ParseCIDR(addr.String())
			if err != nil {
				continue
			}
			if ip.IsLinkLocalUnicast() || ip.IsLinkLocalMulticast() {
				continue
			}
			if ulaPrefix.Contains(ip) {
				continue
			}
			e.Logger().Debug("Detected IP of instance", zap.String("ip", ip.String()))
			return ip, nil
		}
	}
	return nil, errors.New("failed to find IP, set `INSTANCE_IP`")
}
