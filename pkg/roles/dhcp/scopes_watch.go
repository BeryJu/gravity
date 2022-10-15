package dhcp

import (
	"strings"
	"time"

	"beryju.io/gravity/pkg/roles/dhcp/types"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

func (r *Role) handleScopeOp(t mvccpb.Event_EventType, kv *mvccpb.KeyValue) bool {
	prefix := r.i.KV().Key(types.KeyRole, types.KeyScopes).Prefix(true)
	relKey := strings.TrimPrefix(string(kv.Key), prefix.String())
	// we only care about scope-level updates, everything underneath doesn't matter
	if strings.Contains(relKey, "/") {
		return false
	}
	if t == mvccpb.DELETE {
		r.log.Debug("removed scope", zap.String("key", relKey))
		r.scopesM.Lock()
		defer r.scopesM.Unlock()
		delete(r.scopes, relKey)
	} else if t == mvccpb.PUT {
		z, err := r.scopeFromKV(kv)
		if err != nil {
			r.log.Warn("failed to convert scope from event", zap.Error(err))
		} else {
			r.log.Debug("added scope", zap.String("name", z.Name))
			r.scopesM.Lock()
			defer r.scopesM.Unlock()
			r.scopes[z.Name] = z
		}
	}
	return true
}

func (r *Role) loadInitialScopes() {
	scopes, err := r.i.KV().Get(
		r.ctx,
		r.i.KV().Key(
			types.KeyRole,
			types.KeyScopes,
		).Prefix(true).String(),
		clientv3.WithPrefix(),
	)
	if err != nil {
		r.log.Warn("failed to list initial scopes", zap.Error(err))
		time.Sleep(5 * time.Second)
		r.startWatchScopes()
		return
	}
	for _, scope := range scopes.Kvs {
		r.handleScopeOp(mvccpb.PUT, scope)
	}
}

func (r *Role) startWatchScopes() {
	watchChan := r.i.KV().Watch(
		r.ctx,
		r.i.KV().Key(types.KeyRole, types.KeyScopes).Prefix(true).String(),
		clientv3.WithPrefix(),
	)
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			if r.handleScopeOp(event.Type, event.Kv) {
				r.log.Debug("scope watch update", zap.String("key", string(event.Kv.Key)))
			}
		}
	}
}
