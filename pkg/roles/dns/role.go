package dns

import (
	"context"
	"fmt"
	"runtime"
	"sync"

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
		svc.Post("/api/v1/dns/zones/{zone}", r.apiHandlerZonesPut())
		svc.Delete("/api/v1/dns/zones/{zone}", r.apiHandlerZonesDelete())
		svc.Get("/api/v1/dns/zones/{zone}/records", r.apiHandlerZoneRecordsGet())
		svc.Post("/api/v1/dns/zones/{zone}/records/{hostname}", r.apiHandlerZoneRecordsPut())
		svc.Delete("/api/v1/dns/zones/{zone}/records/{hostname}", r.apiHandlerZoneRecordsDelete())
		svc.Get("/api/v1/roles/dns", r.apiHandlerRoleConfigGet())
		svc.Post("/api/v1/roles/dns", r.apiHandlerRoleConfigPut())
	})
	r.loadInitialZones()
	return r
}

func (r *Role) Start(ctx context.Context, config []byte) error {
	r.ctx = ctx
	r.cfg = r.decodeRoleConfig(config)

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
	wg := sync.WaitGroup{}
	wg.Add(2)

	listen := extconfig.Get().Listen(r.cfg.Port)
	if runtime.GOOS == "darwin" {
		listen = fmt.Sprintf(":%d", r.cfg.Port)
	}

	srv := func(proto string) {
		defer wg.Done()
		server := &dns.Server{
			Addr:    listen,
			Net:     proto,
			Handler: dnsMux,
		}
		r.servers = append(r.servers, server)
		r.log.WithField("listen", listen).WithField("proto", proto).Info("starting DNS Server")
		err := server.ListenAndServe()
		if err != nil {
			r.log.WithField("listen", listen).WithField("proto", proto).WithError(err).Warning("failed to start dns server")
		}
	}

	go srv("udp")
	go srv("tcp")
	wg.Wait()
	return nil
}

func (r *Role) Stop() {
	for _, server := range r.servers {
		server.Shutdown()
	}
}
