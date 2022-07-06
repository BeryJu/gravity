package dhcp

import (
	"context"
	"strings"

	"beryju.io/ddet/pkg/roles/dhcp/types"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (r *DHCPRole) handleScopeOp(t mvccpb.Event_EventType, kv *mvccpb.KeyValue) {
	prefix := r.i.GetKV().Key(types.KeyRole, types.KeyScopes, "")
	relKey := strings.TrimPrefix(string(kv.Key), prefix)
	// we only care about scope-level updates, everything underneath doesn't matter
	if strings.Contains(relKey, "/") {
		return
	}
	if t == mvccpb.DELETE {
		r.log.WithField("name", r.scopes[relKey].Name).Trace("removed scope")
		delete(r.scopes, relKey)
	} else if t == mvccpb.PUT {
		z, err := r.scopeFromKV(kv)
		if err != nil {
			r.log.WithError(err).Warning("failed to convert scope from event")
		} else {
			r.log.WithField("name", z.Name).Debug("added scope")
			r.scopes[z.Name] = z
		}
	}
}

func (r *DHCPRole) startWatchScopes() {
	prefix := r.i.GetKV().Key(types.KeyRole, types.KeyScopes, "")
	scopes, err := r.i.GetKV().Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		r.log.WithError(err).Warning("failed to list initial scopes")
		r.startWatchScopes()
		return
	}
	for _, scope := range scopes.Kvs {
		r.handleScopeOp(mvccpb.PUT, scope)
	}

	watchChan := r.i.GetKV().Watch(
		context.Background(),
		prefix,
		clientv3.WithPrefix(),
		clientv3.WithProgressNotify(),
	)
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			r.handleScopeOp(event.Type, event.Kv)
			r.log.WithField("key", string(event.Kv.Key)).Trace("scope watch update")
		}
	}
}
