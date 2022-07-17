package dns

import (
	"context"
	"errors"
	"math/rand"
	"net"
	"net/netip"
	"strconv"
	"strings"
	"time"

	"beryju.io/gravity/pkg/roles/dns/utils"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

type IPForwarderHandler struct {
	CacheTTL int

	r   *net.Resolver
	z   *Zone
	log *log.Entry
}

func NewIPForwarderHandler(z *Zone, config map[string]string) *IPForwarderHandler {
	ipf := &IPForwarderHandler{
		z: z,
	}
	ipf.log = z.log.WithField("handler", ipf.Identifier())

	rawTtl := config["cache_ttl"]
	cacheTtl, err := strconv.Atoi(rawTtl)
	if err != nil && rawTtl != "" {
		ipf.log.WithField("config", config).WithError(err).Warning("failed to parse cache_ttl, defaulting to 0")
		cacheTtl = 0
	}
	ipf.CacheTTL = cacheTtl

	forwarders := strings.Split(config["to"], ";")
	ipf.r = &net.Resolver{
		PreferGo: true,
		Dial: func(ctx context.Context, network, address string) (net.Conn, error) {
			d := net.Dialer{
				Timeout: time.Millisecond * 5,
			}
			resolver := forwarders[rand.Intn(len(forwarders))]
			if !strings.Contains(resolver, ":") {
				resolver += ":53"
			}
			return d.DialContext(ctx, network, resolver)
		},
	}
	return ipf
}

func (ipf *IPForwarderHandler) cacheToEtcd(query dns.Question, ans dns.RR, idx int) {
	if ans == nil {
		return
	}
	cacheTtl := ans.Header().Ttl
	// never cache if set to -1
	if ipf.CacheTTL == -1 {
		return
	}
	// If CacheTTL is set to -2 we don't expire the entry at all
	// this can be used to forward a specific zone and import all the records
	// on the fly
	if ipf.CacheTTL == -2 {
		cacheTtl = 0
	} else {
		// Try to set cache expiry based on TTL of answer
		// if no TTL set, default to CacheTTL, and if that's
		// not set, then don't cache at all
		if cacheTtl < 1 {
			cacheTtl = uint32(ipf.CacheTTL)
			if cacheTtl < 1 {
				return
			}
		}
	}
	name := strings.TrimSuffix(query.Name, utils.EnsureLeadingPeriod(ipf.z.Name))
	record := ipf.z.newRecord(name, dns.TypeToString[ans.Header().Rrtype])
	switch v := ans.(type) {
	case *dns.A:
		record.Data = v.A.String()
	case *dns.AAAA:
		record.Data = v.AAAA.String()
	case *dns.PTR:
		record.Data = v.Ptr
	case *dns.CNAME:
		record.Data = v.Target
	case *dns.MX:
		record.Data = v.Mx
		record.MXPreference = v.Preference
	case *dns.SRV:
		record.Data = v.Target
		record.SRVPort = v.Port
		record.SRVPriority = v.Priority
		record.SRVWeight = v.Weight
	}
	record.TTL = ans.Header().Ttl
	record.uid = strconv.Itoa(idx)
	err := record.put(context.Background(), int64(cacheTtl))
	if err != nil {
		ipf.log.WithError(err).Warning("failed to cache answer")
	}
}

func (ipf *IPForwarderHandler) Identifier() string {
	return "forward_ip"
}

func (ipf *IPForwarderHandler) Handle(w *utils.FakeDNSWriter, r *dns.Msg) *dns.Msg {
	if len(r.Question) < 1 {
		ipf.log.Error("No question")
		return nil
	}
	question := r.Question[0]
	ips, err := ipf.r.LookupHost(context.Background(), question.Name)
	m := new(dns.Msg)
	m.SetReply(r)

	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		m.SetRcode(r, dns.RcodeNameError)
		return m
	} else if err != nil {
		ipf.log.WithError(err).Warning("failed to forward")
		m.SetRcode(r, dns.RcodeServerFailure)
		return m
	}
	m.Answer = make([]dns.RR, 0)
	for idx, rawIp := range ips {
		ip, err := netip.ParseAddr(rawIp)
		if err != nil {
			ipf.log.WithError(err).Warning("failed to parse response IP")
			continue
		}
		var ans dns.RR
		if ip.Is6() && question.Qtype == dns.TypeAAAA {
			ans = &dns.AAAA{
				Hdr: dns.RR_Header{
					Name:   question.Name,
					Rrtype: dns.TypeAAAA,
					Class:  dns.ClassINET,
					Ttl:    0,
				},
				AAAA: net.ParseIP(ip.String()),
			}
		} else if ip.Is4() && question.Qtype == dns.TypeA {
			ans = &dns.A{
				Hdr: dns.RR_Header{
					Name:   question.Name,
					Rrtype: dns.TypeA,
					Class:  dns.ClassINET,
					Ttl:    0,
				},
				A: net.ParseIP(ip.String()),
			}
		}
		if ans == nil {
			continue
		}
		m.Answer = append(m.Answer, ans)
		go ipf.cacheToEtcd(question, ans, idx)
	}
	m.RecursionAvailable = true
	if len(m.Answer) < 1 {
		m.SetRcode(r, dns.RcodeNameError)
	} else {
		m.SetRcode(r, dns.RcodeSuccess)
	}
	return m
}
