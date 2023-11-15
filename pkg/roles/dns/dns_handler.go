package dns

import (
	"context"
	"net"
	"strings"

	"beryju.io/gravity/pkg/roles/dns/utils"
	"github.com/getsentry/sentry-go"
	"github.com/miekg/dns"
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

	span := sentry.StartTransaction(
		context.TODO(),
		"gravity.dns.request",
	)
	var clientIP = ""
	switch addr := w.RemoteAddr().(type) {
	case *net.UDPAddr:
		clientIP = addr.IP.String()
	case *net.TCPAddr:
		clientIP = addr.IP.String()
	}
	hub := sentry.GetHubFromContext(span.Context())
	if hub == nil {
		hub = sentry.CurrentHub()
	}
	hub.Scope().SetUser(sentry.User{
		IPAddress: clientIP,
	})
	defer span.Finish()

	for _, question := range r.Question {
		span.SetTag("gravity.dns.query.type", dns.TypeToString[question.Qtype])
		ro.zonesM.RLock()
		for name, zone := range ro.zones {
			// Zone doesn't have the correct suffix for the question
			if !strings.HasSuffix(strings.ToLower(question.Name), strings.ToLower(name)) {
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
		err := w.WriteMsg(m)
		if err != nil {
			ro.log.Warn("failed to send answer", zap.Error(err))
		}
		return
	}
	ro.log.Debug("routing request to zone", zap.String("zone", longestZone.etcdKey))
	span.SetTag("gravity.dns.zone", longestZone.Name)
	longestZone.resolve(w, utils.NewRequest(r, span.Context()), span)
}
