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
		return "get"
	} else if op.IsPut() {
		return "put"
	} else if op.IsDelete() {
		return "delete"
	} else {
		return "unknown"
	}
}

func (kv traceKV) trace(ctx context.Context, op clientv3.Op) func() {
	tx := sentry.TransactionFromContext(ctx)
	if tx == nil {
		kv.opWithoutSpan(op)
		return func() {}
	}
	span := tx.StartChild("db")
	span.Description = string(op.KeyBytes())
	span.SetTag("db.system", "other_sql")
	span.SetTag("db.operation.name", NameFromOp(op))
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
