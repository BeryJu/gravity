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

const CNameRecursionMax = 10

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
		ans := rec.ToDNS(question.Name)
		if ans != nil {
			answers = append(answers, ans)
		}
	}
	return answers
}

func (eh *EtcdHandler) findWildcard(r *utils.DNSRequest, relRecordName string, question dns.Question) []dns.RR {
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
			return wildcardAns
		}
	}
	return []dns.RR{}
}

func (eh *EtcdHandler) handleSingleQuestion(question dns.Question, r *utils.DNSRequest, level int) []dns.RR {
	answers := []dns.RR{}
	if level >= CNameRecursionMax {
		eh.log.Info("stopping cname recursion at level", zap.Int("max", CNameRecursionMax))
		return answers
	}
	relRecordName := strings.TrimSuffix(strings.ToLower(question.Name), strings.ToLower(utils.EnsureLeadingPeriod(eh.z.Name)))
	directRecordKey := eh.z.inst.KV().Key(
		eh.z.etcdKey,
		strings.ToLower(relRecordName),
	)
	if question.Qtype != dns.TypeNone {
		// If we're looking for a specific key, include that in the etcd key
		directRecordKey = directRecordKey.Add(dns.Type(question.Qtype).String())
	} else {
		// otherwise turn the query into a prefix
		// this is most likely to happen due to a CNAME query, where we don't know if its
		// name points to an A, AAAA or anything else (realistically it should be either
		// of those)
		directRecordKey = directRecordKey.Prefix(true)
	}
	// Look for direct matches first
	answers = append(answers, eh.lookupKey(
		directRecordKey,
		question,
		r,
	)...)
	// Look for CNAMEs
	// This section has room for optimization, as for we currently check for the existence
	// of CNAMEs for each question, which is a bit pointless if there are questions for foo.bar/A
	// and foo.bar/AAAA, and we also do the recursion for each question
	cnames := eh.lookupKey(
		eh.z.inst.KV().Key(
			eh.z.etcdKey,
			strings.ToLower(relRecordName),
			dns.Type(dns.TypeCNAME).String(),
		),
		question,
		r,
	)
	if len(cnames) > 0 {
		answers = append(answers, cnames...)
		// For each cname, lookup the actual record
		for _, _cn := range cnames {
			cn := _cn.(*dns.CNAME)
			nq := dns.Question{
				Name: cn.Target,
			}
			answers = append(answers, eh.handleSingleQuestion(nq, r, level+1)...)
		}
	}
	// If we don't find an answer for the direct key lookup, try a wildcard lookup
	if len(answers) < 1 {
		answers = append(answers, eh.findWildcard(r, relRecordName, question)...)
	}
	return answers
}

func (eh *EtcdHandler) Handle(w *utils.FakeDNSWriter, r *utils.DNSRequest) *dns.Msg {
	m := new(dns.Msg)
	m.Authoritative = eh.z.Authoritative
	for _, question := range r.Question {
		m.Answer = append(m.Answer, eh.handleSingleQuestion(question, r, 0)...)
	}
	if len(m.Answer) < 1 {
		return nil
	}
	return m
}
