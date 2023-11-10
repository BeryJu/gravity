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

func (r *Role) eventHandlerDHCPLeasePut(ev *roles.Event) {
	if ev.Payload.Data["hostname"] == "" {
		return
	}
	// Forward record
	r.eventHandlerCreateForward(ev)
	// Reverse record
	r.eventHandlerCreateReverse(ev)
}

// eventHandlerCreateForward Event handler for `roles.dns.record.create_forward`
// also used by the DHCP lease put event handler
// requires these payload data attributes:
// - hostname for the device name
// - fqdn for the dns zone to put the record in
// - identifier for the record UID
// - address for the actual address
func (r *Role) eventHandlerCreateForward(ev *roles.Event) {
	hostname := ev.Payload.Data["hostname"].(string)
	fqdn := ev.Payload.Data["fqdn"].(string)
	identifier := ev.Payload.Data["identifier"].(string)
	forwardZone := r.FindZone(utils.EnsureTrailingPeriod(fqdn))
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
	var rec *Record
	if ip.Is4() {
		rec = forwardZone.newRecord(hostname, types.DNSRecordTypeA)
	} else {
		rec = forwardZone.newRecord(hostname, types.DNSRecordTypeAAAA)
	}
	rec.Data = ip.String()
	rec.uid = identifier
	rec.TTL = forwardZone.DefaultTTL
	err = rec.put(ev.Context, 0, ev.Payload.RelatedObjectOptions...)
	if err != nil {
		r.log.Warn("failed to save dns record", zap.Error(err))
		return
	}
	r.log.Debug("put record", zap.String("record", rec.Name), zap.String("zone", forwardZone.Name))
}

// eventHandlerCreateReverse Event handler for `roles.dns.record.create_reverse`
// also used by the DHCP lease put event handler
// requires these payload data attributes:
// - hostname for the device name
// - fqdn for the dns zone to put the record in
// - identifier for the record UID
// - address for the actual address
func (r *Role) eventHandlerCreateReverse(ev *roles.Event) {
	fqdn := ev.Payload.Data["fqdn"].(string)
	rawAddr := ev.Payload.Data["address"].(string)
	identifier := ev.Payload.Data["identifier"].(string)
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

	reverseZone := r.FindZone(utils.EnsureTrailingPeriod(rev))
	if reverseZone == nil {
		r.log.Debug("No zone found for hostname", zap.Any("event", ev), zap.String("fqdn", fqdn))
		return
	}

	relName := strings.TrimSuffix(rev, utils.EnsureLeadingPeriod(reverseZone.Name))
	rec := reverseZone.newRecord(relName, types.DNSRecordTypePTR)
	rec.uid = identifier
	rec.Data = fqdn
	rec.TTL = reverseZone.DefaultTTL
	err = rec.put(ev.Context, 0, ev.Payload.RelatedObjectOptions...)
	if err != nil {
		r.log.Warn("failed to save dns record", zap.Error(err))
		return
	}
	r.log.Debug("put record", zap.String("record", rec.Name), zap.String("zone", reverseZone.Name))
}
