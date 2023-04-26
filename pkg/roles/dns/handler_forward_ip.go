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
	"github.com/getsentry/sentry-go"
	"github.com/miekg/dns"
	"go.uber.org/zap"
)

const IPForwarderType = "forward_ip"

type IPForwarderHandler struct {
	r        *net.Resolver
	z        *Zone
	log      *zap.Logger
	CacheTTL int
}

func NewIPForwarderHandler(z *Zone, config map[string]string) *IPForwarderHandler {
	ipf := &IPForwarderHandler{
		z: z,
	}
	ipf.log = z.log.With(zap.String("handler", ipf.Identifier()))

	rawTtl := config["cache_ttl"]
	cacheTtl, err := strconv.Atoi(rawTtl)
	if err != nil && rawTtl != "" {
		ipf.log.Warn("failed to parse cache_ttl, defaulting to 0", zap.Error(err), zap.Any("config", config))
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

func (ipf *IPForwarderHandler) cacheToEtcd(r *utils.DNSRequest, query dns.Question, ans dns.RR, idx int) {
	cs := sentry.TransactionFromContext(r.Context()).StartChild("gravity.dns.handler.forward_ip.cache")
	cs.SetTag("gravity.dns.handler.forward_ip.cache.query", query.String())
	cs.SetTag("gravity.dns.handler.forward_ip.cache.ans", ans.String())
	defer cs.Finish()
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
	case *dns.TXT:
		record.Data = strings.Join(v.Txt, TXTSeparator)
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
	err := record.put(r.Context(), int64(cacheTtl))
	if err != nil {
		ipf.log.Warn("failed to cache answer", zap.Error(err))
	}
}

func (ipf *IPForwarderHandler) Identifier() string {
	return IPForwarderType
}

func (ipf *IPForwarderHandler) Handle(w *utils.FakeDNSWriter, r *utils.DNSRequest) *dns.Msg {
	if len(r.Question) < 1 {
		ipf.log.Error("No question")
		return nil
	}
	question := r.Question[0]
	fs := sentry.TransactionFromContext(r.Context()).StartChild("gravity.dns.handler.forward_ip.lookup")
	ips, err := ipf.r.LookupHost(r.Context(), question.Name)
	fs.Finish()
	m := new(dns.Msg)
	m.SetReply(r.Msg)

	var dnsErr *net.DNSError
	if errors.As(err, &dnsErr) {
		m.SetRcode(r.Msg, dns.RcodeNameError)
		return m
	} else if err != nil {
		ipf.log.Warn("failed to forward", zap.Error(err))
		m.SetRcode(r.Msg, dns.RcodeServerFailure)
		return m
	}
	m.Answer = make([]dns.RR, 0)
	for idx, rawIp := range ips {
		ip, err := netip.ParseAddr(rawIp)
		if err != nil {
			ipf.log.Warn("failed to parse response IP", zap.Error(err))
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
		go ipf.cacheToEtcd(r, question, ans, idx)
	}
	m.RecursionAvailable = true
	if len(m.Answer) < 1 {
		m.SetRcode(r.Msg, dns.RcodeNameError)
	} else {
		m.SetRcode(r.Msg, dns.RcodeSuccess)
	}
	return m
}
