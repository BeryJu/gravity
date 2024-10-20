package dns_test

import (
	"errors"
	"net"

	"github.com/miekg/dns"
)

type NullDNSWriter struct {
	closed bool
	msg    *dns.Msg
}

func (nw *NullDNSWriter) LocalAddr() net.Addr {
	return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 53}
}

func (nw *NullDNSWriter) RemoteAddr() net.Addr {
	return &net.TCPAddr{IP: net.ParseIP("127.0.0.1"), Port: 1053}
}
func (nw *NullDNSWriter) WriteMsg(msg *dns.Msg) error { nw.msg = msg; return nil }
func (nw *NullDNSWriter) Write([]byte) (int, error)   { return 0, errors.New("Write not supported") }
func (nw *NullDNSWriter) Close() error                { nw.closed = true; return nil }
func (nw *NullDNSWriter) TsigStatus() error           { return nil }
func (nw *NullDNSWriter) TsigTimersOnly(v bool)       {}
func (nw *NullDNSWriter) Hijack()                     {}
func (nw *NullDNSWriter) Msg() *dns.Msg               { return nw.msg }

func NewNullDNSWriter() *NullDNSWriter {
	return &NullDNSWriter{}
}
