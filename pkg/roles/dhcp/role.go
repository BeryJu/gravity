package dhcp

import (
	"context"
	"fmt"
	"net"
	"strings"
	"sync"

	"beryju.io/gravity/pkg/roles"
	apitypes "beryju.io/gravity/pkg/roles/api/types"
	"beryju.io/gravity/pkg/roles/dhcp/oui"
	"beryju.io/gravity/pkg/roles/dhcp/types"
	"github.com/getsentry/sentry-go"
	"github.com/swaggest/rest/web"
	"go.uber.org/zap"
	"golang.org/x/net/ipv4"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
)

type Role struct {
	scopes  map[string]*Scope
	scopesM sync.RWMutex
	leases  map[string]*Lease
	leasesM sync.RWMutex

	cfg *RoleConfig

	s4  *handler4
	log *zap.Logger
	i   roles.Instance
	ctx context.Context

	oui *oui.OuiDb
}

func New(instance roles.Instance) *Role {
	r := &Role{
		log:     instance.Log(),
		i:       instance,
		scopes:  make(map[string]*Scope),
		scopesM: sync.RWMutex{},
		leases:  make(map[string]*Lease),
		leasesM: sync.RWMutex{},
	}
	r.s4 = &handler4{
		role: r,
	}
	r.i.AddEventListener(types.EventTopicDHCPCreateLease, r.eventCreateLease)
	r.i.AddEventListener(apitypes.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		svc.Get("/api/v1/dhcp/scopes", r.APIScopesGet())
		svc.Post("/api/v1/dhcp/scopes", r.APIScopesPut())
		svc.Delete("/api/v1/dhcp/scopes", r.APIScopesDelete())
		svc.Get("/api/v1/dhcp/scopes/leases", r.APILeasesGet())
		svc.Post("/api/v1/dhcp/scopes/leases", r.APILeasesPut())
		svc.Post("/api/v1/dhcp/scopes/leases/wol", r.APILeasesWOL())
		svc.Delete("/api/v1/dhcp/scopes/leases", r.APILeasesDelete())
		svc.Get("/api/v1/roles/dhcp", r.APIRoleConfigGet())
		svc.Post("/api/v1/roles/dhcp", r.APIRoleConfigPut())
	})
	r.initOUI()
	return r
}

func (r *Role) Handler4(re *Request4) *dhcpv4.DHCPv4 {
	return r.s4.HandleRequest(re)
}

func (r *Role) Start(ctx context.Context, config []byte) error {
	r.ctx = ctx
	r.cfg = r.decodeRoleConfig(config)

	start := sentry.StartSpan(ctx, "gravity.dhcp.start")
	defer start.Finish()
	r.loadInitialScopes(start.Context())
	r.loadInitialLeases(start.Context())

	// Since scope usage relies on r.leases, but r.leases is loaded after the scopes,
	// manually update the usage
	for _, s := range r.scopes {
		s.calculateUsage()
	}

	go r.startWatchScopes()
	go r.startWatchLeases()

	err := r.initServer4()
	if err != nil {
		r.log.Warn("failed to setup server", zap.Error(err))
		return err
	}
	go func() {
		err := r.startServer4()
		if err != nil {
			r.log.Warn("failed to listen", zap.Error(err))
		}
	}()
	return nil
}

func (r *Role) initServer4() error {
	laddr := &net.UDPAddr{
		IP:   net.ParseIP("0.0.0.0"),
		Port: r.cfg.Port,
	}
	var err error
	udpConn, err := server4.NewIPv4UDPConn(laddr.Zone, laddr)
	if err != nil {
		return err
	}
	r.s4.pc = ipv4.NewPacketConn(udpConn)
	var ifi *net.Interface
	if laddr.Zone != "" {
		ifi, err = net.InterfaceByName(laddr.Zone)
		if err != nil {
			return fmt.Errorf("DHCPv4: Listen could not find interface %s: %v", laddr.Zone, err)
		}
		r.s4.iface = *ifi
	} else {
		// When not bound to an interface, we need the information in each
		// packet to know which interface it came on
		err = r.s4.pc.SetControlMessage(ipv4.FlagInterface, true)
		if err != nil {
			return err
		}
	}

	if laddr.IP.IsMulticast() {
		err = r.s4.pc.JoinGroup(ifi, laddr)
		if err != nil {
			return err
		}
	}
	return nil
}

var useOfClosedErrMsg = "use of closed network connection"

// isErrNetClosing checks whether is an ErrNetClosing error
func isErrNetClosing(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), useOfClosedErrMsg)
}

func (r *Role) startServer4() error {
	r.log.Info("starting DHCP Server", zap.Int("port", r.cfg.Port))
	err := r.s4.Serve()
	if !isErrNetClosing(err) {
		return err
	}
	return nil
}

func (r *Role) Stop() {
	if r.s4 != nil && r.s4.pc != nil {
		r.s4.pc.Close()
	}
}

func (r *Role) DeviceIdentifier(m *dhcpv4.DHCPv4) string {
	return m.ClientHWAddr.String()
}
