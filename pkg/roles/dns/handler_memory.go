package dns

import (
	"strings"

	"beryju.io/gravity/pkg/roles/dns/utils"
	"github.com/miekg/dns"
	"go.uber.org/zap"
)

type MemoryHandler struct {
	log *zap.Logger
	z   *Zone
}

func NewMemoryHandler(z *Zone, config map[string]string) *MemoryHandler {
	eh := &MemoryHandler{
		z: z,
	}
	eh.log = z.log.With(zap.String("handler", eh.Identifier()))
	return eh
}

func (eh *MemoryHandler) Identifier() string {
	return "memory"
}

func (eh *MemoryHandler) Handle(w *utils.FakeDNSWriter, r *dns.Msg) *dns.Msg {
	m := new(dns.Msg)
	m.Authoritative = eh.z.Authoritative
	for _, question := range r.Question {
		relRecordName := strings.TrimSuffix(question.Name, utils.EnsureLeadingPeriod(eh.z.Name))
		fullRecordKey := eh.z.inst.KV().Key(eh.z.etcdKey, relRecordName, dns.Type(question.Qtype).String()).String()
		eh.z.recordsSync.RLock()
		recs, ok := eh.z.records[fullRecordKey]
		eh.z.recordsSync.RUnlock()
		if ok {
			if len(recs) < 1 {
				continue
			}
			eh.log.Debug("got record in in-memory cache", zap.String("key", fullRecordKey))
			for _, rec := range recs {
				ans := rec.ToDNS(question.Name, question.Qtype)
				if ans != nil {
					m.Answer = append(m.Answer, ans)
				}
			}
		}
	}
	if len(m.Answer) < 1 {
		if eh.z.Authoritative {
			return eh.z.soa()
		}
		return nil
	}
	return m
}
