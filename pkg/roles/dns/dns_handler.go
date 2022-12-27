package dns

import (
	"context"
	"net"
	"strings"

	"github.com/miekg/dns"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.uber.org/zap"
)

// Find a zone for the given fqdn
func (r *Role) FindZone(fqdn string) *Zone {
	lastLongest := 0
	var longestZone *Zone
	for name, zone := range r.zones {
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

func (ro *Role) Handler(w dns.ResponseWriter, r *dns.Msg) {
	lastLongest := 0
	var longestZone *Zone

	_, span := otel.Tracer("").Start(context.TODO(), "gravity.roles.dns.request")
	defer span.End()

	var clientIP = ""
	switch addr := w.RemoteAddr().(type) {
	case *net.UDPAddr:
		clientIP = addr.IP.String()
	case *net.TCPAddr:
		clientIP = addr.IP.String()
	}
	span.SetAttributes(attribute.String("ip", clientIP))
	defer span.End()

	for _, question := range r.Question {
		span.SetAttributes(attribute.String("gravity.dns.query.type", dns.TypeToString[question.Qtype]))
		ro.zonesM.RLock()
		for name, zone := range ro.zones {
			// Zone doesn't have the correct suffix for the question
			if !strings.HasSuffix(question.Name, name) {
				continue
			}
			if len(name) > lastLongest {
				lastLongest = len(name)
				longestZone = zone
			}
		}
		ro.zonesM.RUnlock()
	}
	if longestZone == nil {
		longestZone = ro.zones["."]
	}
	if longestZone == nil {
		ro.log.Error("no matching zone and no global zone")
		m := new(dns.Msg)
		m.SetRcode(r, dns.RcodeNameError)
		w.WriteMsg(m)
		return
	}
	ro.log.Debug("routing request to zone", zap.String("zone", longestZone.etcdKey))
	span.SetAttributes(attribute.String("gravity.dns.zone", longestZone.Name))
	longestZone.resolve(w, r, span)
}
