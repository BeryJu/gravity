package ip

import (
	"context"
	"math/rand"
	"net"
	"strings"
	"time"

	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

type IPForwarder struct {
	r   *net.Resolver
	log *log.Entry
}

func New(config map[string]string) *IPForwarder {
	forwarders := strings.Split(config["to"], ";")
	r := &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * time.Duration(10000),
			}
			resolver := forwarders[rand.Intn(len(forwarders))]
			if !strings.Contains(resolver, ":") {
				resolver += ":53"
			}
			return d.DialContext(ctx, network, resolver)
		},
	}

	return &IPForwarder{
		r:   r,
		log: log.WithField("forwarder", "ip"),
	}
}

func (ipf *IPForwarder) Handle(w dns.ResponseWriter, r *dns.Msg) {
	if len(r.Question) < 1 {
		ipf.log.Error("No question")
		return
	}
	question := r.Question[0]
	ips, err := ipf.r.LookupHost(context.Background(), question.Name)
	m := new(dns.Msg)
	m.SetReply(r)

	if err != nil {
		ipf.log.WithError(err).Warning("failed to forward")
		m.SetRcode(r, dns.RcodeServerFailure)
		w.WriteMsg(m)
		return
	}
	m.Answer = make([]dns.RR, len(ips))
	for idx, rawIp := range ips {
		ip := net.ParseIP(rawIp)
		if len(ip) > 15 {
			m.Answer[idx] = &dns.AAAA{
				Hdr: dns.RR_Header{
					Name:   question.Name,
					Rrtype: dns.TypeAAAA,
					Class:  dns.ClassINET,
					Ttl:    0,
				},
				AAAA: ip,
			}
		} else {
			m.Answer[idx] = &dns.A{
				Hdr: dns.RR_Header{
					Name:   question.Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    0,
				},
				A: ip,
			}
		}
	}
	m.SetRcode(r, dns.RcodeSuccess)

	w.WriteMsg(m)
}
