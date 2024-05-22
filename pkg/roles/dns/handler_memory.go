package dns

import (
	"strings"

	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/roles/dns/utils"
	"beryju.io/gravity/pkg/storage"
	"github.com/getsentry/sentry-go"
	"github.com/miekg/dns"
	"go.uber.org/zap"
)

const MemoryType = "memory"

type MemoryHandler struct {
	log *zap.Logger
	z   *Zone
}

func NewMemoryHandler(z *Zone, config map[string]string) *MemoryHandler {
	mh := &MemoryHandler{
		z: z,
	}
	mh.log = z.log.With(zap.String("handler", mh.Identifier()))
	return mh
}

func (mh *MemoryHandler) Identifier() string {
	return MemoryType
}

// lookupKey Lookup direct key and fetch all UID entries below it
func (mh *MemoryHandler) lookupKey(k *storage.Key, qname string) []dns.RR {
	mh.z.recordsSync.RLock()
	recs, ok := mh.z.records[k.String()]
	mh.z.recordsSync.RUnlock()
	answers := []dns.RR{}
	if !ok {
		return answers
	}
	for _, rec := range recs {
		ans := rec.ToDNS(qname)
		if ans != nil {
			answers = append(answers, ans)
		}
	}
	return answers
}

func (mh *MemoryHandler) handleSingleQuestion(question dns.Question) []dns.RR {
	answers := []dns.RR{}
	relRecordName := strings.TrimSuffix(strings.ToLower(question.Name), strings.ToLower(utils.EnsureLeadingPeriod(mh.z.Name)))
	directRecordKey := mh.z.inst.KV().Key(
		mh.z.etcdKey,
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
	answers = append(answers, mh.lookupKey(
		directRecordKey,
		question.Name,
	)...)
	// If we don't find an answer for the direct key lookup, try a wildcard lookup
	if len(answers) < 1 {
		answers = append(answers, mh.findWildcard(relRecordName, question)...)
	}
	return answers
}

func (mh *MemoryHandler) findWildcard(relRecordName string, question dns.Question) []dns.RR {
	// Assuming the question is foo.bar.baz and the zone is baz,
	// we'll try replacing all names from left to right by starts and query with that
	wildcardName := relRecordName
	parts := strings.Split(relRecordName, ".")
	for _, part := range parts {
		// Replace the current dot part with a wildcard (make sure to only replace 1 occurrence,
		// since we replace from left to right)
		wildcardName = strings.Replace(wildcardName, part, types.DNSWildcard, 1)
		wildcardKey := mh.z.inst.KV().Key(mh.z.etcdKey, strings.ToLower(wildcardName), dns.Type(question.Qtype).String())
		wildcardAns := mh.lookupKey(wildcardKey, question.Name)
		// If we do get an answer from this wildcard key, stop going further
		if len(wildcardAns) > 0 {
			return wildcardAns
		}
	}
	return []dns.RR{}
}

func (mh *MemoryHandler) Handle(w *utils.FakeDNSWriter, r *utils.DNSRequest) *dns.Msg {
	m := new(dns.Msg)
	m.Authoritative = mh.z.Authoritative
	ms := sentry.TransactionFromContext(r.Context()).StartChild("gravity.dns.handler.memory.get")
	defer ms.Finish()
	uniqueQuestionNames := map[string]struct{}{}
	for _, question := range r.Question {
		m.Answer = append(m.Answer, mh.handleSingleQuestion(question)...)
		uniqueQuestionNames[question.Name] = struct{}{}
	}
	for _, question := range r.Question {
		relRecordName := strings.TrimSuffix(strings.ToLower(question.Name), strings.ToLower(utils.EnsureLeadingPeriod(mh.z.Name)))
		fullRecordKey := mh.z.inst.KV().Key(mh.z.etcdKey, strings.ToLower(relRecordName), dns.Type(question.Qtype).String()).String()
		mh.z.recordsSync.RLock()
		recs, ok := mh.z.records[fullRecordKey]
		mh.z.recordsSync.RUnlock()
		if ok {
			if len(recs) < 1 {
				continue
			}
			mh.log.Debug("got record in in-memory cache", zap.String("key", fullRecordKey))
			for _, rec := range recs {
				ans := rec.ToDNS(question.Name)
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
