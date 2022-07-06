package dns

import (
	"fmt"
	"time"

	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

func (ro *DNSRole) loggingHandler(inner dns.HandlerFunc) dns.HandlerFunc {
	return func(w dns.ResponseWriter, m *dns.Msg) {
		start := time.Now()
		inner(w, m)
		duration := float64(time.Since(start)) / float64(time.Millisecond)
		ro.log.WithFields(log.Fields{
			"runtime_ms": fmt.Sprintf("%0.3f", duration),
			"client":     w.RemoteAddr().String(),
		}).Info("DNS request")
	}
}
