package dns

import (
	"context"
	"strings"

	"beryju.io/gravity/pkg/roles/dns/utils"
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

func (eh *EtcdHandler) Handle(w *utils.FakeDNSWriter, r *utils.DNSRequest) *dns.Msg {
	m := new(dns.Msg)
	m.Authoritative = eh.z.Authoritative
	ctx := context.Background()
	for _, question := range r.Question {
		relRecordName := strings.TrimSuffix(question.Name, utils.EnsureLeadingPeriod(eh.z.Name))
		fullRecordKey := eh.z.inst.KV().Key(eh.z.etcdKey, relRecordName, dns.Type(question.Qtype).String()).String()
		eh.log.Debug("fetching kv key", zap.String("key", fullRecordKey))
		es := sentry.StartSpan(r.Context(), "gravity.dns.handler.etcd.get")
		es.SetTag("gravity.dns.handler.etcd.key", fullRecordKey)
		res, err := eh.z.inst.KV().Get(ctx, fullRecordKey, clientv3.WithPrefix())
		es.Finish()
		if err != nil || len(res.Kvs) < 1 {
			continue
		}
		for _, kv := range res.Kvs {
			rec, err := eh.z.recordFromKV(kv)
			if err != nil {
				continue
			}
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
