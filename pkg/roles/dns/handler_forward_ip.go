package dns

import (
	"math/rand"
	"strconv"
	"strings"

	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/roles/dns/utils"
	"github.com/getsentry/sentry-go"
	"github.com/miekg/dns"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/structpb"
)

const IPForwarderType = "forward_ip"

type IPForwarderHandler struct {
	c         *dns.Client
	resolvers []string

	z        *ZoneContext
	log      *zap.Logger
	CacheTTL int
}

func NewIPForwarderHandler(z *ZoneContext, config map[string]*structpb.Value) *IPForwarderHandler {
	net, ok := config["net"]
	if !ok {
		net = structpb.NewStringValue("")
	}

	ipf := &IPForwarderHandler{
		z: z,
		c: &dns.Client{
			Net:     net.GetStringValue(),
			Timeout: types.DefaultUpstreamTimeout,
		},
		resolvers: strings.Split(config["to"].GetStringValue(), ";"),
	}
	ipf.log = z.log.With(zap.String("handler", ipf.Identifier()))

	rawTtl := config["cache_ttl"]
	cacheTtl, err := strconv.Atoi(rawTtl.GetStringValue())
	if err != nil && rawTtl.GetStringValue() != "" {
		ipf.log.Warn("failed to parse cache_ttl, defaulting to 0", zap.Error(err), zap.Any("config", config))
		cacheTtl = 0
	}
	ipf.CacheTTL = cacheTtl
	return ipf
}

func (ipf *IPForwarderHandler) Identifier() string {
	return IPForwarderType
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
		record.MxPreference = uint32(v.Preference)
	case *dns.SRV:
		record.Data = v.Target
		record.SrvPort = uint32(v.Port)
		record.SrvPriority = uint32(v.Priority)
		record.SrvWeight = uint32(v.Weight)
	}
	record.uid = strconv.Itoa(idx)
	err := record.put(r.Context(), int64(cacheTtl))
	if err != nil {
		ipf.log.Warn("failed to cache answer", zap.Error(err))
	}
}

func (ipf *IPForwarderHandler) pickResolver() string {
	resolver := ipf.resolvers[rand.Intn(len(ipf.resolvers))]
	if !strings.Contains(resolver, ":") {
		resolver += ":53"
	}
	return resolver
}

func (ipf *IPForwarderHandler) Handle(w *utils.FakeDNSWriter, r *utils.DNSRequest) *dns.Msg {
	if len(r.Question) < 1 {
		ipf.log.Error("No question")
		return nil
	}
	question := r.Question[0]
	fs := sentry.TransactionFromContext(r.Context()).StartChild("gravity.dns.handler.forward_ip.lookup")
	resolver := ipf.pickResolver()
	fs.SetTag("resolver", resolver)
	ipf.log.Debug("sending message to resolve", zap.String("resolver", resolver))
	m, rtt, err := ipf.c.ExchangeContext(r.Context(), r.Msg, resolver)
	ipf.log.Debug("dns rtt", zap.Duration("rtt", rtt))
	fs.Finish()

	if err != nil {
		ipf.log.Warn("failed to forward", zap.Error(err))
		m := new(dns.Msg)
		m.SetRcode(r.Msg, dns.RcodeServerFailure)
		return m
	}
	m.RecursionAvailable = true
	for idx, ans := range m.Answer {
		if ans == nil {
			continue
		}
		go ipf.cacheToEtcd(r, question, ans, idx)
	}
	if len(m.Answer) < 1 {
		m.SetRcode(r.Msg, dns.RcodeNameError)
	} else {
		m.SetRcode(r.Msg, dns.RcodeSuccess)
	}
	return m
}
