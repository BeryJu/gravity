package dns

import (
	"fmt"
	"net"
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles/dns/utils"
	"github.com/getsentry/sentry-go"
	"github.com/miekg/dns"
	"go.uber.org/zap"
)

func (r *Role) recoverMiddleware(inner dns.HandlerFunc) dns.HandlerFunc {
	return func(w dns.ResponseWriter, m *dns.Msg) {
		defer func() {
			err := extconfig.RecoverWrapper(recover())
			if err == nil {
				return
			}
			if e, ok := err.(error); ok {
				r.log.Warn("recover in dns handler", zap.Error(e))
				sentry.CaptureException(e)
			} else {
				r.log.Warn("recover in dns handler", zap.Any("panic", err))
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

func (r *Role) loggingMiddleware(inner dns.HandlerFunc) dns.HandlerFunc {
	return func(w dns.ResponseWriter, m *dns.Msg) {
		start := time.Now()
		fw := utils.NewFakeDNSWriter(w)
		inner(fw, m)
		w.WriteMsg(fw.Msg())
		var clientIP = ""
		switch addr := w.RemoteAddr().(type) {
		case *net.UDPAddr:
			clientIP = addr.IP.String()
		case *net.TCPAddr:
			clientIP = addr.IP.String()
		}
		f := []zap.Field{
			zap.Duration("runtime", time.Since(start)),
			zap.String("client", clientIP),
			zap.String("response", dns.RcodeToString[fw.Msg().Rcode]),
		}
		msg := "DNS Request"
		if len(m.Question) > 0 {
			msg = m.Question[0].Name
		}
		for idx, a := range fw.Msg().Answer {
			f = append(f, zap.String(fmt.Sprintf("answer[%d]", idx), dns.TypeToString[a.Header().Rrtype]))
		}
		r.log.With(f...).Info(msg)
	}
}
