package discovery

import (
	"fmt"
	"time"

	"beryju.io/ddet/pkg/roles/discovery/types"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (r *DiscoveryRole) startDiscovery(raw *mvccpb.KeyValue) {
	sub, err := r.subnetFromKV(raw)
	if err != nil {
		r.log.WithError(err).Warning("failed to parse subnet")
		return
	}
	go sub.RunDiscovery()
}

func (r *DiscoveryRole) startWatchSubnets() {
	fmt.Println("test")
	prefix := r.i.KV().Key(types.KeyRole, types.KeySubnets, "")
	subnets, err := r.i.KV().Get(r.ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		r.log.WithError(err).Warning("failed to list initial subnets")
		time.Sleep(5 * time.Second)
		r.startWatchSubnets()
		return
	}
	for _, subnet := range subnets.Kvs {
		r.startDiscovery(subnet)
	}

	r.log.Trace(prefix)
	watchChan := r.i.KV().Watch(
		r.ctx,
		prefix,
		clientv3.WithPrefix(),
		clientv3.WithProgressNotify(),
	)
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			if event.Type == clientv3.EventTypeDelete {
				continue
			}
			r.startDiscovery(event.Kv)
		}
	}
}
