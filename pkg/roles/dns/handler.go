package dns

import (
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

type Handler interface {
	Handle(w *fakeDNSWriter, r *dns.Msg) *dns.Msg
	Log() *log.Entry
}
