package dns

import (
	"fmt"
	"net"
	"time"

	"beryju.io/gravity/pkg/roles/dns/utils"
	"github.com/getsentry/sentry-go"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

func (r *DNSRole) recoverMiddleware(inner dns.HandlerFunc) dns.HandlerFunc {
	return func(w dns.ResponseWriter, m *dns.Msg) {
		defer func() {
			err := recover()
			if err == nil {
				return
			}
			if e, ok := err.(error); ok {
				r.log.WithError(e).Warning("recover in dns handler")
				sentry.CaptureException(e)
			} else {
				r.log.WithField("panic", err).Warning("recover in dns handler")
			}
			// ensure DNS query gets some sort of response to prevent
			// clients hanging
			fallback := new(dns.Msg)
			fallback.SetReply(m)
			fallback.SetRcode(m, dns.RcodeServerFailure)
			w.WriteMsg(fallback)
		}()
		inner(w, m)
	}
}

func (r *DNSRole) loggingMiddleware(inner dns.HandlerFunc) dns.HandlerFunc {
	return func(w dns.ResponseWriter, m *dns.Msg) {
		start := time.Now()
		fw := utils.NewFakeDNSWriter(w)
		inner(fw, m)
		w.WriteMsg(fw.Msg())
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
			"response":  dns.RcodeToString[fw.Msg().Rcode],
		}
		msg := "DNS Request"
		if len(m.Question) > 0 {
			msg = m.Question[0].Name
		}
		for idx, a := range fw.Msg().Answer {
			f[fmt.Sprintf("answer[%d]", idx)] = dns.TypeToString[a.Header().Rrtype]
		}
		r.log.WithFields(f).Info(msg)
	}
}
