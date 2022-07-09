package dns

import (
	"errors"
	"net"

	"github.com/miekg/dns"
)

type fakeDNSWriter struct {
	inner dns.ResponseWriter
	msg   *dns.Msg
}

func (fd *fakeDNSWriter) LocalAddr() net.Addr         { return fd.inner.LocalAddr() }
func (fd *fakeDNSWriter) RemoteAddr() net.Addr        { return fd.inner.RemoteAddr() }
func (fd *fakeDNSWriter) WriteMsg(msg *dns.Msg) error { fd.msg = msg; return nil }
func (fd *fakeDNSWriter) Write([]byte) (int, error)   { return 0, errors.New("Write not supported") }
func (fd *fakeDNSWriter) Close() error                { return fd.inner.Close() }
func (fd *fakeDNSWriter) TsigStatus() error           { return fd.inner.TsigStatus() }
func (fd *fakeDNSWriter) TsigTimersOnly(v bool)       { fd.inner.TsigTimersOnly(v) }
func (fd *fakeDNSWriter) Hijack()                     { fd.inner.Hijack() }

func NewFakeDNSWriter(w dns.ResponseWriter) *fakeDNSWriter {
	return &fakeDNSWriter{
		inner: w,
	}
}
