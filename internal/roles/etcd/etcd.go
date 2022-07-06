package etcd

import (
	log "github.com/sirupsen/logrus"

	"go.etcd.io/etcd/server/v3/embed"
)

type EmbeddedEtcd struct {
	e      *embed.Etcd
	cfg    *embed.Config
	log    *log.Entry
	prefix string
}

func New(prefix string) *EmbeddedEtcd {
	cfg := embed.NewConfig()
	cfg.Dir = "default.etcd"
	cfg.LogLevel = "warn"
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
	ee.e.Server.Stop()
	ee.e.Close()
}
