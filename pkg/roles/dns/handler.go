package dns

import (
	"github.com/miekg/dns"
)

type Handler interface {
	Handle(w *fakeDNSWriter, r *dns.Msg) *dns.Msg
	Identifier() string
}
