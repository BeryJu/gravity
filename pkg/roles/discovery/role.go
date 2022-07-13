package discovery

import (
	"context"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	apitypes "beryju.io/gravity/pkg/roles/api/types"
	"beryju.io/gravity/pkg/roles/discovery/types"

	log "github.com/sirupsen/logrus"
	"github.com/swaggest/rest/web"
)

type DiscoveryRole struct {
	log *log.Entry
	i   roles.Instance
	cfg *DiscoveryRoleConfig
	ctx context.Context
}

func New(instance roles.Instance) *DiscoveryRole {
	r := &DiscoveryRole{
		log: instance.GetLogger().WithField("role", types.KeyRole),
		i:   instance,
	}
	r.i.AddEventListener(apitypes.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		svc.Get("/api/v1/discovery/subnets", r.apiHandlerSubnets())
		svc.Get("/api/v1/discovery/devices", r.apiHandlerDevices())
		svc.Post("/api/v1/discovery/devices/apply", r.apiHandlerDeviceApply())
	})
	return r
}

func (r *DiscoveryRole) Start(ctx context.Context, config []byte) error {
	r.ctx = ctx
	r.cfg = r.decodeDiscoveryRoleConfig(config)
	if !r.cfg.Enabled || extconfig.Get().ListenOnlyMode {
		r.log.Info("Not enabling discovery")
		return nil
	}
	r.startWatchSubnets()
	return nil
}

func (r *DiscoveryRole) Stop() {
}
