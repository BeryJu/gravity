package extconfig

import (
	"os"

	env "github.com/Netflix/go-env"
	log "github.com/sirupsen/logrus"
)

type ExtConfig struct {
	Etcd struct {
		Prefix   string `env:"ETCD_PREFIX,default=/ddet"`
		Endpoint string `env:"ETCD_ENDPOINT,default=localhost:2379"`
	}
	BootstrapRoles     string `env:"BOOTSTRAP_ROLES,default=dns;api;etcd'"`
	InstanceIdentifier string `env:"INSTANCE_IDENTIFIER"`
}

var globalExtConfig *ExtConfig

func Get() *ExtConfig {
	if globalExtConfig != nil {
		return globalExtConfig
	}
	var cfg ExtConfig
	_, err := env.UnmarshalFromEnviron(&cfg)
	if err != nil {
		log.WithError(err).Warning("failed to load external config")
		return nil
	}
	cfg.defaults()
	globalExtConfig = &cfg
	return &cfg
}

func (e *ExtConfig) defaults() {
	if e.InstanceIdentifier == "" {
		h, err := os.Hostname()
		if err != nil {
			panic(err)
		}
		e.InstanceIdentifier = h
	}
}
