package storage

import (
	"context"
	"time"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/namespace"
)

type Client struct {
	*clientv3.Client
	config clientv3.Config
	prefix string
	log    *zap.Logger
}

func NewClient(prefix string, logger *zap.Logger, endpoints ...string) *Client {
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
	cli.KV = namespace.NewKV(cli.KV, prefix)
	cli.Watcher = namespace.NewWatcher(cli.Watcher, prefix)
	cli.Lease = namespace.NewLease(cli.Lease, prefix)

	return &Client{
		Client: cli,
		log:    logger,
		prefix: prefix,
		config: config,
	}
}

func (c *Client) Config() clientv3.Config {
	return c.config
}

func (c *Client) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	span := sentry.StartSpan(ctx, "etcd.get")
	span.SetTag("etcd.key", key)
	defer span.Finish()
	return c.Client.Get(ctx, key, opts...)
}

func (c *Client) Put(ctx context.Context, key string, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	span := sentry.StartSpan(ctx, "etcd.put")
	span.SetTag("etcd.key", key)
	defer span.Finish()
	return c.Client.Put(ctx, key, val, opts...)
}

func (c *Client) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	span := sentry.StartSpan(ctx, "etcd.delete")
	span.SetTag("etcd.key", key)
	defer span.Finish()
	return c.Client.Delete(ctx, key, opts...)
}
