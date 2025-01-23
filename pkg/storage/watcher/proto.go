package watcher

import (
	"beryju.io/gravity/pkg/storage"
	"go.etcd.io/etcd/api/v3/mvccpb"
	"google.golang.org/protobuf/proto"
)

func NewProto[T proto.Message](
	client *storage.Client,
	prefix *storage.Key,
	opts ...func(w *Watcher[T]),
) *Watcher[T] {
	_ = new(T)
	return New(
		func(kv *mvccpb.KeyValue) (T, error) {
			obj := new(T)
			err := proto.Unmarshal(kv.Value, *obj)
			if err != nil {
				return *obj, err
			}
			return *obj, nil
		},
		client, prefix, opts...,
	)
}
