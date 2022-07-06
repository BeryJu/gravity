package discovery

import (
	"beryju.io/ddet/pkg/roles"

	log "github.com/sirupsen/logrus"
)

const (
	KeyRole = "discovery"
)

type DiscoveryRole struct {
	log *log.Entry
	i   roles.Instance
}

func New(instance roles.Instance) *DiscoveryRole {
	return &DiscoveryRole{
		log: log.WithField("role", "discovery"),
		i:   instance,
	}
}

func (r *DiscoveryRole) Start(config []byte) error {
	return nil
}

func (r *DiscoveryRole) Stop() {
}
