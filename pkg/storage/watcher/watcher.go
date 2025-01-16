package watcher

import (
	"context"
	"errors"
	"strings"
	"sync"
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/storage"
	"github.com/getsentry/sentry-go"
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

	keyFunc func(string) string

	watchCancel context.CancelFunc
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
		log:         extconfig.Get().Logger().Named("storage.watcher").With(zap.String("prefix", prefix.String())),
		constructor: constructor,
		prefix:      prefix,
		client:      client,
		withPrefix:  false,
		keyFunc: func(s string) string {
			return s
		},
	}
	for _, opt := range opts {
		opt(w)
	}
	return w
}

func (w *Watcher[T]) Prefix() *storage.Key {
	return w.prefix.Copy().Prefix(false)
}

func (w *Watcher[T]) Start(ctx context.Context) {
	w.log.Debug("Starting watcher")
	w.loadInitial(ctx)
	cctx, cancel := context.WithCancel(ctx)
	w.watchCancel = cancel
	go w.startWatch(cctx)
}

func (w *Watcher[T]) Stop() {
	if w.watchCancel != nil {
		w.watchCancel()
	}
}

func (w *Watcher[T]) loadInitial(ctx context.Context) {
	w.log.Debug("Loading initial")
	tx := sentry.StartTransaction(ctx, "gravity.storage.watcher.loadInitial")
	defer tx.Finish()
	entries, err := w.client.Get(tx.Context(), w.prefix.String(), clientv3.WithPrefix())
	if err != nil {
		w.log.Warn("failed to list entries", zap.Error(err))
		if !errors.Is(err, context.Canceled) {
			time.Sleep(5 * time.Second)
			w.loadInitial(tx.Context())
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
	watcherEvents.WithLabelValues(w.prefix.String(), mvccpb.Event_EventType_name[int32(t)]).Inc()
	key := w.keyFunc(string(kv.Key))
	// we only care about scope-level updates, everything underneath doesn't matter
	relKey := strings.TrimPrefix(key, w.prefix.String())
	if !w.withPrefix && strings.Contains(relKey, "/") {
		return false
	}
	if w.beforeUpdate != nil {
		w.mutex.RLock()
		old := w.entries[key]
		w.beforeUpdate(old)
		w.mutex.RUnlock()
	}
	if t == mvccpb.DELETE {
		w.log.Debug("removed entry", zap.String("key", key))
		w.mutex.Lock()
		defer w.mutex.Unlock()
		delete(w.entries, key)
	} else if t == mvccpb.PUT {
		e, err := w.constructor(kv)
		if err != nil {
			w.log.Warn("failed to construct entry", zap.Error(err))
			return false
		}
		w.mutex.Lock()
		w.entries[key] = e
		w.mutex.Unlock()
		w.log.Debug("added entry", zap.String("key", key))
	}
	return true
}
