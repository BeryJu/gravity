package tsdb

import (
	"context"
	"net/http"
	"strconv"
	"sync"
	"time"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/roles"
	debugTypes "beryju.io/gravity/pkg/roles/debug/types"
	"beryju.io/gravity/pkg/roles/tsdb/types"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type Role struct {
	log *log.Entry
	i   roles.Instance
	ctx context.Context
	m   map[string]int
	ms  sync.RWMutex
}

func New(instance roles.Instance) *Role {
	r := &Role{
		log: instance.Log(),
		i:   instance,
		m:   make(map[string]int),
		ms:  sync.RWMutex{},
	}
	r.i.AddEventListener(debugTypes.EventTopicDebugMuxSetup, func(ev *roles.Event) {
		mux := ev.Payload.Data["mux"].(*mux.Router)
		mux.HandleFunc("/debug/tsdb/write", func(w http.ResponseWriter, re *http.Request) {
			r.i.DispatchEvent(types.EventTopicTSDBWrite, roles.NewEvent(
				re.Context(),
				map[string]interface{}{},
			))
		})
	})
	return r
}

func (r *Role) Start(ctx context.Context, config []byte) error {
	r.ctx = ctx
	r.i.AddEventListener(types.EventTopicTSDBWrite, func(ev *roles.Event) {
		r.write()
	})
	r.i.AddEventListener(types.EventTopicTSDBSet, func(ev *roles.Event) {
		r.ms.Lock()
		defer r.ms.Unlock()
		key := ev.Payload.Data["key"].(string)
		r.log.WithField("key", key).Trace("tsdb inc")
		r.m[key] = ev.Payload.Data["value"].(int)
	})
	r.i.AddEventListener(types.EventTopicTSDBInc, func(ev *roles.Event) {
		r.ms.Lock()
		defer r.ms.Unlock()
		key := ev.Payload.Data["key"].(string)
		r.log.WithField("key", key).Trace("tsdb inc")
		val, ok := r.m[key]
		if !ok {
			val = 0
		}
		r.m[key] = val + 1
	})
	go func() {
		for {
			select {
			case <-r.ctx.Done():
				return
			default:
				r.write()
				time.Sleep(30 * time.Second)
			}
		}
	}()
	return nil
}

func (r *Role) write() {
	r.log.Trace("writing metrics")
	r.ms.RLock()
	defer r.ms.RUnlock()
	// Don't bother granting a lease if we don't have any metrics
	if len(r.m) < 1 {
		return
	}
	lease, err := r.i.KV().Grant(r.ctx, 60*30)
	if err != nil {
		r.log.WithError(err).Warning("failed to grant lease, skipping write")
		return
	}
	for key, value := range r.m {
		key := r.i.KV().Key(
			types.KeyRole,
			key,
			extconfig.Get().Instance.Identifier,
			strconv.FormatInt(time.Now().Unix(), 10),
		).String()
		_, err := r.i.KV().Put(
			r.ctx,
			key,
			strconv.Itoa(value),
			clientv3.WithLease(lease.ID),
		)
		if err != nil {
			r.log.WithError(err).WithField("key", key).Warning("failed to put value")
		}
	}
}

func (r *Role) Stop() {
}
