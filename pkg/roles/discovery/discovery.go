package discovery

import (
	"context"

	"beryju.io/ddet/pkg/extconfig"
	"beryju.io/ddet/pkg/roles"
	apitypes "beryju.io/ddet/pkg/roles/api/types"
	"beryju.io/ddet/pkg/roles/discovery/types"

	log "github.com/sirupsen/logrus"
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
	r.i.AddEventListener(apitypes.EventTopicAPIMuxSetup, r.eventHandlerAPIMux)
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
