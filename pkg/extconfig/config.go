package extconfig

import (
	"os"

	env "github.com/Netflix/go-env"
	log "github.com/sirupsen/logrus"
)

type ExtConfig struct {
	Debug    bool   `env:"DEBUG,default=false"`
	DataPath string `env:"DATA_PATH,default=./data"`
	Etcd     struct {
		Prefix    string `env:"ETCD_PREFIX,default=/ddet"`
		Endpoint  string `env:"ETCD_ENDPOINT,default=localhost:2379"`
		Discovery string `env:"ETCD_DISCOVERY"`
	}
	BootstrapRoles string `env:"BOOTSTRAP_ROLES,default=dns;dhcp;api;etcd;discovery"`
	Instance       struct {
		Identifier string `env:"INSTANCE_IDENTIFIER"`
		IP         string `env:"INSTANCE_IP"`
	}
	ListenOnlyMode bool `env:"LISTEN_ONLY,default=false"`
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
