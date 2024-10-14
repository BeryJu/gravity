package etcd

import (
	"strings"

	"beryju.io/gravity/pkg/roles/dns/handlers"
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
	log           *zap.Logger
	z             handlers.HandlerZoneContext
	LookupKeyFunc func(k *storage.Key, qname string, r *utils.DNSRequest) []dns.RR
}

func NewEtcdHandler(z handlers.HandlerZoneContext, config map[string]string) *EtcdHandler {
	eh := &EtcdHandler{
		z: z,
	}
	eh.LookupKeyFunc = func(k *storage.Key, qname string, r *utils.DNSRequest) []dns.RR {
		answers := []dns.RR{}
		es := sentry.TransactionFromContext(r.Context()).StartChild("gravity.dns.handler.etcd.get")
		defer es.Finish()
		key := k.String()
		eh.log.Debug("fetching kv key", zap.String("key", key))
		es.SetTag("gravity.dns.handler.etcd.key", key)
		res, err := eh.z.RoleInstance().KV().Get(r.Context(), key, clientv3.WithPrefix())
		if err != nil || len(res.Kvs) < 1 {
			return answers
		}
		for _, kv := range res.Kvs {
			rec, err := eh.z.RecordFromKV(kv)
			if err != nil {
				continue
			}
			ans := rec.ToDNS(qname)
			if ans != nil {
				answers = append(answers, ans)
			}
		}
		return answers
	}
	eh.log = z.Log().With(zap.String("handler", eh.Identifier()))
	return eh
}

func (eh *EtcdHandler) Identifier() string {
	return EtcdType
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
		wildcardKey := eh.z.RoleInstance().KV().Key(eh.z.EtcdKey(), strings.ToLower(wildcardName), dns.Type(question.Qtype).String())
		wildcardAns := eh.LookupKeyFunc(wildcardKey, question.Name, r)
		// If we do get an answer from this wildcard key, stop going further
		if len(wildcardAns) > 0 {
			return wildcardAns
		}
	}
	return []dns.RR{}
}

func (eh *EtcdHandler) handleSingleQuestion(question dns.Question, r *utils.DNSRequest) []dns.RR {
	answers := []dns.RR{}
	// Remove zone from query name
	relRecordName := strings.TrimSuffix(strings.ToLower(question.Name), strings.ToLower(eh.z.Name))
	if relRecordName == "" {
		// If the query name was the zone, the query should look for a record at the root
		relRecordName = types.DNSRoot
	} else {
		// Otherwise the relative record name still has a dot at the end which is not what we store
		// in the database
		relRecordName = strings.TrimSuffix(relRecordName, ".")
	}
	directRecordKey := eh.z.RoleInstance().KV().Key(
		eh.z.EtcdKey(),
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
	answers = append(answers, eh.LookupKeyFunc(
		directRecordKey,
		question.Name,
		r,
	)...)
	// If we don't find an answer for the direct key lookup, try a wildcard lookup
	if len(answers) < 1 {
		answers = append(answers, eh.findWildcard(r, relRecordName, question)...)
	}
	// If not explicitly looking for CNAME and
	// any of the answers are a CNAME, look those up too
	if question.Qtype != dns.TypeCNAME {
		for _, ans := range answers {
			if cn, ok := ans.(*dns.CNAME); ok {
				answers = append(answers, eh.handleSingleQuestion(dns.Question{
					Name:   cn.Header().Name,
					Qtype:  question.Qtype,
					Qclass: question.Qclass,
				}, r)...)
			}
		}
	}
	return answers
}

func (eh *EtcdHandler) Handle(w *utils.FakeDNSWriter, r *utils.DNSRequest) *dns.Msg {
	m := new(dns.Msg)
	m.Authoritative = eh.z.Authoritative
	uniqueQuestionNames := map[string]struct{}{}
	for _, question := range r.Question {
		m.Answer = append(m.Answer, eh.handleSingleQuestion(question, r)...)
		uniqueQuestionNames[question.Name] = struct{}{}
	}
	for un := range uniqueQuestionNames {
		// Look for CNAMEs
		relRecordName := strings.TrimSuffix(strings.ToLower(un), strings.ToLower(utils.EnsureLeadingPeriod(eh.z.Name)))
		cnames := eh.LookupKeyFunc(
			eh.z.RoleInstance().KV().Key(
				eh.z.EtcdKey(),
				strings.ToLower(relRecordName),
				dns.Type(dns.TypeCNAME).String(),
			),
			un,
			r,
		)
		if len(cnames) > 0 {
			m.Answer = append(m.Answer, cnames...)
			// For each cname, lookup the actual record
			for _, _cn := range cnames {
				cn := _cn.(*dns.CNAME)
				nq := dns.Question{
					Name: cn.Target,
				}
				m.Answer = append(m.Answer, eh.handleSingleQuestion(nq, r)...)
			}
		}
	}
	if len(m.Answer) < 1 {
		return nil
	}
	// If none of our answers match the types that were queried for, and we have more handlers in the chain
	// don't return an answer
	answerForType := false
	for _, a := range m.Answer {
		for _, q := range r.Question {
			if a.Header().Rrtype == q.Qtype {
				answerForType = true
			}
		}
	}
	if !answerForType && r.Meta().HasMoreHandlers {
		return nil
	}
	return m
}
