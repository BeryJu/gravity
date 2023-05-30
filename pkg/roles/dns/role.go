package dns

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"beryju.io/gravity/pkg/extconfig"
	instanceTypes "beryju.io/gravity/pkg/instance/types"
	"beryju.io/gravity/pkg/roles"
	apiTypes "beryju.io/gravity/pkg/roles/api/types"
	dhcpTypes "beryju.io/gravity/pkg/roles/dhcp/types"
	"beryju.io/gravity/pkg/roles/dns/types"
	"github.com/getsentry/sentry-go"
	"github.com/miekg/dns"
	"github.com/swaggest/rest/web"
	"go.uber.org/zap"
	"google.golang.org/protobuf/types/known/structpb"
)

type Role struct {
	i     roles.Instance
	ctx   context.Context
	zones map[string]*ZoneContext

	cfg     *RoleConfig
	log     *zap.Logger
	servers []*dns.Server
	zonesM  sync.RWMutex
}

func New(instance roles.Instance) *Role {
	r := &Role{
		servers: make([]*dns.Server, 0),
		zones:   make(map[string]*ZoneContext, 0),
		zonesM:  sync.RWMutex{},
		log:     instance.Log(),
		i:       instance,
		ctx:     instance.Context(),
	}
	r.i.AddEventListener(dhcpTypes.EventTopicDHCPLeasePut, r.eventHandlerDHCPLeaseGiven)
	r.i.AddEventListener(types.EventTopicDNSRecordCreateForward, r.eventCreateForward)
	r.i.AddEventListener(types.EventTopicDNSRecordCreateReverse, r.eventCreateReverse)
	r.i.AddEventListener(instanceTypes.EventTopicInstanceFirstStart, func(ev *roles.Event) {
		// On first start create a zone that'll forward to a reasonable default
		zone := r.newZone(".")
		zone.HandlerConfigs = []*structpb.Struct{
			{
				Fields: map[string]*structpb.Value{
					"type": structpb.NewStringValue("memory"),
				},
			},
			{
				Fields: map[string]*structpb.Value{
					"type": structpb.NewStringValue("etcd"),
				},
			},
			{
				Fields: map[string]*structpb.Value{
					"type":      structpb.NewStringValue("forward_ip"),
					"to":        structpb.NewStringValue(extconfig.Get().FallbackDNS),
					"cache_ttl": structpb.NewStringValue("3600"),
				},
			},
		}
		zone.put(ev.Context)
	})
	r.i.AddEventListener(apiTypes.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		svc.Get("/api/v1/dns/metrics", r.APIMetrics())
		svc.Get("/api/v1/dns/zones", r.APIZonesGet())
		svc.Post("/api/v1/dns/zones", r.APIZonesPut())
		svc.Delete("/api/v1/dns/zones", r.APIZonesDelete())
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

	r.loadInitialZones(start.Context())
	go r.startWatchZones(start.Context())

	dnsMux := dns.NewServeMux()
	dnsMux.HandleFunc(
		".",
		r.recoverMiddleware(
			r.loggingMiddleware(
				r.Handler,
			),
		),
	)
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
			Handler: dnsMux,
		},
		{
			Addr:    listen,
			Net:     "tcp",
			Handler: dnsMux,
		},
	}
	go srv(0)
	go srv(1)
	return nil
}

func (r *Role) Stop() {
	for _, server := range r.servers {
		server.Shutdown()
	}
}
