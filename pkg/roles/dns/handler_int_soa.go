package dns

import (
	"fmt"

	"beryju.io/gravity/pkg/roles/dns/utils"
	"github.com/miekg/dns"
)

type SOAIntHandler struct {
	zone *Zone
}

func (sih *SOAIntHandler) Identifier() string {
	return "int_soa"
}

func (sih *SOAIntHandler) Handle(w *utils.FakeDNSWriter, r *utils.DNSRequest) *dns.Msg {
	if !sih.zone.Authoritative {
		return nil
	}
	m := new(dns.Msg)
	m.Authoritative = true
	m.Answer = []dns.RR{
		&dns.SOA{
			Hdr: dns.RR_Header{
				Name:   sih.zone.Name,
				Rrtype: dns.TypeSOA,
				Class:  dns.ClassINET,
				Ttl:    sih.zone.DefaultTTL,
			},
			Ns:      sih.zone.Name,
			Mbox:    fmt.Sprintf("root.%s", sih.zone.Name),
			Serial:  1337,
			Refresh: 600,
			Retry:   15,
			Expire:  5,
			Minttl:  sih.zone.DefaultTTL,
		},
	}
	for _, qs := range r.Question {
		if qs.Qtype == dns.TypeSOA {
			return m
		}
	}
	return nil
}
