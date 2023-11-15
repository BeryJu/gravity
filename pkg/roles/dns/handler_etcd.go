package dns

import (
	"strings"

	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/roles/dns/utils"
	"beryju.io/gravity/pkg/storage"
	"github.com/getsentry/sentry-go"
	"github.com/miekg/dns"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

const EtcdType = "etcd"

type EtcdHandler struct {
	log *zap.Logger
	z   *Zone
}

func NewEtcdHandler(z *Zone, config map[string]string) *EtcdHandler {
	eh := &EtcdHandler{
		z: z,
	}
	eh.log = z.log.With(zap.String("handler", eh.Identifier()))
	return eh
}

func (eh *EtcdHandler) Identifier() string {
	return EtcdType
}

// lookupKey Lookup direct key and fetch all UID entries below it
func (eh *EtcdHandler) lookupKey(k *storage.Key, question dns.Question, r *utils.DNSRequest) []dns.RR {
	answers := []dns.RR{}
	es := sentry.TransactionFromContext(r.Context()).StartChild("gravity.dns.handler.etcd.get")
	defer es.Finish()
	key := k.String()
	eh.log.Debug("fetching kv key", zap.String("key", key))
	es.SetTag("gravity.dns.handler.etcd.key", key)
	res, err := eh.z.inst.KV().Get(r.Context(), key, clientv3.WithPrefix())
	if err != nil || len(res.Kvs) < 1 {
		return answers
	}
	for _, kv := range res.Kvs {
		rec, err := eh.z.recordFromKV(kv)
		if err != nil {
			continue
		}
		ans := rec.ToDNS(question.Name, question.Qtype)
		if ans != nil {
			answers = append(answers, ans)
		}
	}
	return answers
}

func (eh *EtcdHandler) Handle(w *utils.FakeDNSWriter, r *utils.DNSRequest) *dns.Msg {
	m := new(dns.Msg)
	m.Authoritative = eh.z.Authoritative
	for _, question := range r.Question {
		relRecordName := strings.TrimSuffix(strings.ToLower(question.Name), strings.ToLower(utils.EnsureLeadingPeriod(eh.z.Name)))
		fullRecordKey := eh.z.inst.KV().Key(eh.z.etcdKey, strings.ToLower(relRecordName), dns.Type(question.Qtype).String())
		ans := eh.lookupKey(fullRecordKey, question, r)
		// If we don't find an answer for the direct key lookup, try a wildcard lookup
		if len(ans) < 1 {
			// Assuming the question is foo.bar.baz and the zone is baz,
			// we'll try replacing all names from left to right by starts and query with that
			wildcardName := relRecordName
			parts := strings.Split(relRecordName, ".")
			for _, part := range parts {
				// Replace the current dot part with a wildcard (make sure to only replace 1 occurrence,
				// since we replace from left to right)
				wildcardName = strings.Replace(wildcardName, part, types.DNSWildcard, 1)
				wildcardKey := eh.z.inst.KV().Key(eh.z.etcdKey, strings.ToLower(wildcardName), dns.Type(question.Qtype).String())
				wildcardAns := eh.lookupKey(wildcardKey, question, r)
				// If we do get an answer from this wildcard key, stop going further
				if len(wildcardAns) > 0 {
					m.Answer = append(m.Answer, wildcardAns...)
					break
				}
			}
		} else {
			m.Answer = append(m.Answer, ans...)
		}
	}
	if len(m.Answer) < 1 {
		return nil
	}
	return m
}
