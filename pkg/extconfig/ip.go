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

func (e *ExtConfig) checkInterface(i net.Interface) ([]net.IP, error) {
	ips := []net.IP{}
	if i.Flags&net.FlagLoopback != 0 {
		return ips, errors.New("interface is loopback")
	}
	if i.Flags&net.FlagUp == 0 {
		return ips, errors.New("interface is down")
	}
	addrs, err := i.Addrs()
	if err != nil {
		return ips, err
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
		ips = append(ips, ip)
	}
	return ips, nil
}

func (e *ExtConfig) GetInterfaceForIP(forIp net.IP) (*net.Interface, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range ifaces {
		addrs, err := e.checkInterface(i)
		if err != nil {
			e.Logger().Debug("failed to get IPs from interface", zap.Error(err), zap.String("if", i.Name))
			continue
		}
		for _, addr := range addrs {
			if addr.Equal(forIp) {
				return &i, nil
			}
		}
	}
	return nil, nil
}

func (e *ExtConfig) GetIP() (net.IP, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return nil, err
	}
	for _, i := range ifaces {
		addrs, err := e.checkInterface(i)
		if err != nil {
			e.Logger().Debug("failed to get IPs from interface", zap.Error(err), zap.String("if", i.Name))
			continue
		}
		if len(addrs) < 1 {
			continue
		}
		e.Logger().Debug("Detected IP of instance", zap.String("ip", addrs[0].String()))
		return addrs[0], nil
	}
	return nil, errors.New("failed to find IP, set `INSTANCE_IP`")
}
