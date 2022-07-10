package dns

import (
	"context"
	"fmt"
	"runtime"
	"sync"

	"beryju.io/ddet/pkg/extconfig"
	"beryju.io/ddet/pkg/roles"
	apitypes "beryju.io/ddet/pkg/roles/api/types"
	dhcptypes "beryju.io/ddet/pkg/roles/dhcp/types"
	"beryju.io/ddet/pkg/roles/dns/types"
	log "github.com/sirupsen/logrus"

	"github.com/miekg/dns"
)

type DNSRole struct {
	servers []*dns.Server
	zones   map[string]*Zone

	log *log.Entry
	i   roles.Instance
	ctx context.Context
}

func New(instance roles.Instance) *DNSRole {
	return &DNSRole{
		servers: make([]*dns.Server, 0),
		zones:   make(map[string]*Zone, 0),
		log:     instance.GetLogger().WithField("role", types.KeyRole),
		i:       instance,
	}
}

func (r *DNSRole) Start(ctx context.Context, config []byte) error {
	r.ctx = ctx
	cfg := r.decodeDNSRoleConfig(config)
	r.i.AddEventListener(dhcptypes.EventTopicDHCPLeaseGiven, r.eventHandlerDHCPLeaseGiven)
	r.i.AddEventListener(apitypes.EventTopicAPIMuxSetup, r.eventHandlerAPIMux)

	go r.startWatchZones()

	dnsMux := dns.NewServeMux()
	dnsMux.HandleFunc(".", r.loggingHandler(r.handler))
	wg := sync.WaitGroup{}
	wg.Add(2)

	listen := fmt.Sprintf("%s:%d", extconfig.Get().Instance.IP, cfg.Port)
	if runtime.GOOS == "darwin" {
		listen = fmt.Sprintf(":%d", cfg.Port)
	}

	srv := func(proto string) {
		defer wg.Done()
		server := &dns.Server{
			Addr:    listen,
			Net:     proto,
			Handler: dnsMux,
		}
		r.servers = append(r.servers, server)
		r.log.WithField("port", cfg.Port).WithField("proto", proto).Info("Starting DNS Server")
		err := server.ListenAndServe()
		if err != nil {
			r.log.WithField("port", cfg.Port).WithField("proto", proto).WithError(err).Warning("failed to start dns server")
		}
	}
	go srv("udp")
	go srv("tcp")
	wg.Wait()
	return nil
}

func (r *DNSRole) Stop() {
	for _, server := range r.servers {
		server.Shutdown()
	}
}
