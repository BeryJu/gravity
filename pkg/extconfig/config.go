package extconfig

import (
	"os"

	env "github.com/Netflix/go-env"
	log "github.com/sirupsen/logrus"
)

type ExtConfig struct {
	// TODO Change default
	Debug bool `env:"DEBUG,default=true"`
	Etcd  struct {
		Prefix   string `env:"ETCD_PREFIX,default=/ddet"`
		Endpoint string `env:"ETCD_ENDPOINT,default=localhost:2379"`
	}
	BootstrapRoles string `env:"BOOTSTRAP_ROLES,default=dns;dhcp;api;etcd;discovery"`
	Instance       struct {
		Identifier string `env:"INSTANCE_IDENTIFIER"`
		IP         string `env:"INSTANCE_IP"`
	}
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
	if e.Instance.Identifier == "" {
		h, err := os.Hostname()
		if err != nil {
			panic(err)
		}
		e.Instance.Identifier = h
	}
	if e.Instance.IP == "" {
		instIp, err := GetIP()
		if err != nil {
			panic(err)
		}
		e.Instance.IP = instIp.String()
	}
}
