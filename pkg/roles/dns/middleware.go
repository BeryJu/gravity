package dns

import (
	"fmt"
	"net"
	"time"

	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

func (ro *DNSRole) loggingHandler(inner dns.HandlerFunc) dns.HandlerFunc {
	return func(w dns.ResponseWriter, m *dns.Msg) {
		start := time.Now()
		fw := NewFakeDNSWriter(w)
		inner(fw, m)
		w.WriteMsg(fw.msg)
		duration := float64(time.Since(start)) / float64(time.Millisecond)
		var clientIP = ""
		switch addr := w.RemoteAddr().(type) {
		case *net.UDPAddr:
			clientIP = addr.IP.String()
		case *net.TCPAddr:
			clientIP = addr.IP.String()
		}
		f := log.Fields{
			"runtimeMS": fmt.Sprintf("%0.3f", duration),
			"client":    clientIP,
			"response":  dns.RcodeToString[fw.msg.Rcode],
		}
		msg := "DNS Request"
		if len(m.Question) > 0 {
			msg = m.Question[0].Name
		}
		for idx, a := range fw.msg.Answer {
			f[fmt.Sprintf("answer[%d]", idx)] = dns.TypeToString[a.Header().Rrtype]
		}
		ro.log.WithFields(f).Info(msg)
	}
}
