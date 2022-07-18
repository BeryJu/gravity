package discovery

import (
	"context"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	apitypes "beryju.io/gravity/pkg/roles/api/types"

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
	r.i.AddEventListener(apitypes.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		svc.Get("/api/v1/discovery/subnets", r.apiHandlerSubnets())
		svc.Post("/api/v1/discovery/subnets", r.apiHandlerSubnetsPut())
		svc.Delete("/api/v1/discovery/subnets", r.apiHandlerSubnetsDelete())
		svc.Get("/api/v1/discovery/devices", r.apiHandlerDevices())
		svc.Post("/api/v1/discovery/devices/apply", r.apiHandlerDeviceApply())
		svc.Delete("/api/v1/discovery/devices/delete", r.apiHandlerDevicesDelete())
		svc.Get("/api/v1/roles/discovery", r.apiHandlerRoleConfigGet())
		svc.Post("/api/v1/roles/discovery", r.apiHandlerRoleConfigPut())
	})
	return r
}

func (r *Role) Start(ctx context.Context, config []byte) error {
	r.ctx = ctx
	r.cfg = r.decodeRoleConfig(config)
	if !r.cfg.Enabled || extconfig.Get().ListenOnlyMode {
		r.log.Info("Not enabling discovery")
		return nil
	}
	r.startWatchSubnets()
	return nil
}

func (r *Role) Stop() {
}
