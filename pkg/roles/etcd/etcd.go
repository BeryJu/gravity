package etcd

import (
	"fmt"
	"net/url"

	"beryju.io/ddet/pkg/extconfig"
	log "github.com/sirupsen/logrus"

	"go.etcd.io/etcd/server/v3/embed"
)

type EmbeddedEtcd struct {
	e      *embed.Etcd
	cfg    *embed.Config
	log    *log.Entry
	prefix string
}

func urlMustParse(raw string) *url.URL {
	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}
	return u
}

func New(prefix string) *EmbeddedEtcd {
	cfg := embed.NewConfig()
	cfg.Dir = "data/etcd/"
	cfg.LogLevel = "warn"
	cfg.LPUrls = []url.URL{
		*urlMustParse(fmt.Sprintf("http://%s:2380", extconfig.Get().Instance.IP)),
	}
	cfg.APUrls = []url.URL{
		*urlMustParse(fmt.Sprintf("http://%s:2380", extconfig.Get().Instance.IP)),
	}
	cfg.Name = extconfig.Get().Instance.Identifier
	cfg.InitialCluster = fmt.Sprintf("%s=http://%s:2380", cfg.Name, extconfig.Get().Instance.IP)
	return &EmbeddedEtcd{
		cfg:    cfg,
		log:    log.WithField("role", "embedded-etcd"),
		prefix: prefix,
	}
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
	// ee.e.Server.Stop()
	ee.e.Close()
}
