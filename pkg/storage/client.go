package storage

import (
	"time"

	"beryju.io/gravity/pkg/storage/trace"
	"go.uber.org/zap"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/namespace"
)

type Client struct {
	*clientv3.Client
	log    *zap.Logger
	config clientv3.Config
	prefix string
	debug  bool
}

func NewClient(prefix string, logger *zap.Logger, debug bool, endpoints ...string) *Client {
	config := clientv3.Config{
		Endpoints:            endpoints,
		DialTimeout:          2 * time.Second,
		DialKeepAliveTime:    2 * time.Second,
		DialKeepAliveTimeout: 2 * time.Second,
		Logger:               logger,
	}
	cli, err := clientv3.New(config)
	if err != nil {
		logger.Panic("failed to setup etcd client", zap.Error(err))
	}
	cli.KV = trace.NewKV(namespace.NewKV(cli.KV, prefix), func(op clientv3.Op) {
		if debug {
			logger.Warn("etcd op without transaction", zap.String("key", string(op.KeyBytes())), zap.String("op", trace.NameFromOp(op)))
		}
	})
	cli.Watcher = namespace.NewWatcher(cli.Watcher, prefix)
	cli.Lease = namespace.NewLease(cli.Lease, prefix)

	return &Client{
		Client: cli,
		log:    logger,
		prefix: prefix,
		config: config,
		debug:  debug,
	}
}

func (c *Client) Config() clientv3.Config {
	return c.config
}

func (c *Client) AddGlobalHook() {

}
