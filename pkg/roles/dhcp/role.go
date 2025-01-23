package dhcp

import (
	"context"
	"fmt"
	"net"
	"strings"

	"beryju.io/gravity/pkg/extconfig"
	instanceTypes "beryju.io/gravity/pkg/instance/types"
	"beryju.io/gravity/pkg/roles"
	apitypes "beryju.io/gravity/pkg/roles/api/types"
	"beryju.io/gravity/pkg/roles/dhcp/options"
	optTypes "beryju.io/gravity/pkg/roles/dhcp/options/types"
	"beryju.io/gravity/pkg/roles/dhcp/oui"
	"beryju.io/gravity/pkg/roles/dhcp/types"
	"beryju.io/gravity/pkg/storage/watcher"
	"github.com/getsentry/sentry-go"
	"github.com/swaggest/rest/web"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.uber.org/zap"
	"golang.org/x/net/ipv4"

	"github.com/insomniacslk/dhcp/dhcpv4"
	"github.com/insomniacslk/dhcp/dhcpv4/server4"
)

type Role struct {
	i   roles.Instance
	ctx context.Context

	scopes            *watcher.Watcher[*Scope]
	leases            *watcher.Watcher[*Lease]
	optionDefinitions *watcher.Watcher[*optTypes.OptionDefinition]

	cfg *RoleConfig

	s4  *handler4
	log *zap.Logger

	oui *oui.OuiDb
}

func init() {
	roles.Register("dhcp", func(i roles.Instance) roles.Role {
		return New(i)
	})
}

func New(instance roles.Instance) *Role {
	r := &Role{
		log: instance.Log(),
		i:   instance,
		ctx: instance.Context(),
	}
	r.scopes = watcher.New(
		func(kv *mvccpb.KeyValue) (*Scope, error) {
			s, err := r.scopeFromKV(kv)
			if err != nil {
				return nil, err
			}
			s.calculateUsage()
			return s, nil
		},
		r.i.KV(),
		r.i.KV().Key(
			types.KeyRole,
			types.KeyScopes,
		).Prefix(true),
	)

	r.leases = watcher.New(
		func(kv *mvccpb.KeyValue) (*Lease, error) {
			s, err := r.leaseFromKV(kv)
			if err != nil {
				return nil, err
			}
			return s, nil
		},
		r.i.KV(),
		r.i.KV().Key(
			types.KeyRole,
			types.KeyLeases,
		).Prefix(true), watcher.WithAfterInitialLoad[*Lease](func() {
			// Re-calculate scope usage after all leases are loaded
			for _, s := range r.scopes.Iter() {
				s.calculateUsage()
			}
		}),
	)

	r.optionDefinitions = watcher.NewProto[*optTypes.OptionDefinition](
		r.i.KV(),
		r.i.KV().Key(
			types.KeyRole,
			types.KeyOptionDefinitions,
		).Prefix(true),
	)

	r.s4 = &handler4{
		role: r,
	}
	r.i.AddEventListener(instanceTypes.EventTopicInstanceBootstrapped, options.Bootstrap(r.i))
	r.i.AddEventListener(types.EventTopicDHCPCreateLease, r.eventCreateLease)
	r.i.AddEventListener(apitypes.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		svc.Get("/api/v1/dhcp/scopes", r.APIScopesGet())
		svc.Post("/api/v1/dhcp/scopes", r.APIScopesPut())
		svc.Delete("/api/v1/dhcp/scopes", r.APIScopesDelete())
		svc.Post("/api/v1/dhcp/scopes/import", r.APIScopesImport())
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

// Deprecated: FunctionName is deprecated.
func (r *Role) Handle4(re *Request4) *dhcpv4.DHCPv4 {
	return r.s4.HandleRequest(re)
}

func (r *Role) Handler4() *handler4 {
	return r.s4
}

func (r *Role) Start(ctx context.Context, config []byte) error {
	r.cfg = r.decodeRoleConfig(config)

	start := sentry.TransactionFromContext(ctx).StartChild("gravity.dhcp.start")
	defer start.Finish()

	r.scopes.Start(r.ctx)
	r.leases.Start(r.ctx)

	if r.cfg.Port < 1 {
		return nil
	}

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
	ifName := extconfig.Get().Instance.Interface
	udpConn, err := server4.NewIPv4UDPConn(ifName, laddr)
	if err != nil {
		return err
	}
	r.s4.pc = ipv4.NewPacketConn(udpConn)
	var ifi *net.Interface
	if ifName != "" {
		ifi, err = net.InterfaceByName(ifName)
		if err != nil {
			return fmt.Errorf("DHCPv4: Listen could not find interface %s: %v", ifName, err)
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
	r.log.Info("starting DHCP Server", zap.Int("port", r.cfg.Port), zap.String("interface", extconfig.Get().Instance.Interface))
	err := r.s4.Serve()
	if !isErrNetClosing(err) {
		return err
	}
	return nil
}

func (r *Role) Stop() {
	r.scopes.Stop()
	r.leases.Stop()
	if r.s4 != nil && r.s4.pc != nil {
		err := r.s4.pc.Close()
		if err != nil {
			r.log.Warn("Failed to stop packet conn", zap.Error(err))
		}
	}
}

func (r *Role) DeviceIdentifier(m *dhcpv4.DHCPv4) string {
	return m.ClientHWAddr.String()
}
