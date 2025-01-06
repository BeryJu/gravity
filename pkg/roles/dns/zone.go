package dns

import (
	"context"
	"encoding/json"
	"net"
	"net/netip"
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
)

const (
	DefaultTTL = 3600
)

type Zone struct {
	inst roles.Instance
	role *Role

	records         map[string]map[string]*Record
	recordsWatchCtx context.CancelFunc

	log  *zap.Logger
	Name string `json:"-"`

	etcdKey        string
	HandlerConfigs []map[string]interface{} `json:"handlerConfigs"`

	h []Handler

	recordsSync sync.RWMutex
	DefaultTTL  uint32 `json:"defaultTTL"`

	Authoritative bool   `json:"authoritative"`
	Hook          string `json:"hook"`
}

func (z *Zone) Handlers() []Handler {
	return append([]Handler{
		// Internal SOA handler
		&SOAIntHandler{z},
	}, z.h...)
}

func (z *Zone) resolveUpdateMetrics(dur time.Duration, q *utils.DNSRequest, h Handler) {
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
		).Observe(dur.Seconds())
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

func getIP(addr net.Addr) *netip.Addr {
	clientIP := ""
	switch addr := addr.(type) {
	case *net.UDPAddr:
		clientIP = addr.IP.String()
	case *net.TCPAddr:
		clientIP = addr.IP.String()
	}
	i, err := netip.ParseAddr(clientIP)
	if err != nil {
		return nil
	}
	return &i
}

func (z *Zone) resolve(w dns.ResponseWriter, r *utils.DNSRequest, span *sentry.Span) {
	z.inst.ExecuteHook(roles.HookOptions{
		Source: z.Hook,
		Method: "onDNSRequestBefore",
	}, r)
	for idx, handler := range z.Handlers() {
		ss := span.StartChild("gravity.dns.request.handler")
		ss.Description = handler.Identifier()
		z.log.Debug("sending request to handler", zap.String("handler", handler.Identifier()))
		ss.SetTag("gravity.dns.handler", handler.Identifier())
		// Create new request for handler with the correct context
		hr := r.Chain(r.Msg, ss.Context(), utils.DNSRoutingMeta{
			HandlerIdx:      idx,
			HasMoreHandlers: len(z.h)-(idx+1) > 0,
			ResolveRequest: func(w dns.ResponseWriter, r *utils.DNSRequest) {
				z.log.Debug("Next lookup iteration", zap.Int("iter", r.Iteration()+1))
				z.role.rootHandler(w, r)
			},
		})

		handlerReply := handler.Handle(utils.NewFakeDNSWriter(w), hr)
		ss.Finish()

		if handlerReply != nil {
			z.log.Debug("returning reply from handler", zap.String("handler", handler.Identifier()))
			if i := getIP(w.RemoteAddr()); i != nil && (i.IsPrivate() || i.IsLoopback()) {
				handlerReply.RecursionAvailable = r.Msg.RecursionDesired
			}
			handlerReply.SetEdns0(4000, false)
			handlerReply.SetReply(r.Msg)
			z.inst.ExecuteHook(roles.HookOptions{
				Source: z.Hook,
				Method: "onDNSRequestAfter",
			}, r, handlerReply)
			err := w.WriteMsg(handlerReply)
			if err != nil {
				z.log.Warn("failed to write response", zap.Error(err))
			}
			z.resolveUpdateMetrics(ss.EndTime.Sub(ss.StartTime), r, handler)
			return
		}
		z.log.Debug("no reply, trying next handler", zap.String("handler", handler.Identifier()))
	}
	z.log.Debug("no handler has a reply")
	fallback := new(dns.Msg)
	fallback.SetReply(r.Msg)
	err := w.WriteMsg(fallback)
	if err != nil {
		z.log.Warn("failed to write response", zap.Error(err))
	}
}

func (r *Role) newZone(name string) *Zone {
	return &Zone{
		Name:        strings.ToLower(name),
		DefaultTTL:  DefaultTTL,
		inst:        r.i,
		h:           make([]Handler, 0),
		records:     make(map[string]map[string]*Record),
		recordsSync: sync.RWMutex{},
		role:        r,
	}
}

func (r *Role) zoneFromKV(raw *mvccpb.KeyValue) (*Zone, error) {
	prefix := r.i.KV().Key(types.KeyRole, types.KeyZones).Prefix(true).String()
	name := strings.TrimPrefix(string(raw.Key), prefix)
	z := r.newZone(name)
	z.etcdKey = string(raw.Key)
	err := json.Unmarshal(raw.Value, &z)
	if err != nil {
		return nil, err
	}
	z.log = r.log.With(zap.String("zone", z.Name))
	return z, nil
}

func (z *Zone) Init(ctx context.Context) {
	for _, handlerCfg := range z.HandlerConfigs {
		t := handlerCfg["type"].(string)
		hc, ok := HandlerRegistry.Find(t)
		if !ok {
			z.log.Warn("invalid forwarder type", zap.String("type", t))
			continue
		}
		z.h = append(z.h, hc(z, handlerCfg))
	}

	// start watching all records in this zone, in case etcd goes down
	z.watchZoneRecords(ctx)
}

func (z *Zone) watchZoneRecords(ctx context.Context) {
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
	go func() {
		for watchResp := range watchChan {
			for _, event := range watchResp.Events {
				go evtHandler(event)
			}
		}
	}()
}

func (z *Zone) StopWatchingRecords() {
	if z != nil && z.recordsWatchCtx != nil {
		z.recordsWatchCtx()
	}
}

func (z *Zone) put(ctx context.Context) error {
	raw, err := json.Marshal(&z)
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
