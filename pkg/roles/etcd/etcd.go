package etcd

import (
	"fmt"
	"net/url"
	"path"

	"beryju.io/ddet/pkg/extconfig"
	"beryju.io/ddet/pkg/roles"
	log "github.com/sirupsen/logrus"

	"go.etcd.io/etcd/server/v3/embed"
)

type EmbeddedEtcd struct {
	etcdDir string
	certDir string

	e   *embed.Etcd
	cfg *embed.Config
	log *log.Entry
	i   roles.Instance
}

func urlMustParse(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return u
}

func New(instance roles.Instance) *EmbeddedEtcd {
	etcdDir := path.Join(extconfig.Get().DataPath, "etcd/")
	certDir := path.Join(extconfig.Get().DataPath, "cert/")
	cfg := embed.NewConfig()
	cfg.Dir = etcdDir
	cfg.LogLevel = "warn"
	cfg.LPUrls = []url.URL{
		*urlMustParse(fmt.Sprintf("https://%s:2380", extconfig.Get().Instance.IP)),
	}
	cfg.APUrls = []url.URL{
		*urlMustParse(fmt.Sprintf("https://%s:2380", extconfig.Get().Instance.IP)),
	}
	cfg.Name = extconfig.Get().Instance.Identifier
	cfg.InitialCluster = ""
	ee := &EmbeddedEtcd{
		cfg:     cfg,
		log:     instance.GetLogger().WithField("role", "embedded-etcd"),
		i:       instance,
		etcdDir: etcdDir,
		certDir: certDir,
	}
	cfg.InitialCluster = fmt.Sprintf("%s=https://%s:2380", cfg.Name, extconfig.Get().Instance.IP)
	if extconfig.Get().Etcd.JoinCluster != "" {
		cfg.ClusterState = "existing"
		cfg.InitialCluster = fmt.Sprintf(
			"%s,%[2]s=https://%[2]s:2380",
			cfg.InitialCluster,
			extconfig.Get().Etcd.JoinCluster,
		)
	}
	cfg.PeerAutoTLS = true
	cfg.PeerTLSInfo.ClientCertFile = path.Join(certDir, relInstCertPath)
	cfg.PeerTLSInfo.ClientKeyFile = path.Join(certDir, relInstKeyPath)
	cfg.PeerTLSInfo.ClientCertAuth = true
	cfg.SelfSignedCertValidity = 1
	return ee
}

func (ee *EmbeddedEtcd) Start(ready func()) error {
	e, err := embed.StartEtcd(ee.cfg)
	if err != nil {
		return err
	}
	ee.e = e
	go func() {
		<-e.Server.ReadyNotify()
		ee.log.Info("Embedded etcd Ready!")
		ready()
	}()
	return <-e.Err()
}

func (ee *EmbeddedEtcd) Stop() {
	if ee.e == nil {
		return
	}
	ee.log.Info("Stopping etcd")
	ee.e.Close()
}
