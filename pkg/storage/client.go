package storage

import (
	"time"

	log "github.com/sirupsen/logrus"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/namespace"
)

type Client struct {
	*clientv3.Client
	prefix string
	log    *log.Entry
}

func NewClient(prefix string, endpoints ...string) *Client {
	l := log.WithField("component", "etcd-client")
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:            endpoints,
		DialTimeout:          2 * time.Second,
		DialKeepAliveTime:    2 * time.Second,
		DialKeepAliveTimeout: 2 * time.Second,
	})
	if err != nil {
		l.Panic(err)
	}
	cli.KV = namespace.NewKV(cli.KV, prefix)
	cli.Watcher = namespace.NewWatcher(cli.Watcher, prefix)
	cli.Lease = namespace.NewLease(cli.Lease, prefix)

	return &Client{
		Client: cli,
		log:    l,
		prefix: prefix,
	}
}
