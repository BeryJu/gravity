package watcher

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"beryju.io/gravity/pkg/storage"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

type Watcher[T any] struct {
	entries     map[string]T
	mutex       sync.RWMutex
	log         *zap.Logger
	constructor func(*mvccpb.KeyValue) (T, error)
	prefix      *storage.Key
	client      *storage.Client

	withPrefix       bool
	afterInitialLoad func()
	beforeUpdate     func(entry T)
}

func WithPrefix[T any]() func(*Watcher[T]) {
	return func(w *Watcher[T]) {
		w.withPrefix = true
	}
}

func WithAfterInitialLoad[T any](callback func()) func(*Watcher[T]) {
	return func(w *Watcher[T]) {
		w.afterInitialLoad = callback
	}
}

func WithBeforeUpdate[T any](callback func(entry T)) func(*Watcher[T]) {
	return func(w *Watcher[T]) {
		w.beforeUpdate = callback
	}
}

func New[T any](
	constructor func(*mvccpb.KeyValue) (T, error),
	client *storage.Client,
	prefix *storage.Key,
	opts ...func(w *Watcher[T]),
) *Watcher[T] {
	w := &Watcher[T]{
		entries:     make(map[string]T),
		mutex:       sync.RWMutex{},
		constructor: constructor,
		prefix:      prefix,
		client:      client,
		withPrefix:  false,
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func (w *Watcher[T]) Get(key string) T {
	w.mutex.RLock()
	defer w.mutex.RUnlock()
	return w.entries[key]
}

type KV[T any] struct {
	Key   string
	Value T
}

func (w *Watcher[T]) Iter() chan KV[T] {
	c := make(chan KV[T])
	go func() {
		w.mutex.RLock()
		defer w.mutex.RUnlock()
		for k, v := range w.entries {
			select {
			case c <- KV[T]{Key: k, Value: v}:
			case <-c:
				close(c)
				return
			}
		}
		close(c)
	}()
	return c
}

func (w *Watcher[T]) Start(ctx context.Context) {
	w.loadInitial(ctx)
	w.startWatch(ctx)
}

func (w *Watcher[T]) loadInitial(ctx context.Context) {
	entries, err := w.client.Get(
		ctx, w.prefix.String(), clientv3.WithPrefix())
	if err != nil {
		w.log.Warn("failed to list entries", zap.Error(err))
		if !errors.Is(err, context.Canceled) {
			time.Sleep(5 * time.Second)
			w.loadInitial(ctx)
		}
		return
	}
	for _, entry := range entries.Kvs {
		w.handleEvent(mvccpb.PUT, entry)
	}
	if w.afterInitialLoad != nil {
		w.afterInitialLoad()
	}
}

func (w *Watcher[T]) startWatch(ctx context.Context) {
	ch := w.client.Watch(ctx, w.prefix.String(), clientv3.WithPrefix())
	for watchResp := range ch {
		for _, event := range watchResp.Events {
			w.handleEvent(event.Type, event.Kv)
		}
	}
}

func (w *Watcher[T]) handleEvent(t mvccpb.Event_EventType, kv *mvccpb.KeyValue) bool {
	relKey := strings.TrimPrefix(string(kv.Key), w.prefix.String())
	// we only care about scope-level updates, everything underneath doesn't matter
	if !w.withPrefix && strings.Contains(relKey, "/") {
		return false
	}
	if w.beforeUpdate != nil {
		w.mutex.RLock()
		old := w.entries[string(kv.Key)]
		w.beforeUpdate(old)
		w.mutex.RUnlock()
	}
	if t == mvccpb.DELETE {
		w.log.Debug("removed entry", zap.String("key", string(kv.Key)))
		w.mutex.Lock()
		defer w.mutex.Unlock()
		delete(w.entries, string(kv.Key))
	} else if t == mvccpb.PUT {
		e, err := w.constructor(kv)
		if err != nil {
			w.log.Warn("failed to construct entry", zap.Error(err))
			return false
		}
		// s.calculateUsage()
		w.mutex.Lock()
		w.entries[string(kv.Key)] = e
		w.mutex.Unlock()
		w.log.Debug("added entry", zap.String("key", string(kv.Key)))
	}
	return true
}
