package dns

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/roles/dns/utils"
	tsdbTypes "beryju.io/gravity/pkg/roles/tsdb/types"
	"github.com/getsentry/sentry-go"
	"github.com/miekg/dns"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"
)

const (
	DefaultTTL = 3600
)

type ZoneContext struct {
	*types.Zone
	inst            roles.Instance
	records         map[string]map[string]*Record
	recordsWatchCtx context.CancelFunc
	log             *zap.Logger
	etcdKey         string
	h               []Handler
	recordsSync     sync.RWMutex
}

func (z *ZoneContext) soa() *dns.Msg {
	m := new(dns.Msg)
	m.Authoritative = z.Authoritative
	m.Ns = []dns.RR{
		&dns.SOA{
			Hdr: dns.RR_Header{
				Name:   z.Name,
				Rrtype: dns.TypeSOA,
				Class:  dns.ClassINET,
				Ttl:    z.DefaultTTL,
			},
			Ns:      z.Name,
			Mbox:    fmt.Sprintf("root.%s", z.Name),
			Serial:  1337,
			Refresh: 600,
			Retry:   15,
			Expire:  5,
			Minttl:  z.DefaultTTL,
		},
	}
	return m
}

func (z *ZoneContext) resolveUpdateMetrics(dur time.Duration, q *utils.DNSRequest, h Handler, rep *dns.Msg) {
	for _, question := range q.Question {
		dnsQueries.WithLabelValues(
			dns.TypeToString[question.Qtype],
			h.Identifier(),
			z.Name,
		).Inc()
		dnsQueryDuration.WithLabelValues(
			dns.TypeToString[question.Qtype],
			h.Identifier(),
			z.Name,
		).Observe(float64(dur.Milliseconds()))
		go z.inst.DispatchEvent(tsdbTypes.EventTopicTSDBInc, roles.NewEvent(
			context.Background(),
			map[string]interface{}{
				"key": z.inst.KV().Key(
					types.KeyRole,
					h.Identifier(),
				).String(),
				"default": tsdbTypes.Metric{
					ResetOnWrite: true,
				},
			},
		))
	}
}

func (z *ZoneContext) resolve(w dns.ResponseWriter, r *utils.DNSRequest, span *sentry.Span) {
	for _, handler := range z.h {
		ss := span.StartChild("gravity.dns.request.handler")
		ss.Description = handler.Identifier()
		z.log.Debug("sending request to handler", zap.String("handler", handler.Identifier()))
		ss.SetTag("gravity.dns.handler", handler.Identifier())
		// Create new request for handler with the correct context
		hr := utils.NewRequest(r.Msg, ss.Context())

		handlerReply := handler.Handle(utils.NewFakeDNSWriter(w), hr)
		ss.Finish()

		if handlerReply != nil {
			z.log.Debug("returning reply from handler", zap.String("handler", handler.Identifier()))
			handlerReply.SetReply(r.Msg)
			w.WriteMsg(handlerReply)
			z.resolveUpdateMetrics(ss.EndTime.Sub(ss.StartTime), r, handler, handlerReply)
			return
		}
		z.log.Debug("no reply, trying next handler", zap.String("handler", handler.Identifier()))
	}
	if z.Authoritative {
		soa := z.soa()
		soa.SetReply(r.Msg)
		w.WriteMsg(soa)
		return
	}
	z.log.Debug("no handler has a reply, fallback back to NX")
	fallback := new(dns.Msg)
	fallback.SetReply(r.Msg)
	fallback.SetRcode(r.Msg, dns.RcodeNameError)
	w.WriteMsg(fallback)
}

func (r *Role) newZone(name string) *ZoneContext {
	return &ZoneContext{
		Zone: &types.Zone{
			Name:       strings.ToLower(name),
			DefaultTTL: DefaultTTL,
		},
		inst:        r.i,
		h:           make([]Handler, 0),
		records:     make(map[string]map[string]*Record),
		recordsSync: sync.RWMutex{},
	}
}

func (r *Role) zoneFromKV(raw *mvccpb.KeyValue) (*ZoneContext, error) {
	prefix := r.i.KV().Key(types.KeyRole, types.KeyZones).Prefix(true).String()
	name := strings.TrimPrefix(string(raw.Key), prefix)
	z := r.newZone(name)
	z.log = r.log.With(zap.String("zone", name))
	z.etcdKey = string(raw.Key)

	// Try loading protobuf first
	err := proto.Unmarshal(raw.Value, z.Zone)
	if err == nil {
		return z, nil
	}
	// Otherwise try json
	err = json.Unmarshal(raw.Value, &z)
	if err != nil {
		return nil, err
	}
	return z, nil
}

func (z *ZoneContext) Init(ctx context.Context) {
	for _, handlerCfg := range z.HandlerConfigs {
		t := handlerCfg.Fields["type"].GetStringValue()
		var handler Handler
		switch t {
		case CoreDNSType:
			handler = NewCoreDNS(z, handlerCfg.Fields)
		case BlockyForwarderType:
			handler = NewBlockyForwarder(z, handlerCfg.Fields)
		case IPForwarderType:
			handler = NewIPForwarderHandler(z, handlerCfg.Fields)
		case EtcdType:
			handler = NewEtcdHandler(z, handlerCfg.Fields)
		case MemoryType:
			handler = NewMemoryHandler(z, handlerCfg.Fields)
		default:
			z.log.Warn("invalid forwarder type", zap.String("type", t))
		}
		z.h = append(z.h, handler)
	}

	// start watching all records in this zone, in case etcd goes down
	go z.watchZoneRecords(ctx)
}

func (z *ZoneContext) watchZoneRecords(ctx context.Context) {
	evtHandler := func(ev *clientv3.Event) {
		z.recordsSync.Lock()
		defer z.recordsSync.Unlock()
		rec, err := z.recordFromKV(ev.Kv)
		if _, ok := z.records[rec.recordKey]; !ok {
			z.records[rec.recordKey] = make(map[string]*Record)
		}
		if ev.Type == clientv3.EventTypeDelete {
			delete(z.records[rec.recordKey], rec.uid)
			dnsRecordsMetric.WithLabelValues(z.Name).Dec()
		} else {
			// Check if the record parsed above actually was parsed correctly,
			// we don't care for that when removing, but prevent adding
			// empty records
			if err != nil {
				return
			}
			if _, ok := z.records[rec.recordKey][rec.uid]; !ok {
				dnsRecordsMetric.WithLabelValues(z.Name).Inc()
			}
			z.records[rec.recordKey][rec.uid] = rec
		}
	}
	ctx, canc := context.WithCancel(ctx)
	z.recordsWatchCtx = canc

	prefix := z.inst.KV().Key(z.etcdKey).Prefix(true).String()

	records, err := z.inst.KV().Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		z.log.Warn("failed to list initial records", zap.Error(err))
		time.Sleep(5 * time.Second)
		z.watchZoneRecords(ctx)
		return
	}
	for _, record := range records.Kvs {
		evtHandler(&clientv3.Event{
			Type: mvccpb.PUT,
			Kv:   record,
		})
	}

	watchChan := z.inst.KV().Watch(
		ctx,
		prefix,
		clientv3.WithPrefix(),
	)
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			go evtHandler(event)
		}
	}
}

func (z *ZoneContext) StopWatchingRecords() {
	if z.recordsWatchCtx != nil {
		z.recordsWatchCtx()
	}
}

func (z *ZoneContext) put(ctx context.Context) error {
	raw, err := proto.Marshal(z.Zone)
	if err != nil {
		return err
	}
	_, err = z.inst.KV().Put(
		ctx,
		z.inst.KV().Key(
			types.KeyRole,
			types.KeyZones,
			strings.ToLower(z.Name),
		).String(),
		string(raw),
	)
	if err != nil {
		return err
	}
	return nil
}
