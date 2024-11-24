package trace

import (
	"context"

	"github.com/getsentry/sentry-go"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type opWithoutSpan func(op clientv3.Op)

type traceKV struct {
	clientv3.KV
	opWithoutSpan opWithoutSpan
}

func NewKV(c clientv3.KV, opWithoutSpan opWithoutSpan) clientv3.KV {
	return traceKV{c, opWithoutSpan}
}

func NameFromOp(op clientv3.Op) string {
	if op.IsGet() {
		return "etcd.get"
	} else if op.IsPut() {
		return "etcd.put"
	} else if op.IsDelete() {
		return "etcd.delete"
	} else {
		return "etcd.unknown"
	}
}

func (kv traceKV) trace(ctx context.Context, op clientv3.Op) func() {
	tx := sentry.TransactionFromContext(ctx)
	if tx == nil {
		kv.opWithoutSpan(op)
		return func() {}
	}
	span := tx.StartChild(NameFromOp(op))
	span.Description = string(op.KeyBytes())
	span.SetTag("etcd.key", span.Description)
	return func() {
		span.Finish()
	}
}

func (kv traceKV) Get(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	defer kv.trace(ctx, clientv3.OpGet(key, opts...))()
	return kv.KV.Get(ctx, key, opts...)
}

func (kv traceKV) Put(ctx context.Context, key string, val string, opts ...clientv3.OpOption) (*clientv3.PutResponse, error) {
	defer kv.trace(ctx, clientv3.OpPut(key, val, opts...))()
	return kv.KV.Put(ctx, key, val, opts...)
}

func (kv traceKV) Delete(ctx context.Context, key string, opts ...clientv3.OpOption) (*clientv3.DeleteResponse, error) {
	defer kv.trace(ctx, clientv3.OpDelete(key, opts...))()
	return kv.KV.Delete(ctx, key, opts...)
}
