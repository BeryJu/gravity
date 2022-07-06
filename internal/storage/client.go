package storage

import (
	"strings"
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

func NewClient(endpoint string, prefix string) *Client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{endpoint},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		panic(err)
	}
	cli.KV = namespace.NewKV(cli.KV, prefix)
	cli.Watcher = namespace.NewWatcher(cli.Watcher, prefix)
	cli.Lease = namespace.NewLease(cli.Lease, prefix)

	return &Client{
		Client: cli,
		log:    log.WithField("component", "etcd-client"),
		prefix: prefix,
	}
}

func (c *Client) Key(parts ...string) string {
	b := strings.Builder{}
	b.WriteString("/")
	for idx, part := range parts {
		b.WriteString(part)
		if idx != len(parts)-1 {
			b.WriteString("/")
		}
	}
	return b.String()
}
