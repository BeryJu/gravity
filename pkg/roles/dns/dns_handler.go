package dns

import (
	"context"
	"net"
	"strings"

	"beryju.io/gravity/pkg/o11y"
	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/roles/dns/utils"
	"github.com/getsentry/sentry-go"
	"github.com/miekg/dns"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
)

const IterationMax = 20

// Find a zone for the given fqdn
func (r *Role) FindZone(fqdn string) *Zone {
	lastLongest := 0
	var longestZone *Zone
	for name, zone := range r.zones.IterRelativeKey() {
		// Zone doesn't have the correct suffix for the question
		if !strings.HasSuffix(fqdn, name) {
			continue
		}
		if len(name) > lastLongest {
			lastLongest = len(name)
			longestZone = zone
		}
	}
	return longestZone
}

func (ro *Role) handler(w dns.ResponseWriter, r *dns.Msg) {
	ctx, span := o11y.Tracer.Start(context.TODO(), "gravity.dns.request")
	clientIP := ""
	switch addr := w.RemoteAddr().(type) {
	case *net.UDPAddr:
		clientIP = addr.IP.String()
	case *net.TCPAddr:
		clientIP = addr.IP.String()
	}
	hub := sentry.GetHubFromContext(ctx)
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	hub.Scope().SetUser(sentry.User{
		IPAddress: clientIP,
	})
	defer span.End()
	req := utils.NewRequest(r, ctx, utils.DNSRoutingMeta{})
	ro.rootHandler(w, req)
}

func (ro *Role) rootHandler(w dns.ResponseWriter, r *utils.DNSRequest) {
	if r.Iteration() > IterationMax {
		ro.log.Error("exceeded maximum iteration count")
		m := new(dns.Msg)
		m.SetRcode(r.Msg, dns.RcodeNameError)
		err := w.WriteMsg(m)
		if err != nil {
			ro.log.Warn("failed to send answer", zap.Error(err))
		}
		return
	}
	lastLongest := 0
	var longestZone *Zone
	span := trace.SpanFromContext(r.Context())
	for _, question := range r.Question {
		span.SetName(question.Name)
		span.SetAttributes(
			attribute.String("http.request.method", dns.TypeToString[question.Qtype]),
			attribute.String("gravity.dns.query.type", dns.TypeToString[question.Qtype]),
		)
		for name, zone := range ro.zones.IterRelativeKey() {
			// Zone doesn't have the correct suffix for the question
			if !strings.HasSuffix(strings.ToLower(question.Name), strings.ToLower(name)) {
				continue
			}
			if len(name) > lastLongest {
				lastLongest = len(name)
				longestZone = zone
			}
		}
	}
	if longestZone == nil {
		longestZone, _ = ro.zones.GetPrefix(types.DNSRootZone)
	}
	if longestZone == nil {
		ro.log.Warn("no matching zone and no global zone")
		m := new(dns.Msg)
		m.SetRcode(r.Msg, dns.RcodeNameError)
		err := w.WriteMsg(m)
		if err != nil {
			ro.log.Warn("failed to send answer", zap.Error(err))
		}
		return
	}
	ro.log.Debug("routing request to zone", zap.String("zone", longestZone.etcdKey))
	span.SetAttributes(attribute.String("gravity.dns.zone", longestZone.Name))
	longestZone.resolve(w, r, span)
}
