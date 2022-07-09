package discovery

import (
	"beryju.io/ddet/pkg/roles"
	"beryju.io/ddet/pkg/roles/discovery/types"

	log "github.com/sirupsen/logrus"
)

type DiscoveryRole struct {
	log *log.Entry
	i   roles.Instance
	cfg *DiscoveryRoleConfig
}

func New(instance roles.Instance) *DiscoveryRole {
	return &DiscoveryRole{
		log: instance.GetLogger().WithField("role", types.KeyRole),
		i:   instance,
	}
}

func (r *DiscoveryRole) Start(config []byte) error {
	r.cfg = r.decodeDiscoveryRoleConfig(config)
	if !r.cfg.Enabled {
		return nil
	}
	r.startWatchSubnets()
	return nil
}

func (r *DiscoveryRole) Stop() {
}
