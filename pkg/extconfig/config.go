package extconfig

import (
	"fmt"
	"os"
	"path"

	"beryju.io/gravity/pkg/storage"
	env "github.com/Netflix/go-env"
	log "github.com/sirupsen/logrus"
)

type ExtConfig struct {
	Debug    bool   `env:"DEBUG,default=false"`
	LogLevel string `env:"LOG_LEVEL,default=info"`
	DataPath string `env:"DATA_PATH,default=./data"`
	Etcd     struct {
		Prefix      string `env:"ETCD_PREFIX,default=/gravity"`
		Endpoint    string `env:"ETCD_ENDPOINT,default=localhost:2379"`
		JoinCluster string `env:"ETCD_JOIN_CLUSTER"`
	}
	BootstrapRoles string `env:"BOOTSTRAP_ROLES,default=dns;dhcp;api;etcd;discovery;backup;monitoring"`
	Instance       struct {
		Identifier string `env:"INSTANCE_IDENTIFIER"`
		IP         string `env:"INSTANCE_IP"`
		Listen     string `env:"INSTANCE_LISTEN"`
	}
	ListenOnlyMode bool   `env:"LISTEN_ONLY,default=false"`
	FallbackDNS    string `env:"FALLBACK_DNS,default=1.1.1.1:53"`
}

type ExtConfigDirs struct {
	EtcdDir   string `json:"etcdDir"`
	CertDir   string `json:"certDir"`
	BackupDir string `json:"backupDir"`
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

func (e *ExtConfig) EtcdClient() *storage.Client {
	return storage.NewClient(e.Etcd.Prefix, e.Etcd.Endpoint)
}

func (e *ExtConfig) Listen(port int32) string {
	listen := e.Instance.IP
	if e.Instance.Listen != "" {
		listen = e.Instance.Listen
	}
	return fmt.Sprintf("%s:%d", listen, port)
}

func (e *ExtConfig) defaults() {
	if e.Debug {
		log.SetLevel(log.TraceLevel)
		log.SetFormatter(&log.TextFormatter{})
	} else {
		l, err := log.ParseLevel(e.LogLevel)
		if err != nil {
			l = log.WarnLevel
		}
		log.SetLevel(l)
		log.SetFormatter(&log.JSONFormatter{
			DisableHTMLEscape: true,
		})
	}
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
