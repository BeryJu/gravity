package discovery

import (
	"context"
	"fmt"
	"net/netip"

	"beryju.io/gravity/pkg/extconfig"
	instanceTypes "beryju.io/gravity/pkg/instance/types"
	"beryju.io/gravity/pkg/roles"
	apiTypes "beryju.io/gravity/pkg/roles/api/types"
	"beryju.io/gravity/pkg/roles/discovery/types"
	"beryju.io/gravity/pkg/storage/watcher"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"go.uber.org/zap"

	"github.com/swaggest/rest/web"
)

type Role struct {
	log     *zap.Logger
	i       roles.Instance
	cfg     *RoleConfig
	ctx     context.Context
	watcher *watcher.Watcher[*Subnet]
}

func init() {
	roles.Register("discovery", func(i roles.Instance) roles.Role {
		return New(i)
	})
}

func New(instance roles.Instance) *Role {
	r := &Role{
		log: instance.Log(),
		i:   instance,
		ctx: instance.Context(),
	}
	r.watcher = watcher.New(
		func(kv *mvccpb.KeyValue) (*Subnet, error) {
			sub, err := r.subnetFromKV(kv)
			if err != nil {
				r.log.Warn("failed to parse subnet", zap.Error(err))
				return nil, err
			}
			go sub.RunDiscovery(context.Background())
			return sub, err
		},
		r.i.KV(),
		r.i.KV().Key(types.KeyRole, types.KeySubnets).Prefix(true),
	)
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
		subnet := r.NewSubnet(fmt.Sprintf("instance-subnet-%s", extconfig.Get().Instance.Identifier))
		ip := netip.MustParseAddr(extconfig.Get().Instance.IP)
		// TODO: Check interface to see if it's actually a /24
		prefix, err := ip.Prefix(24)
		if err != nil {
			r.log.Warn("failed to get prefix", zap.Error(err))
			return
		}
		subnet.CIDR = prefix.String()
		subnet.DNSResolver = extconfig.Get().FallbackDNS
		subnet.DiscoveryTTL = 86400
		err = subnet.put(ev.Context)
		if err != nil {
			r.log.Warn("failed to put subnet", zap.Error(err))
			return
		}
		go subnet.RunDiscovery(context.Background())
	})
	return r
}

func (r *Role) Start(ctx context.Context, config []byte) error {
	r.cfg = r.decodeRoleConfig(config)
	if !r.cfg.Enabled || extconfig.Get().ListenOnlyMode {
		return roles.ErrRoleNotConfigured
	}
	r.watcher.Start(r.ctx)
	return nil
}

func (r *Role) Stop() {
	r.watcher.Stop()
}
