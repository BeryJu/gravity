package dns

import (
	"context"
	"fmt"
	"runtime"

	"beryju.io/gravity/pkg/extconfig"
	instanceTypes "beryju.io/gravity/pkg/instance/types"
	"beryju.io/gravity/pkg/roles"
	apiTypes "beryju.io/gravity/pkg/roles/api/types"
	dhcpTypes "beryju.io/gravity/pkg/roles/dhcp/types"
	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/storage/watcher"
	"github.com/getsentry/sentry-go"
	"github.com/miekg/dns"
	"github.com/swaggest/rest/web"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.uber.org/zap"
)

type Role struct {
	i     roles.Instance
	ctx   context.Context
	zones *watcher.Watcher[*Zone]

	cfg *RoleConfig
	log *zap.Logger

	servers []*dns.Server
	m       *dns.ServeMux
}

func init() {
	roles.Register("dns", func(i roles.Instance) roles.Role {
		return New(i)
	})
}

func New(instance roles.Instance) *Role {
	r := &Role{
		servers: make([]*dns.Server, 0),
		log:     instance.Log(),
		i:       instance,
		ctx:     instance.Context(),
		m:       dns.NewServeMux(),
	}
	r.zones = watcher.New(
		func(kv *mvccpb.KeyValue) (*Zone, error) {
			z, err := r.zoneFromKV(kv)
			if err != nil {
				return nil, err
			}
			z.Init(r.i.Context())
			return z, nil
		},
		r.i.KV(),
		r.i.KV().Key(types.KeyRole, types.KeyZones).Prefix(true),
		watcher.WithBeforeUpdate(func(zone *Zone, direction mvccpb.Event_EventType) {
			zone.StopWatchingRecords()
		}),
	)
	r.i.AddEventListener(dhcpTypes.EventTopicDHCPLeasePut, r.eventHandlerDHCPLeasePut)
	r.i.AddEventListener(types.EventTopicDNSRecordCreateForward, r.eventHandlerCreateForward)
	r.i.AddEventListener(types.EventTopicDNSRecordCreateReverse, r.eventHandlerCreateReverse)
	r.i.AddEventListener(instanceTypes.EventTopicInstanceFirstStart, func(ev *roles.Event) {
		// On first start create a zone that'll forward to a reasonable default
		zone := r.newZone(types.DNSRootZone)
		zone.HandlerConfigs = []map[string]interface{}{
			{
				"type": "memory",
			},
			{
				"type": "etcd",
			},
			{
				"type":      "forward_ip",
				"to":        []string{extconfig.Get().FallbackDNS},
				"cache_ttl": 3600,
			},
		}
		err := zone.put(ev.Context)
		if err != nil {
			r.log.Warn("failed to write startup zone config", zap.Error(err))
		}
	})
	r.i.AddEventListener(apiTypes.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		svc.Get("/api/v1/dns/zones", r.APIZonesGet())
		svc.Post("/api/v1/dns/zones", r.APIZonesPut())
		svc.Delete("/api/v1/dns/zones", r.APIZonesDelete())
		svc.Post("/api/v1/dns/zones/import", r.APIZonesImport())
		svc.Get("/api/v1/dns/zones/records", r.APIRecordsGet())
		svc.Post("/api/v1/dns/zones/records", r.APIRecordsPut())
		svc.Delete("/api/v1/dns/zones/records", r.APIRecordsDelete())
		svc.Get("/api/v1/roles/dns", r.APIRoleConfigGet())
		svc.Post("/api/v1/roles/dns", r.APIRoleConfigPut())
	})
	return r
}

func (r *Role) Start(ctx context.Context, config []byte) error {
	r.cfg = r.decodeRoleConfig(config)

	start := sentry.TransactionFromContext(ctx).StartChild("gravity.dns.start")
	defer start.Finish()

	r.zones.Start(start.Context())

	r.m.HandleFunc(
		types.DNSRootZone,
		r.recoverMiddleware(
			r.loggingMiddleware(
				r.handler,
			),
		),
	)
	if r.cfg.Port < 1 {
		return nil
	}
	listen := extconfig.Get().Listen(r.cfg.Port)
	if runtime.GOOS == "darwin" {
		listen = fmt.Sprintf(":%d", r.cfg.Port)
	}

	srv := func(idx int) {
		server := r.servers[idx]
		r.log.Info("starting DNS Server", zap.String("listen", listen), zap.String("proto", server.Net))
		err := server.ListenAndServe()
		if err != nil {
			r.log.Warn("failed to start dns server", zap.String("listen", listen), zap.String("proto", server.Net), zap.Error(err))
		}
	}

	r.servers = []*dns.Server{
		{
			Addr:    listen,
			Net:     "udp",
			Handler: r.m,
		},
		{
			Addr:    listen,
			Net:     "tcp",
			Handler: r.m,
		},
	}
	go srv(0)
	go srv(1)
	return nil
}

func (r *Role) Handler(w dns.ResponseWriter, req *dns.Msg) {
	r.m.ServeDNS(w, req)
}

func (r *Role) Stop() {
	r.zones.Stop()
	for _, server := range r.servers {
		err := server.Shutdown()
		if err != nil && err.Error() != "dns: server not started" {
			r.log.Warn("failed to stop server", zap.Error(err))
		}
	}
}
