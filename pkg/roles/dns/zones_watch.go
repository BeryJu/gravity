package dns

import (
	"context"
	"strings"
	"time"

	"beryju.io/ddet/pkg/roles/dns/types"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (r *DNSRole) handleZoneOp(t mvccpb.Event_EventType, kv *mvccpb.KeyValue) {
	prefix := r.i.KV().Key(types.KeyRole, types.KeyZones, "")
	relKey := strings.TrimPrefix(string(kv.Key), prefix)
	// we only care about zone-level updates, everything underneath doesn't matter
	if strings.Contains(relKey, "/") {
		return
	}
	if t == mvccpb.DELETE {
		r.log.WithField("name", r.zones[relKey].Name).Trace("removed zone")
		delete(r.zones, relKey)
	} else if t == mvccpb.PUT {
		z, err := r.zoneFromKV(kv)
		if err != nil {
			r.log.WithError(err).Warning("failed to convert zone from event")
		} else {
			r.log.WithField("name", z.Name).Debug("added zone")
			r.zones[z.Name] = z
		}
	}
}

func (r *DNSRole) startWatchZones() {
	prefix := r.i.KV().Key(types.KeyRole, types.KeyZones, "")
	zones, err := r.i.KV().Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		r.log.WithError(err).Warning("failed to list initial zones")
		time.Sleep(5 * time.Second)
		r.startWatchZones()
		return
	}
	for _, zone := range zones.Kvs {
		r.handleZoneOp(mvccpb.PUT, zone)
	}

	watchChan := r.i.KV().Watch(
		context.Background(),
		prefix,
		clientv3.WithPrefix(),
		clientv3.WithProgressNotify(),
	)
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			r.handleZoneOp(event.Type, event.Kv)
			r.log.WithField("key", string(event.Kv.Key)).Trace("zone watch update")
		}
	}
}
