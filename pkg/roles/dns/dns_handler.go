package dns

import (
	"strings"

	"github.com/miekg/dns"
)

func (ro *DNSRole) handler(w dns.ResponseWriter, r *dns.Msg) {
	lastLongest := 0
	var longestZone *Zone

	for _, question := range r.Question {
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
	}
	if longestZone == nil {
		longestZone = ro.zones["."]
	}
	if longestZone == nil {
		go ro.createFallbackGlobalZone()
		ro.log.Error("no matching zone and no global zone")
		m := new(dns.Msg)
		m.SetRcode(r, dns.RcodeNameError)
		w.WriteMsg(m)
		return
	}
	ro.log.WithField("zone", longestZone.etcdKey).Trace("routing request to zone")
	longestZone.resolve(w, r)
}

func (r *DNSRole) createFallbackGlobalZone() {
	z := &Zone{
		Name: ".",
		HandlerConfigs: []map[string]string{
			{
				"type": "forward_ip",
				"to":   "8.8.8.8:53",
			},
		},
		inst: r.i,
	}
	z.put()
}
