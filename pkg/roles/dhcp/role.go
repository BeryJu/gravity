package dhcp

import (
	"context"
	"fmt"
	"net"

	"beryju.io/gravity/pkg/roles"
	apitypes "beryju.io/gravity/pkg/roles/api/types"
	"beryju.io/gravity/pkg/roles/dhcp/types"
	log "github.com/sirupsen/logrus"
	"github.com/swaggest/rest/web"
	"golang.org/x/net/ipv4"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
)

type Role struct {
	scopes map[string]*Scope
	leases map[string]*Lease

	cfg *RoleConfig

	s4  *handler4
	log *log.Entry
	i   roles.Instance
	ctx context.Context
}

func New(instance roles.Instance) *Role {
	r := &Role{
		log:    instance.Log(),
		i:      instance,
		scopes: make(map[string]*Scope),
		leases: make(map[string]*Lease),
	}
	r.i.AddEventListener(types.EventTopicDHCPCreateLease, r.eventCreateLease)
	r.i.AddEventListener(apitypes.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		svc.Get("/api/v1/dhcp/scopes", r.apiHandlerScopesGet())
		svc.Post("/api/v1/dhcp/scopes", r.apiHandlerScopesPut())
		svc.Delete("/api/v1/dhcp/scopes", r.apiHandlerScopesDelete())
		svc.Get("/api/v1/dhcp/scopes/leases", r.apiHandlerLeasesGet())
		svc.Post("/api/v1/dhcp/scopes/leases", r.apiHandlerLeasesPut())
		svc.Post("/api/v1/dhcp/scopes/leases/wol", r.apiHandlerLeasesWOL())
		svc.Delete("/api/v1/dhcp/scopes/leases", r.apiHandlerLeasesDelete())
		svc.Get("/api/v1/roles/dhcp", r.apiHandlerRoleConfigGet())
		svc.Post("/api/v1/roles/dhcp", r.apiHandlerRoleConfigPut())
	})
	return r
}

func (r *Role) Start(ctx context.Context, config []byte) error {
	r.ctx = ctx
	r.cfg = r.decodeRoleConfig(config)

	r.loadInitialScopes()
	r.loadInitialLeases()

	go r.startWatchScopes()
	go r.startWatchLeases()

	go func() {
		err := r.startServer4()
		if err != nil {
			r.log.WithError(err).Warning("failed to listen")
		}
	}()
	return nil
}

func (r *Role) startServer4() error {
	laddr := &net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: r.cfg.Port,
	}
	var err error
	l4 := &handler4{
		role: r,
	}
	udpConn, err := server4.NewIPv4UDPConn(laddr.Zone, laddr)
	if err != nil {
		return err
	}
	l4.pc = ipv4.NewPacketConn(udpConn)
	var ifi *net.Interface
	if laddr.Zone != "" {
		ifi, err = net.InterfaceByName(laddr.Zone)
		if err != nil {
			return fmt.Errorf("DHCPv4: Listen could not find interface %s: %v", laddr.Zone, err)
		}
		l4.iface = *ifi
	} else {
		// When not bound to an interface, we need the information in each
		// packet to know which interface it came on
		err = l4.pc.SetControlMessage(ipv4.FlagInterface, true)
		if err != nil {
			return err
		}
	}

	if laddr.IP.IsMulticast() {
		err = l4.pc.JoinGroup(ifi, laddr)
		if err != nil {
			return err
		}
	}
	r.s4 = l4
	r.log.WithField("port", r.cfg.Port).Info("starting DHCP Server")
	return l4.Serve()
}

func (r *Role) Stop() {
	if r.s4 != nil {
		r.s4.pc.Close()
	}
}

func (r *Role) DeviceIdentifier(m *dhcpv4.DHCPv4) string {
	return m.ClientHWAddr.String()
}
