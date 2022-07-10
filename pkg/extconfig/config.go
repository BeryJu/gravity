package extconfig

import (
	"fmt"
	"os"
	"path"

	env "github.com/Netflix/go-env"
	log "github.com/sirupsen/logrus"
)

type ExtConfig struct {
	Debug    bool   `env:"DEBUG,default=false"`
	DataPath string `env:"DATA_PATH,default=./data"`
	Etcd     struct {
		Prefix      string `env:"ETCD_PREFIX,default=/ddet"`
		Endpoint    string `env:"ETCD_ENDPOINT,default=localhost:2379"`
		JoinCluster string `env:"ETCD_JOIN_CLUSTER"`
	}
	BootstrapRoles string `env:"BOOTSTRAP_ROLES,default=dns;dhcp;api;etcd;discovery;backup"`
	Instance       struct {
		Identifier string `env:"INSTANCE_IDENTIFIER"`
		IP         string `env:"INSTANCE_IP"`
		Listen     string `env:"INSTANCE_LISTEN"`
	}
	ListenOnlyMode bool `env:"LISTEN_ONLY,default=false"`
}

type ExtConfigDirs struct {
	EtcdDir   string
	CertDir   string
	BackupDir string
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

func (e *ExtConfig) Dirs() *ExtConfigDirs {
	return &ExtConfigDirs{
		EtcdDir:   path.Join(e.DataPath, "etcd/"),
		CertDir:   path.Join(e.DataPath, "cert/"),
		BackupDir: path.Join(e.DataPath, "backup/"),
	}
}

func (e *ExtConfig) Listen(port int32) string {
	listen := e.Instance.IP
	if e.Instance.Listen != "" {
		listen = e.Instance.Listen
	}
	return fmt.Sprintf("%s:%d", listen, port)
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
