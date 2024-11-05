package dhcp

import (
	"context"
	"errors"
	"strings"
	"time"

	"beryju.io/gravity/pkg/roles/dhcp/types"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

func (r *Role) handleLeaseOp(ev *clientv3.Event) {
	rec, err := r.leaseFromKV(ev.Kv)
	if ev.Type == clientv3.EventTypeDelete {
		r.leasesM.Lock()
		defer r.leasesM.Unlock()
		delete(r.leases, rec.Identifier)
	} else {
		// Check if the lease parsed above actually was parsed correctly,
		// we don't care for that when removing, but prevent adding
		// empty leases
		if err != nil {
			r.log.Warn("failed to parse lease", zap.Error(err))
			return
		}
		r.leasesM.Lock()
		defer r.leasesM.Unlock()
		r.leases[rec.Identifier] = rec
	}
}

func (r *Role) loadInitialLeases(ctx context.Context) {
	prefix := r.i.KV().Key(
		types.KeyRole,
		types.KeyScopes,
	).Prefix(true).String()
	leases, err := r.i.KV().Get(
		ctx,
		prefix,
		clientv3.WithPrefix(),
	)
	if err != nil {
		r.log.Warn("failed to list initial leases", zap.Error(err))
		if !errors.Is(err, context.Canceled) {
			time.Sleep(5 * time.Second)
			r.loadInitialLeases(ctx)
		}
		return
	}
	for _, lease := range leases.Kvs {
		relKey := strings.ReplaceAll(string(lease.Key), prefix, "")
		if !strings.Contains("/", relKey) {
			continue
		}
		r.handleLeaseOp(&clientv3.Event{
			Type: mvccpb.PUT,
			Kv:   lease,
		})
	}
}

func (r *Role) startWatchLeases() {
	prefix := r.i.KV().Key(
		types.KeyRole,
		types.KeyScopes,
	).Prefix(true).String()
	watchChan := r.i.KV().Watch(
		r.ctx,
		prefix,
		clientv3.WithPrefix(),
	)
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			relKey := strings.ReplaceAll(string(event.Kv.Key), prefix, "")
			if !strings.Contains("/", relKey) {
				continue
			}
			r.handleLeaseOp(event)
		}
	}
}
