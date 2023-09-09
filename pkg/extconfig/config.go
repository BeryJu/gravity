package extconfig

import (
	"fmt"
	"net"
	"os"
	"path"

	"beryju.io/gravity/pkg/storage"
	env "github.com/Netflix/go-env"
	"go.uber.org/zap"
)

type ExtConfig struct {
	logger *zap.Logger
	Etcd   struct {
		Prefix      string `env:"ETCD_PREFIX,default=/gravity"`
		Endpoint    string `env:"ETCD_ENDPOINT,default=localhost:2379"`
		PeerPort    int32  `env:"ETCD_PEER_PORT,default=2380"`
		JoinCluster string `env:"ETCD_JOIN_CLUSTER"`
	}
	Instance struct {
		Identifier string `env:"INSTANCE_IDENTIFIER"`
		IP         string `env:"INSTANCE_IP"`
		Interface  string `env:"INSTANCE_INTERFACE"`
		Listen     string `env:"INSTANCE_LISTEN"`
	}
	LogLevel       string `env:"LOG_LEVEL,default=info"`
	DataPath       string `env:"DATA_PATH,default=./data"`
	BootstrapRoles string `env:"BOOTSTRAP_ROLES,default=dns;dhcp;api;etcd;discovery;backup;monitoring;tsdb"`
	FallbackDNS    string `env:"FALLBACK_DNS,default=1.1.1.1:53"`

	Sentry struct {
		DSN     string `env:"SENTRY_DSN,default=https://731a93aa4a1a42a2960ac9eecee628c5@sentry.beryju.org/2"`
		Enabled bool   `env:"SENTRY_ENABLED,default=false"`
	}

	Debug          bool `env:"DEBUG,default=false"`
	ListenOnlyMode bool `env:"LISTEN_ONLY,default=false"`
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
		panic(err)
	}
	cfg.load()
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
	return storage.NewClient(
		e.Etcd.Prefix,
		e.BuildLoggerWithLevel(zap.WarnLevel).Named("etcd.client"),
		e.Debug,
		e.Etcd.Endpoint,
	)
}

func Listen(addr string, port int32) string {
	ip := net.ParseIP(addr)
	if ip.To4() != nil {
		return fmt.Sprintf("%s:%d", ip.String(), port)
	}
	return fmt.Sprintf("[%s]:%d", ip.String(), port)
}

func (e *ExtConfig) Listen(port int32) string {
	listen := e.Instance.IP
	if e.Instance.Listen != "" {
		listen = e.Instance.Listen
	}
	return Listen(listen, port)
}

func (e *ExtConfig) load() {
	e.logger = e.BuildLogger()
	if e.Instance.Identifier == "" {
		h, err := os.Hostname()
		if err != nil {
			panic(err)
		}
		e.Instance.Identifier = h
	}
	if e.Instance.IP == "" {
		instIp, err := e.GetIP()
		if err != nil {
			panic(err)
		}
		e.Instance.IP = instIp.String()
	}
	if e.Instance.Interface == "" {
		i, err := e.GetInterfaceForIP(net.ParseIP(e.Instance.IP))
		if err != nil || i == nil {
			e.logger.Warn("defaulting to all interfaces", zap.Error(err))
			e.Instance.Interface = ""
		} else {
			e.Instance.Interface = i.Name
		}
	}
	e.logger = e.BuildLogger()
}
