package dns

import (
	"context"
	"strings"

	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
)

type EtcdHandler struct {
	log *log.Entry
	z   Zone
}

func NewEtcdHandler(z Zone, config map[string]string) *EtcdHandler {
	return &EtcdHandler{
		log: log.WithField("handler", "etcd"),
		z:   z,
	}
}

func (eh *EtcdHandler) Handle(w *fakeDNSWriter, r *dns.Msg) *dns.Msg {
	m := new(dns.Msg)
	m.Authoritative = eh.z.Authoritative
	ctx := context.Background()
	for _, question := range r.Question {
		relRecordName := strings.TrimSuffix(question.Name, "."+eh.z.Name)
		fullRecordKey := eh.z.inst.GetKV().Key(eh.z.etcdKey, relRecordName, dns.Type(question.Qtype).String())
		// TODO: Optimise this
		res, err := eh.z.inst.GetKV().Get(ctx, fullRecordKey)
		if err != nil || len(res.Kvs) < 1 {
			continue
		}
		for _, key := range res.Kvs {
			ans := eh.z.kvToDNS(question.Name, key, question.Qtype)
			if ans != nil {
				m.Answer = append(m.Answer, ans)
			}
		}
	}
	if len(m.Answer) < 1 {
		m.SetRcode(r, dns.RcodeNameError)
	}
	return m
}
