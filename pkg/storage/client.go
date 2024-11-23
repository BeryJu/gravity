package storage

import (
	"context"
	"time"

	"beryju.io/gravity/pkg/storage/trace"
	"go.uber.org/zap"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/namespace"
)

type StorageHook struct {
	Request func(ctx context.Context, op clientv3.Op, client *Client) error
}

type Client struct {
	*clientv3.Client
	log    *zap.Logger
	config clientv3.Config
	prefix string
	debug  bool
	hooks  []StorageHook
	parent *Client
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
		hooks:  []StorageHook{},
	}
}

func (c *Client) Config() clientv3.Config {
	return c.config
}

func (c *Client) WithHooks(hooks ...StorageHook) *Client {
	return &Client{
		Client: c.Client,
		log:    c.log,
		prefix: c.prefix,
		config: c.config,
		debug:  c.debug,
		hooks:  hooks,
		parent: c,
	}
}

func (c *Client) Do(ctx context.Context, op clientv3.Op) (clientv3.OpResponse, error) {
	for _, h := range c.hooks {
		err := h.Request(ctx, op, c.parent)
		if err != nil {
			return clientv3.OpResponse{}, err
		}
	}
	return c.Client.Do(ctx, op)
}
