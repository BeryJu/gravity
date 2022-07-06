package dns

import (
	"context"
	"strings"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (r *DNSServerRole) handleZoneOp(t mvccpb.Event_EventType, kv *mvccpb.KeyValue) {
	prefix := r.i.GetKV().Key(KeyRole, KeyZones, "")
	relKey := strings.TrimPrefix(string(kv.Key), prefix)
	// we only care about zone-level updates, everything underneath doesn't matter
	if strings.Contains(relKey, "/") {
		return
	}
	if t == mvccpb.DELETE {
		r.log.WithField("name", r.zones[relKey].Name).Trace("removed zone")
		delete(r.zones, relKey)
	} else if t == mvccpb.PUT {
		z, err := r.zoneFromKV(r.i, kv)
		if err != nil {
			r.log.WithError(err).Warning("failed to convert zone from event")
		} else {
			r.log.WithField("name", z.Name).Trace("added zone")
			r.zones[z.Name] = z
		}
	}
}

func (r *DNSServerRole) startWatchZones() {
	prefix := r.i.GetKV().Key(KeyRole, KeyZones, "")
	zones, err := r.i.GetKV().Get(context.Background(), prefix, clientv3.WithPrefix())
	if err != nil {
		r.log.WithError(err).Warning("failed to list initial zones")
		r.startWatchZones()
		return
	}
	for _, zone := range zones.Kvs {
		r.handleZoneOp(mvccpb.PUT, zone)
	}

	watchChan := r.i.GetKV().Watch(
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
