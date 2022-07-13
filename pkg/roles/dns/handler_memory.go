package dns

import (
	"strings"

	"beryju.io/gravity/pkg/roles/dns/utils"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

type MemoryHandler struct {
	log *log.Entry
	z   *Zone
}

func NewMemoryHandler(z *Zone, config map[string]string) *MemoryHandler {
	eh := &MemoryHandler{
		z: z,
	}
	eh.log = z.log.WithField("handler", eh.Identifier())
	return eh
}

func (eh *MemoryHandler) Identifier() string {
	return "memory"
}

func (eh *MemoryHandler) Handle(w *fakeDNSWriter, r *dns.Msg) *dns.Msg {
	m := new(dns.Msg)
	m.Authoritative = eh.z.Authoritative
	for _, question := range r.Question {
		relRecordName := strings.TrimSuffix(question.Name, utils.EnsureLeadingPeriod(eh.z.Name))
		fullRecordKey := eh.z.inst.KV().Key(eh.z.etcdKey, relRecordName, dns.Type(question.Qtype).String())
		if recs, ok := eh.z.records[fullRecordKey]; ok {
			if len(recs) < 1 {
				continue
			}
			eh.log.WithField("key", fullRecordKey).Trace("got record in in-memory cache")
			for _, rec := range recs {
				ans := rec.ToDNS(question.Name, question.Qtype)
				if ans != nil {
					m.Answer = append(m.Answer, ans)
				}
			}
		}
	}
	if len(m.Answer) < 1 {
		return nil
	}
	return m
}
