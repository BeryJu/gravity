package dhcp

import (
	"context"
	"net"
	"sync"

	"beryju.io/gravity/pkg/roles"
	apitypes "beryju.io/gravity/pkg/roles/api/types"
	"beryju.io/gravity/pkg/roles/dhcp/types"
	log "github.com/sirupsen/logrus"
	"github.com/swaggest/rest/web"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
	"github.com/insomniacslk/dhcp/dhcpv6/server6"
)

type Role struct {
	scopes     map[string]*Scope
	leases     map[string]*Lease
	leasesSync sync.RWMutex

	cfg *RoleConfig

	s4  *server4.Server
	s6  *server6.Server
	log *log.Entry
	i   roles.Instance
	ctx context.Context
}

func New(instance roles.Instance) *Role {
	r := &Role{
		log:        instance.Log(),
		i:          instance,
		scopes:     make(map[string]*Scope),
		leases:     make(map[string]*Lease),
		leasesSync: sync.RWMutex{},
	}
	r.i.AddEventListener(types.EventTopicDHCPCreateLease, r.eventCreateLease)
	r.i.AddEventListener(apitypes.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		svc.Get("/api/v1/dhcp/scopes", r.apiHandlerScopes())
		svc.Post("/api/v1/dhcp/scopes/{scope}", r.apiHandlerScopesPut())
		svc.Delete("/api/v1/dhcp/scopes/{scope}", r.apiHandlerScopesDelete())
		svc.Get("/api/v1/dhcp/scopes/{scope}/leases", r.apiHandlerLeases())
		svc.Post("/api/v1/dhcp/scopes/{scope}/leases/{identifier}", r.apiHandlerLeasesPut())
		svc.Delete("/api/v1/dhcp/scopes/{scope}/leases/{identifier}", r.apiHandlerLeasesDelete())
	})
	return r
}

func (r *Role) Start(ctx context.Context, config []byte) error {
	r.ctx = ctx
	r.cfg = r.decodeRoleConfig(config)

	r.loadInitialScopes()
	r.loadInitialScopes()

	go r.startWatchScopes()
	go r.startWatchLeases()

	return r.startServer4()
}

func (r *Role) startServer4() error {
	laddr := net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: r.cfg.Port,
	}
	server, err := server4.NewServer(
		"", // TODO: specify interface to DHCP?
		&laddr,
		r.recoverMiddleware4(
			r.loggingMiddleware4(
				r.handler4,
			),
		),
	)
	if err != nil {
		return err
	}
	r.s4 = server
	r.log.WithField("port", r.cfg.Port).Info("starting DHCP Server")
	return r.s4.Serve()
}

func (r *Role) Stop() {
	if r.s4 != nil {
		r.s4.Close()
	}
}

func (r *Role) DeviceIdentifier(m *dhcpv4.DHCPv4) string {
	return m.ClientHWAddr.String()
}
