package etcd

import (
	"fmt"
	"net/url"
	"path"
	"strings"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	apitypes "beryju.io/gravity/pkg/roles/api/types"
	log "github.com/sirupsen/logrus"
	"github.com/swaggest/rest/web"

	"go.etcd.io/etcd/server/v3/embed"
)

type Role struct {
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

func New(instance roles.Instance) *Role {
	dirs := extconfig.Get().Dirs()
	cfg := embed.NewConfig()
	cfg.Dir = dirs.EtcdDir
	cfg.LogLevel = "warn"
	cfg.AutoCompactionMode = "periodic"
	cfg.AutoCompactionRetention = "60m"
	cfg.LPUrls = []url.URL{
		*urlMustParse(fmt.Sprintf("https://%s", extconfig.Get().Listen(2380))),
	}
	cfg.APUrls = []url.URL{
		*urlMustParse(fmt.Sprintf("https://%s:2380", extconfig.Get().Instance.IP)),
	}
	cfg.Name = extconfig.Get().Instance.Identifier
	cfg.InitialCluster = ""
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
	ee.i.AddEventListener(apitypes.EventTopicAPIMuxSetup, func(ev *roles.Event) {
		svc := ev.Payload.Data["svc"].(*web.Service)
		svc.Get("/api/v1/etcd/members", ee.apiHandlerMembers())
		svc.Post("/api/v1/etcd/join", ee.apiHandlerJoin())
	})
	return ee
}

func (ee *Role) Start(ready func()) error {
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

func (ee *Role) Stop() {
	if ee.e == nil {
		return
	}
	ee.log.Info("Stopping etcd")
	ee.e.Close()
}
