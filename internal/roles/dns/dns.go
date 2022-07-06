package dns

import (
	"fmt"
	"sync"
	"time"

	"github.com/beryju/dns-dhcp-etcd-thingy/internal/roles"
	log "github.com/sirupsen/logrus"

	"github.com/miekg/dns"
)

const (
	KeyRole  = "dns"
	KeyZones = "zones"
)

type DNSServerRole struct {
	servers map[string]*dns.Server
	zones   map[string]*Zone

	log *log.Entry
	i   roles.Instance
}

func New(instance roles.Instance) *DNSServerRole {
	return &DNSServerRole{
		servers: make(map[string]*dns.Server),
		zones:   make(map[string]*Zone, 0),
		log:     log.WithField("role", "dns"),
		i:       instance,
	}
}

func (r *DNSServerRole) Start(config []byte) error {
	cfg := r.decodeDNSRoleConfig(config)

	go r.startWatchZones()

	dnsMux := dns.NewServeMux()
	dnsMux.HandleFunc(".", r.loggingHandler(r.handler))
	wg := sync.WaitGroup{}
	wg.Add(2)
	srv := func(proto string, wg sync.WaitGroup) {
		defer wg.Done()
		server := &dns.Server{
			Addr:    fmt.Sprintf(":%d", cfg.Port),
			Net:     proto,
			Handler: dnsMux,
		}
		r.servers[proto] = server
		r.log.WithField("port", cfg.Port).WithField("proto", proto).Info("Starting DNS Server")
		err := server.ListenAndServe()
		if err != nil {
			r.log.WithField("port", cfg.Port).WithField("proto", proto).WithError(err).Warning("failed to start dns server")
		}
	}
	go srv("udp", wg)
	time.Sleep(1)
	go srv("tcp", wg)
	wg.Wait()
	return nil
}

func (r *DNSServerRole) Stop() {
	for _, server := range r.servers {
		server.Shutdown()
	}
}

func (r *DNSServerRole) HandleEvent(ev *roles.Event[any]) {
	r.log.Debug(ev)
}
