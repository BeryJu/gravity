package utils

import (
	"errors"
	"net"

	"github.com/miekg/dns"
)

type FakeDNSWriter struct {
	inner dns.ResponseWriter
	msg   *dns.Msg
}

func (fd *FakeDNSWriter) LocalAddr() net.Addr         { return fd.inner.LocalAddr() }
func (fd *FakeDNSWriter) RemoteAddr() net.Addr        { return fd.inner.RemoteAddr() }
func (fd *FakeDNSWriter) WriteMsg(msg *dns.Msg) error { fd.msg = msg; return nil }
func (fd *FakeDNSWriter) Write([]byte) (int, error)   { return 0, errors.New("Write not supported") }
func (fd *FakeDNSWriter) Close() error                { return fd.inner.Close() }
func (fd *FakeDNSWriter) TsigStatus() error           { return fd.inner.TsigStatus() }
func (fd *FakeDNSWriter) TsigTimersOnly(v bool)       { fd.inner.TsigTimersOnly(v) }
func (fd *FakeDNSWriter) Hijack()                     { fd.inner.Hijack() }
func (fd *FakeDNSWriter) Msg() *dns.Msg               { return fd.msg }

func NewFakeDNSWriter(w dns.ResponseWriter) *FakeDNSWriter {
	return &FakeDNSWriter{
		inner: w,
	}
}
