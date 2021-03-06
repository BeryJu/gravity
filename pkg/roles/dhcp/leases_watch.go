package dhcp

import (
	"time"

	"beryju.io/gravity/pkg/roles/dhcp/types"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (r *Role) handleLeaseOp(ev *clientv3.Event) {
	rec, err := r.leaseFromKV(ev.Kv)
	if ev.Type == clientv3.EventTypeDelete {
		delete(r.leases, rec.Identifier)
	} else {
		// Check if the lease parsed above actually was parsed correctly,
		// we don't care for that when removing, but prevent adding
		// empty leases
		if err != nil {
			r.log.WithError(err).Warning("failed to parse lease")
			return
		}
		r.leases[rec.Identifier] = rec
	}
}

func (r *Role) loadInitialLeases() {
	leases, err := r.i.KV().Get(r.ctx, r.i.KV().Key(types.KeyRole, types.KeyLeases).Prefix(true).String(), clientv3.WithPrefix())
	if err != nil {
		r.log.WithError(err).Warning("failed to list initial leases")
		time.Sleep(5 * time.Second)
		r.startWatchLeases()
		return
	}
	for _, lease := range leases.Kvs {
		r.handleLeaseOp(&clientv3.Event{
			Type: mvccpb.PUT,
			Kv:   lease,
		})
	}
}

func (r *Role) startWatchLeases() {
	watchChan := r.i.KV().Watch(
		r.ctx,
		r.i.KV().Key(types.KeyRole, types.KeyLeases).Prefix(true).String(),
		clientv3.WithPrefix(),
		clientv3.WithProgressNotify(),
	)
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			r.handleLeaseOp(event)
		}
	}
}
