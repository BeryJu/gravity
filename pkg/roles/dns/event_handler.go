package dns

import (
	"net/netip"
	"strings"

	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/roles/dns/utils"
	"github.com/miekg/dns"
	"go.uber.org/zap"
)

func (r *Role) eventHandlerDHCPLeaseGiven(ev *roles.Event) {
	if ev.Payload.Data["hostname"] == "" {
		return
	}
	// Forward record
	r.eventCreateForward(ev)
	// Reverse record
	r.eventCreateReverse(ev)
}

func (r *Role) eventCreateForward(ev *roles.Event) {
	hostname := ev.Payload.Data["hostname"].(string)
	fqdn := ev.Payload.Data["fqdn"].(string)
	forwardZone := r.FindZone(fqdn)
	if forwardZone == nil {
		r.log.Debug("No zone found for hostname", zap.Any("event", ev), zap.String("fqdn", fqdn))
		return
	}

	rawAddr := ev.Payload.Data["address"].(string)
	ip, err := netip.ParseAddr(rawAddr)
	if err != nil {
		r.log.Warn("failed to parse address to add dns record", zap.Error(err))
		return
	}
	var rec *RecordContext
	if ip.Is4() {
		rec = forwardZone.newRecord(hostname, types.DNSRecordTypeA)
	} else {
		rec = forwardZone.newRecord(hostname, types.DNSRecordTypeA)
	}
	rec.Data = ip.String()
	rec.Ttl = forwardZone.DefaultTTL
	err = rec.put(ev.Context, 0, ev.Payload.RelatedObjectOptions...)
	if err != nil {
		r.log.Warn("failed to save dns record", zap.Error(err))
	}
}

func (r *Role) eventCreateReverse(ev *roles.Event) {
	fqdn := ev.Payload.Data["fqdn"].(string)
	rawAddr := ev.Payload.Data["address"].(string)
	ip, err := netip.ParseAddr(rawAddr)
	if err != nil {
		r.log.Warn("failed to parse address to add dns record", zap.Error(err))
		return
	}

	rev, err := dns.ReverseAddr(ip.String())
	if err != nil {
		r.log.Warn("failed to get reverse of ip", zap.Error(err))
		return
	}

	forwardZone := r.FindZone(rev)
	if forwardZone == nil {
		r.log.Debug("No zone found for hostname", zap.Any("event", ev), zap.String("fqdn", fqdn))
		return
	}

	relName := strings.TrimSuffix(rev, utils.EnsureLeadingPeriod(forwardZone.Name))
	rec := forwardZone.newRecord(relName, types.DNSRecordTypePTR)
	rec.Data = fqdn
	rec.Ttl = forwardZone.DefaultTTL
	err = rec.put(ev.Context, 0, ev.Payload.RelatedObjectOptions...)
	if err != nil {
		r.log.Warn("failed to save dns record", zap.Error(err))
	}
}
