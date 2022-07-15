package extconfig

import (
	"errors"
	"net"

	log "github.com/sirupsen/logrus"
)

func GetIP() (net.IP, error) {
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
			log.WithError(err).WithField("if", i).Trace("failed to get IPs from interface")
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
			log.Debugf("Detected IP of instance as %s", ip.String())
			return ip, nil
		}
	}
	return nil, errors.New("failed to find IP, set `INSTANCE_IP`")
}
