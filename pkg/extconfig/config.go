package extconfig

import (
	"fmt"
	"os"
	"path"

	"beryju.io/gravity/pkg/storage"
	env "github.com/Netflix/go-env"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
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
	BootstrapRoles string `env:"BOOTSTRAP_ROLES,default=dns;dhcp;api;etcd;discovery;backup;monitoring;tsdb"`
	Instance       struct {
		Identifier string `env:"INSTANCE_IDENTIFIER"`
		IP         string `env:"INSTANCE_IP"`
		Listen     string `env:"INSTANCE_LISTEN"`
	}
	ListenOnlyMode bool   `env:"LISTEN_ONLY,default=false"`
	FallbackDNS    string `env:"FALLBACK_DNS,default=1.1.1.1:53"`

	logger *zap.Logger
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
		e.Etcd.Endpoint,
	)
}

func (e *ExtConfig) Listen(port int32) string {
	listen := e.Instance.IP
	if e.Instance.Listen != "" {
		listen = e.Instance.Listen
	}
	return fmt.Sprintf("%s:%d", listen, port)
}

func (e *ExtConfig) Logger() *zap.Logger {
	return e.logger
}

func (e *ExtConfig) BuildLogger() *zap.Logger {
	l, err := zapcore.ParseLevel(e.LogLevel)
	if err != nil {
		l = zapcore.InfoLevel
	}
	if e.Debug {
		l = zapcore.DebugLevel
	}
	return e.BuildLoggerWithLevel(l)
}

func (e *ExtConfig) BuildLoggerWithLevel(l zapcore.Level) *zap.Logger {
	config := zap.Config{
		Encoding:         "json",
		Development:      false,
		OutputPaths:      []string{"stderr"},
		ErrorOutputPaths: []string{"stderr"},
		EncoderConfig:    zap.NewProductionEncoderConfig(),
	}
	config.Level = zap.NewAtomicLevelAt(l)
	config.DisableCaller = !e.Debug
	if e.Debug {
		config.Development = true
		config.Encoding = "console"
		config.EncoderConfig = zap.NewDevelopmentEncoderConfig()
		config.EncoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}
	config.EncoderConfig.EncodeDuration = zapcore.MillisDurationEncoder
	log, err := config.Build()
	if err != nil {
		panic(err)
	}
	return log
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
}
