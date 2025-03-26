package dhcp

import (
	"fmt"
	"net"

	"beryju.io/gravity/pkg/extconfig"
	"github.com/sabhiram/go-wol/wol"
	"go.uber.org/zap"
)

func (l *Lease) sendWOL() error {
	bcast, err := net.ResolveUDPAddr("udp", "255.255.255.255:9")
	if err != nil {
		return err
	}
	laddr, err := net.ResolveUDPAddr("udp", extconfig.Get().Listen(9))
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
	defer func() {
		err := conn.Close()
		if err != nil {
			l.log.Warn("failed to close WOL connection", zap.Error(err))
		}
	}()

	l.log.Info("Attempting WOL", zap.String("mac", l.Identifier), zap.String("raddr", bcast.String()), zap.String("laddr", laddr.String()))
	n, err := conn.Write(bs)
	if err == nil && n != 102 {
		err = fmt.Errorf("magic packet sent was %d bytes (expected 102 bytes sent)", n)
	}
	if err != nil {
		return err
	}
	return nil
}
