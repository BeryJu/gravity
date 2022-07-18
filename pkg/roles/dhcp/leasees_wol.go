package dhcp

import (
	"fmt"
	"net"

	"beryju.io/gravity/pkg/extconfig"
	"github.com/sabhiram/go-wol/wol"
	"github.com/sirupsen/logrus"
)

func (l *Lease) sendWOL() error {
	bcast, err := net.ResolveUDPAddr("udp", "255.255.255.255:9")
	if err != nil {
		return err
	}
	laddr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:9", extconfig.Get().Instance.IP))
	if err != nil {
		return err
	}
	mp, err := wol.New(l.Identifier)
	if err != nil {
		return err
	}
	bs, err := mp.Marshal()
	if err != nil {
		return err
	}
	conn, err := net.DialUDP("udp", laddr, bcast)
	if err != nil {
		return err
	}
	defer conn.Close()

	l.log.WithFields(logrus.Fields{
		"mac":   l.Identifier,
		"raddr": bcast.String(),
		"laddr": laddr.String(),
	}).Info("Attempting WOL")
	n, err := conn.Write(bs)
	if err == nil && n != 102 {
		err = fmt.Errorf("magic packet sent was %d bytes (expected 102 bytes sent)", n)
	}
	if err != nil {
		return err
	}
	return nil
}
