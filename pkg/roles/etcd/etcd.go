package etcd

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"beryju.io/ddet/pkg/extconfig"
	"beryju.io/ddet/pkg/roles"
	apitypes "beryju.io/ddet/pkg/roles/api/types"
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
	dirs := extconfig.Get().Dirs()
	cfg := embed.NewConfig()
	cfg.Dir = dirs.EtcdDir
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
		etcdDir: dirs.EtcdDir,
		certDir: dirs.CertDir,
	}
	cfg.InitialCluster = fmt.Sprintf("%s=https://%s:2380", cfg.Name, extconfig.Get().Instance.IP)
	if extconfig.Get().Etcd.JoinCluster != "" {
		cfg.ClusterState = embed.ClusterStateFlagExisting
		joinParts := strings.Split(extconfig.Get().Etcd.JoinCluster, ";")
		cfg.InitialCluster = fmt.Sprintf(
			"%s,%s=https://%s:2380",
			cfg.InitialCluster,
			joinParts[0],
			joinParts[1],
		)
	}
	ee.log.Trace(cfg.InitialCluster)
	cfg.PeerAutoTLS = true
	cfg.PeerTLSInfo.ClientCertFile = path.Join(ee.certDir, relInstCertPath)
	cfg.PeerTLSInfo.ClientKeyFile = path.Join(ee.certDir, relInstKeyPath)
	cfg.PeerTLSInfo.ClientCertAuth = true
	cfg.SelfSignedCertValidity = 1
	ee.i.AddEventListener(apitypes.EventTopicAPIMuxSetup, ee.eventHandlerAPIMux)
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
