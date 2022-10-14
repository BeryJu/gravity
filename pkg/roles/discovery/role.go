package discovery

import (
	"context"
	"net/netip"

	"beryju.io/gravity/pkg/extconfig"
	instanceTypes "beryju.io/gravity/pkg/instance/types"
	"beryju.io/gravity/pkg/roles"
	apiTypes "beryju.io/gravity/pkg/roles/api/types"

	log "github.com/sirupsen/logrus"
	"github.com/swaggest/rest/web"
)

type Role struct {
	log *log.Entry
	i   roles.Instance
	cfg *RoleConfig
	ctx context.Context
}

func New(instance roles.Instance) *Role {
	r := &Role{
		log: instance.Log(),
		i:   instance,
	}
	r.i.AddEventListener(apiTypes.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		svc.Get("/api/v1/discovery/subnets", r.APISubnetsGet())
		svc.Post("/api/v1/discovery/subnets", r.APISubnetsPut())
		svc.Post("/api/v1/discovery/subnets/start", r.APISubnetsStart())
		svc.Delete("/api/v1/discovery/subnets", r.APISubnetsDelete())
		svc.Get("/api/v1/discovery/devices", r.APIDevicesGet())
		svc.Post("/api/v1/discovery/devices/apply", r.APIDevicesApply())
		svc.Delete("/api/v1/discovery/devices/delete", r.APIDevicesDelete())
		svc.Get("/api/v1/roles/discovery", r.APIRoleConfigGet())
		svc.Post("/api/v1/roles/discovery", r.APIRoleConfigPut())
	})
	r.i.AddEventListener(instanceTypes.EventTopicInstanceFirstStart, func(ev *roles.Event) {
		// On first start create a subnet based on the instance IP
		subnet := r.NewSubnet("default-instance-subnet")
		ip := netip.MustParseAddr(extconfig.Get().Instance.IP)
		prefix, err := ip.Prefix(24)
		if err != nil {
			r.log.WithError(err).Warning("failed to get prefix")
			return
		}
		subnet.CIDR = prefix.String()
		subnet.DNSResolver = extconfig.Get().FallbackDNS
		subnet.DiscoveryTTL = 86400
		subnet.put(ev.Context)
	})
	return r
}

func (r *Role) Start(ctx context.Context, config []byte) error {
	r.ctx = ctx
	r.cfg = r.decodeRoleConfig(config)
	if !r.cfg.Enabled || extconfig.Get().ListenOnlyMode {
		return roles.ErrRoleNotConfigured
	}
	r.startWatchSubnets()
	return nil
}

func (r *Role) Stop() {
}
