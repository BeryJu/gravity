package dns

import (
	"net"
	"time"

	"beryju.io/gravity/pkg/roles/dns/utils"
	"github.com/getsentry/sentry-go"
	"github.com/miekg/dns"
	"go.uber.org/zap"
)

func (r *Role) recoverMiddleware(inner dns.HandlerFunc) dns.HandlerFunc {
	return func(w dns.ResponseWriter, m *dns.Msg) {
		defer func() {
			err := recover()
			if err == nil {
				return
			}
			if e, ok := err.(error); ok {
				r.log.Error("recover in dns handler", zap.Error(e))
				sentry.CaptureException(e)
			} else {
				r.log.Error("recover in dns handler", zap.Any("panic", err))
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

func (r *Role) dnsRRToValue(ro dns.RR) string {
	switch v := ro.(type) {
	case *dns.A:
		return v.A.String()
	case *dns.AAAA:
		return v.AAAA.String()
	case *dns.PTR:
		return v.Ptr
	case *dns.MX:
		return v.Mx
	case *dns.CNAME:
		return v.Target
	default:
		return ro.String()
	}
}

func (r *Role) loggingMiddleware(inner dns.HandlerFunc) dns.HandlerFunc {
	getIP := func(addr net.Addr) string {
		clientIP := ""
		switch addr := addr.(type) {
		case *net.UDPAddr:
			clientIP = addr.IP.String()
		case *net.TCPAddr:
			clientIP = addr.IP.String()
		}
		return clientIP
	}

	return func(w dns.ResponseWriter, m *dns.Msg) {
		fw := utils.NewFakeDNSWriter(w)
		start := time.Now()
		inner(fw, m)
		finish := time.Since(start)
		w.WriteMsg(fw.Msg())

		queryNames := make([]string, len(m.Question))
		queryTypes := make([]string, len(m.Question))
		answerRecords := make([]string, len(fw.Msg().Answer))
		answerTypes := make([]string, len(fw.Msg().Answer))
		for idx, q := range m.Question {
			queryNames[idx] = q.Name
			queryTypes[idx] = dns.TypeToString[q.Qtype]
		}
		for idx, a := range fw.Msg().Answer {
			answerRecords[idx] = r.dnsRRToValue(a)
			answerTypes[idx] = dns.TypeToString[a.Header().Rrtype]
		}
		f := []zap.Field{
			zap.Duration("runtime", finish),
			zap.String("client", getIP(w.RemoteAddr())),
			zap.String("response", dns.RcodeToString[fw.Msg().Rcode]),
			zap.Strings("queryNames", queryNames),
			zap.Strings("queryTypes", queryTypes),
			zap.Strings("answerRecords", answerRecords),
			zap.Strings("answerTypes", answerTypes),
		}
		r.log.With(f...).Info("DNS Query")
	}
}
