package dns

import (
	"context"
	"strings"

	"beryju.io/gravity/pkg/roles/dns/utils"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type EtcdHandler struct {
	log *log.Entry
	z   *Zone
}

func NewEtcdHandler(z *Zone, config map[string]string) *EtcdHandler {
	eh := &EtcdHandler{
		z: z,
	}
	eh.log = z.log.WithField("handler", eh.Identifier())
	return eh
}

func (eh *EtcdHandler) Identifier() string {
	return "etcd"
}

func (eh *EtcdHandler) Handle(w *fakeDNSWriter, r *dns.Msg) *dns.Msg {
	m := new(dns.Msg)
	m.Authoritative = eh.z.Authoritative
	ctx := context.Background()
	for _, question := range r.Question {
		relRecordName := strings.TrimSuffix(question.Name, utils.EnsureLeadingPeriod(eh.z.Name))
		fullRecordKey := eh.z.inst.KV().Key(eh.z.etcdKey, relRecordName, dns.Type(question.Qtype).String())
		eh.log.WithField("key", fullRecordKey).Trace("fetching kv key")
		res, err := eh.z.inst.KV().Get(ctx, fullRecordKey, clientv3.WithPrefix())
		if err != nil || len(res.Kvs) < 1 {
			continue
		}
		for _, kv := range res.Kvs {
			rec := eh.z.recordFromKV(kv)
			ans := rec.ToDNS(question.Name, question.Qtype)
			if ans != nil {
				m.Answer = append(m.Answer, ans)
			}
		}
	}
	if len(m.Answer) < 1 {
		return nil
	}
	return m
}
