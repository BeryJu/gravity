package dns

import (
	"context"
	"fmt"

	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/roles/dns/utils"
	"github.com/miekg/dns"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type SOAIntHandler struct {
	zone *Zone
}

func (sih *SOAIntHandler) Identifier() string {
	return "int_soa"
}

func (sih *SOAIntHandler) getOrCreateSOA(ctx context.Context) *Record {
	soaKey := sih.zone.inst.KV().Key(sih.zone.etcdKey, types.DNSRootRecord, "SOA")
	res, err := sih.zone.inst.KV().Get(ctx, soaKey.String(), clientv3.WithPrefix())
	if err == nil && len(res.Kvs) > 0 {
		rec, err := sih.zone.recordFromKV(res.Kvs[0])
		if err == nil {
			return rec
		}
		sih.zone.log.Warn("failed to parse SOA record from etcd", zap.Error(err))
	}
	rec := sih.zone.newRecord(types.DNSRootRecord, "SOA")
	rec.uid = "0"
	rec.TTL = sih.zone.DefaultTTL
	rec.Data = sih.zone.Name
	rec.SOAMbox = fmt.Sprintf("root.%s", sih.zone.Name)
	rec.SOASerial = 1337
	rec.SOARefresh = 600
	rec.SOARetry = 15
	rec.SOAExpire = 5
	if err := rec.put(ctx, 0); err != nil {
		sih.zone.log.Warn("failed to persist default SOA record", zap.Error(err))
	}
	return rec
}

func (sih *SOAIntHandler) Handle(w *utils.FakeDNSWriter, r *utils.DNSRequest) *dns.Msg {
	if !sih.zone.Authoritative {
		return nil
	}
	for _, qs := range r.Question {
		if qs.Qtype != dns.TypeSOA {
			continue
		}
		rec := sih.getOrCreateSOA(r.Context())
		m := new(dns.Msg)
		m.Authoritative = true
		m.Answer = []dns.RR{rec.ToDNS(sih.zone.Name)}
		return m
	}
	return nil
}
