package dns

import (
	"context"
	"fmt"
	"runtime"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	apitypes "beryju.io/gravity/pkg/roles/api/types"
	dhcptypes "beryju.io/gravity/pkg/roles/dhcp/types"
	"beryju.io/gravity/pkg/roles/dns/types"
	"github.com/miekg/dns"
	log "github.com/sirupsen/logrus"
	"github.com/swaggest/rest/web"
)

type Role struct {
	servers []*dns.Server
	zones   map[string]*Zone

	cfg *RoleConfig
	log *log.Entry
	i   roles.Instance
	ctx context.Context
}

func New(instance roles.Instance) *Role {
	r := &Role{
		servers: make([]*dns.Server, 0),
		zones:   make(map[string]*Zone, 0),
		log:     instance.Log(),
		i:       instance,
	}
	r.i.AddEventListener(dhcptypes.EventTopicDHCPLeasePut, r.eventHandlerDHCPLeaseGiven)
	r.i.AddEventListener(types.EventTopicDNSRecordCreateForward, r.eventCreateForward)
	r.i.AddEventListener(types.EventTopicDNSRecordCreateReverse, r.eventCreateReverse)
	r.i.AddEventListener(apitypes.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		svc.Get("/api/v1/dns/zones", r.apiHandlerZonesGet())
		svc.Post("/api/v1/dns/zones", r.apiHandlerZonesPut())
		svc.Delete("/api/v1/dns/zones", r.apiHandlerZonesDelete())
		svc.Get("/api/v1/dns/zones/records", r.apiHandlerZoneRecordsGet())
		svc.Post("/api/v1/dns/zones/records", r.apiHandlerZoneRecordsPut())
		svc.Delete("/api/v1/dns/zones/records", r.apiHandlerZoneRecordsDelete())
		svc.Get("/api/v1/roles/dns", r.APIRoleConfigGet())
		svc.Post("/api/v1/roles/dns", r.APIRoleConfigPut())
	})
	return r
}

func (r *Role) Start(ctx context.Context, config []byte) error {
	r.ctx = ctx
	r.cfg = r.decodeRoleConfig(config)

	r.loadInitialZones()
	go r.startWatchZones()

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
		r.log.WithField("listen", listen).WithField("proto", server.Net).Info("starting DNS Server")
		err := server.ListenAndServe()
		if err != nil {
			r.log.WithField("listen", listen).WithField("proto", server.Net).WithError(err).Warning("failed to start dns server")
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
