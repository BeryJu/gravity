package dns

import (
	"context"
	"math/rand"
	"net"
	"net/netip"
	"strings"
	"time"

	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

type IPForwarderHandler struct {
	r   *net.Resolver
	log *log.Entry
}

func NewIPForwarderHandler(z Zone, config map[string]string) *IPForwarderHandler {
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

	return &IPForwarderHandler{
		r:   r,
		log: z.log.WithField("handler", "forward_ip"),
	}
}

func (ipf *IPForwarderHandler) Handle(w *fakeDNSWriter, r *dns.Msg) *dns.Msg {
	if len(r.Question) < 1 {
		ipf.log.Error("No question")
		return nil
	}
	question := r.Question[0]
	ips, err := ipf.r.LookupHost(context.Background(), question.Name)
	m := new(dns.Msg)
	m.SetReply(r)

	if err != nil {
		ipf.log.WithError(err).Warning("failed to forward")
		m.SetRcode(r, dns.RcodeServerFailure)
		return nil
	}
	m.Answer = make([]dns.RR, len(ips))
	for idx, rawIp := range ips {
		ip, err := netip.ParseAddr(rawIp)
		if err != nil {
			ipf.log.WithError(err).Warning("failed to parse response IP")
			continue
		}
		if ip.Is6() {
			m.Answer[idx] = &dns.AAAA{
				Hdr: dns.RR_Header{
					Name:   question.Name,
					Rrtype: dns.TypeAAAA,
					Class:  dns.ClassINET,
					Ttl:    0,
				},
				AAAA: net.ParseIP(ip.String()),
			}
		} else {
			m.Answer[idx] = &dns.A{
				Hdr: dns.RR_Header{
					Name:   question.Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    0,
				},
				A: net.ParseIP(ip.String()),
			}
		}
	}
	if len(m.Answer) < 1 {
		m.SetRcode(r, dns.RcodeNameError)
	} else {
		m.SetRcode(r, dns.RcodeSuccess)
	}
	return m
}
