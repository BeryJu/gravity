package etcd

import (
	"context"
	"fmt"
	"net/url"
	"path"
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"

	"go.etcd.io/etcd/server/v3/embed"
)

const (
	relInstCertPath = "/instance.pem"
	relInstKeyPath  = "/instance_key.pem"
)

type Role struct {
	etcdDir string
	certDir string

	e   *embed.Etcd
	cfg *embed.Config
	log *zap.Logger
	i   roles.Instance
}

func urlMustParse(raw string) url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return *u
}

func New(instance roles.Instance) *Role {
	dirs := extconfig.Get().Dirs()
	cfg := embed.NewConfig()
	cfg.Dir = dirs.EtcdDir
	cfg.LogLevel = "warn"
	cfg.ZapLoggerBuilder = embed.NewZapCoreLoggerBuilder(
		extconfig.Get().BuildLoggerWithLevel(zapcore.WarnLevel).Named("etcd"),
		nil,
		nil,
	)
	cfg.AutoCompactionMode = "periodic"
	cfg.AutoCompactionRetention = "60m"
	cfg.LPUrls = []url.URL{
		urlMustParse(fmt.Sprintf("https://%s", extconfig.Get().Listen(2380))),
	}
	cfg.APUrls = []url.URL{
		urlMustParse(fmt.Sprintf("https://%s:2380", extconfig.Get().Instance.IP)),
	}
	cfg.Name = extconfig.Get().Instance.Identifier
	ee := &Role{
		cfg:     cfg,
		log:     instance.Log(),
		i:       instance,
		etcdDir: dirs.EtcdDir,
		certDir: dirs.CertDir,
	}
	cfg.InitialCluster = fmt.Sprintf("%s=https://%s:2380", cfg.Name, extconfig.Get().Instance.IP)
	if extconfig.Get().Etcd.JoinCluster != "" {
		cfg.ClusterState = embed.ClusterStateFlagExisting
		cfg.InitialCluster = cfg.InitialCluster + "," + extconfig.Get().Etcd.JoinCluster
	}
	cfg.PeerAutoTLS = true
	cfg.PeerTLSInfo.ClientCertFile = path.Join(ee.certDir, relInstCertPath)
	cfg.PeerTLSInfo.ClientKeyFile = path.Join(ee.certDir, relInstKeyPath)
	cfg.PeerTLSInfo.ClientCertAuth = true
	cfg.SelfSignedCertValidity = 1
	return ee
}

func (ee *Role) Start(ctx context.Context, config []byte) error {
	start := time.Now()
	e, err := embed.StartEtcd(ee.cfg)
	if err != nil {
		return err
	}
	ee.e = e
	go func() {
		err := <-e.Err()
		if err != nil {
			ee.log.Warn("failed to start/stop etcd", zap.Error(err))
		}
	}()
	<-e.Server.ReadyNotify()
	ee.log.Info("Embedded etcd Ready!", zap.Duration("runtime", time.Since(start)))
	return nil
}

func (ee *Role) Stop() {
	if ee.e == nil {
		return
	}
	ee.log.Info("Stopping etcd")
	ee.e.Close()
}
