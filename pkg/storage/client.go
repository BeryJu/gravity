package storage

import (
	"context"
	"encoding/json"
	"time"

	"beryju.io/gravity/pkg/storage/trace"
	"go.uber.org/zap"
	"google.golang.org/protobuf/proto"

	clientv3 "go.etcd.io/etcd/client/v3"
	"go.etcd.io/etcd/client/v3/namespace"
)

type StorageHook struct {
	GetPre    func(ctx context.Context, key string, opts ...clientv3.OpOption) error
	PutPre    func(ctx context.Context, key string, val string, opts ...clientv3.OpOption) error
	DeletePre func(ctx context.Context, key string, opts ...clientv3.OpOption) error

	GetPost    func(ctx context.Context, key string, res *clientv3.GetResponse, opts ...clientv3.OpOption) (*clientv3.GetResponse, error)
	PutPost    func(ctx context.Context, key string, val string, res *clientv3.PutResponse, opts ...clientv3.OpOption) (*clientv3.PutResponse, error)
	DeletePost func(ctx context.Context, key string, res *clientv3.DeleteResponse, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error)
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
			logger.DPanic("etcd op without transaction", zap.String("key", string(op.KeyBytes())), zap.String("op", trace.NameFromOp(op)))
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

func (c *Client) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	for _, h := range c.hooks {
		if h.GetPre == nil {
			continue
		}
		err := h.GetPre(ctx, key, opts...)
		if err != nil {
			return nil, err
		}
	}
	res, err := c.Client.Get(ctx, key, opts...)
	if err != nil {
		return res, err
	}
	for _, h := range c.hooks {
		if h.GetPost == nil {
			continue
		}
		_r, err := h.GetPost(ctx, key, res, opts...)
		if err != nil {
			return nil, err
		}
		res = _r
	}
	return res, err
}

func (c *Client) Put(ctx context.Context, key string, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	for _, h := range c.hooks {
		if h.PutPre == nil {
			continue
		}
		err := h.PutPre(ctx, key, val, opts...)
		if err != nil {
			return nil, err
		}
	}
	res, err := c.Client.Put(ctx, key, val, opts...)
	if err != nil {
		return res, err
	}
	for _, h := range c.hooks {
		if h.PutPost == nil {
			continue
		}
		_r, err := h.PutPost(ctx, key, val, res, opts...)
		if err != nil {
			return nil, err
		}
		res = _r
	}
	return res, err
}

func (c *Client) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	for _, h := range c.hooks {
		if h.DeletePre == nil {
			continue
		}
		err := h.DeletePre(ctx, key, opts...)
		if err != nil {
			return nil, err
		}
	}
	res, err := c.Client.Delete(ctx, key, opts...)
	if err != nil {
		return res, err
	}
	for _, h := range c.hooks {
		if h.DeletePost == nil {
			continue
		}
		_r, err := h.DeletePost(ctx, key, res, opts...)
		if err != nil {
			return nil, err
		}
		res = _r
	}
	return res, err
}

func (c *Client) PutObj(ctx context.Context, key string, val any, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	marshalled, err := c.Marshal(val)
	if err != nil {
		return nil, err
	}
	return c.Put(ctx, key, string(marshalled), opts...)
}

func (c *Client) Unmarshal(raw []byte, out any) error {
	err := json.Unmarshal(raw, out)
	if err != nil {
		if p, ok := out.(proto.Message); ok {
			return proto.Unmarshal(raw, p)
		}
	}
	return err
}

func (c *Client) Marshal(in any) ([]byte, error) {
	if p, ok := in.(proto.Message); ok {
		return proto.Marshal(p)
	}
	return json.Marshal(in)
}
