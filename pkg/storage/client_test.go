package storage_test

import (
	"context"
	"errors"
	"testing"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/storage"
	"beryju.io/gravity/pkg/tests"
	"github.com/stretchr/testify/assert"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func TestClient(t *testing.T) {
	tests.Setup(t)
	c := storage.NewClient("/gravity", nil, false, "localhost:2379")
	assert.NotNil(t, c)
	assert.Panics(t, func() {
		storage.NewClient("/gravity", nil, false)
	})
}

func TestClient_Hook_Get(t *testing.T) {
	tests.Setup(t)
	ctx := tests.Context()

	t.Run("empty", func(t *testing.T) {
		tests.ResetEtcd(t)
		c := extconfig.Get().EtcdClient().WithHooks(storage.StorageHook{})
		_, err := c.Get(ctx, "/foo")
		assert.NoError(t, err)
	})

	t.Run("normal", func(t *testing.T) {
		tests.ResetEtcd(t)
		ct := 0
		c := extconfig.Get().EtcdClient().WithHooks(storage.StorageHook{
			GetPre: func(ctx context.Context, key string, opts ...clientv3.OpOption) error {
				ct += 1
				return nil
			},
			GetPost: func(ctx context.Context, key string, res *clientv3.GetResponse, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
				ct += 1
				return res, nil
			},
		})
		_, err := c.Get(ctx, "/foo")
		assert.NoError(t, err)
		assert.Equal(t, 2, ct)
	})

	t.Run("error-pre", func(t *testing.T) {
		tests.ResetEtcd(t)
		ct := 0
		rerr := errors.New("foo")
		c := extconfig.Get().EtcdClient().WithHooks(storage.StorageHook{
			GetPre: func(ctx context.Context, key string, opts ...clientv3.OpOption) error {
				return rerr
			},
			GetPost: func(ctx context.Context, key string, res *clientv3.GetResponse, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
				ct += 1
				return res, nil
			},
		})
		_, err := c.Get(ctx, "/foo")
		assert.Equal(t, rerr, err)
		assert.Equal(t, 0, ct)
	})

	t.Run("error-post", func(t *testing.T) {
		tests.ResetEtcd(t)
		ct := 0
		rerr := errors.New("foo")
		c := extconfig.Get().EtcdClient().WithHooks(storage.StorageHook{
			GetPre: func(ctx context.Context, key string, opts ...clientv3.OpOption) error {
				ct += 1
				return nil
			},
			GetPost: func(ctx context.Context, key string, res *clientv3.GetResponse, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
				ct += 1
				return nil, rerr
			},
		})
		_, err := c.Get(ctx, "/foo")
		assert.Equal(t, rerr, err)
		assert.Equal(t, 2, ct)
	})
}

func TestClient_Hook_Delete(t *testing.T) {
	tests.Setup(t)
	ctx := tests.Context()

	c := extconfig.Get().EtcdClient()
	_, err := c.Put(ctx, "/foo", "bar")
	assert.NoError(t, err)

	t.Run("empty", func(t *testing.T) {
		tests.ResetEtcd(t)
		c := extconfig.Get().EtcdClient().WithHooks(storage.StorageHook{})
		_, err := c.Delete(ctx, "/foo")
		assert.NoError(t, err)
	})

	t.Run("normal", func(t *testing.T) {
		tests.ResetEtcd(t)
		ct := 0
		c := extconfig.Get().EtcdClient().WithHooks(storage.StorageHook{
			DeletePre: func(ctx context.Context, key string, opts ...clientv3.OpOption) error {
				ct += 1
				return nil
			},
			DeletePost: func(ctx context.Context, key string, res *clientv3.DeleteResponse, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
				ct += 1
				return res, nil
			},
		})
		_, err := c.Delete(ctx, "/foo")
		assert.NoError(t, err)
		assert.Equal(t, 2, ct)
	})

	t.Run("error-pre", func(t *testing.T) {
		tests.ResetEtcd(t)
		ct := 0
		rerr := errors.New("foo")
		c := extconfig.Get().EtcdClient().WithHooks(storage.StorageHook{
			DeletePre: func(ctx context.Context, key string, opts ...clientv3.OpOption) error {
				return rerr
			},
			DeletePost: func(ctx context.Context, key string, res *clientv3.DeleteResponse, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
				ct += 1
				return res, nil
			},
		})
		_, err := c.Delete(ctx, "/foo")
		assert.Equal(t, rerr, err)
		assert.Equal(t, 0, ct)
	})

	t.Run("error-post", func(t *testing.T) {
		tests.ResetEtcd(t)
		ct := 0
		rerr := errors.New("foo")
		c := extconfig.Get().EtcdClient().WithHooks(storage.StorageHook{
			DeletePre: func(ctx context.Context, key string, opts ...clientv3.OpOption) error {
				ct += 1
				return nil
			},
			DeletePost: func(ctx context.Context, key string, res *clientv3.DeleteResponse, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
				ct += 1
				return nil, rerr
			},
		})
		_, err := c.Delete(ctx, "/foo")
		assert.Equal(t, rerr, err)
		assert.Equal(t, 2, ct)
	})
}

func TestClient_Hook_Put(t *testing.T) {
	tests.Setup(t)
	ctx := tests.Context()

	t.Run("empty", func(t *testing.T) {
		tests.ResetEtcd(t)
		c := extconfig.Get().EtcdClient().WithHooks(storage.StorageHook{})
		_, err := c.Put(ctx, "/foo", "bar")
		assert.NoError(t, err)
	})

	t.Run("normal", func(t *testing.T) {
		tests.ResetEtcd(t)
		ct := 0
		c := extconfig.Get().EtcdClient().WithHooks(storage.StorageHook{
			PutPre: func(ctx context.Context, key, val string, opts ...clientv3.OpOption) error {
				ct += 1
				return nil
			},
			PutPost: func(ctx context.Context, key string, val string, res *clientv3.PutResponse, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
				ct += 1
				return res, nil
			},
		})
		_, err := c.Put(ctx, "/foo", "bar")
		assert.NoError(t, err)
		assert.Equal(t, 2, ct)
	})

	t.Run("error-pre", func(t *testing.T) {
		tests.ResetEtcd(t)
		ct := 0
		rerr := errors.New("foo")
		c := extconfig.Get().EtcdClient().WithHooks(storage.StorageHook{
			PutPre: func(ctx context.Context, key, val string, opts ...clientv3.OpOption) error {
				return rerr
			},
			PutPost: func(ctx context.Context, key, val string, res *clientv3.PutResponse, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
				ct += 1
				return res, nil
			},
		})
		_, err := c.Put(ctx, "/foo", "bar")
		assert.Equal(t, rerr, err)
		assert.Equal(t, 0, ct)
	})

	t.Run("error-post", func(t *testing.T) {
		tests.ResetEtcd(t)
		ct := 0
		rerr := errors.New("foo")
		c := extconfig.Get().EtcdClient().WithHooks(storage.StorageHook{
			PutPre: func(ctx context.Context, key, val string, opts ...clientv3.OpOption) error {
				ct += 1
				return nil
			},
			PutPost: func(ctx context.Context, key, val string, res *clientv3.PutResponse, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
				ct += 1
				return nil, rerr
			},
		})
		_, err := c.Put(ctx, "/foo", "bar")
		assert.Equal(t, rerr, err)
		assert.Equal(t, 2, ct)
	})
}
