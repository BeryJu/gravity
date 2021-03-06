package dns

import (
	"context"
	"strings"
	"time"

	"beryju.io/gravity/pkg/roles/dns/types"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (r *Role) handleZoneOp(t mvccpb.Event_EventType, kv *mvccpb.KeyValue) bool {
	prefix := r.i.KV().Key(types.KeyRole, types.KeyZones).Prefix(true).String()
	relKey := strings.TrimPrefix(string(kv.Key), prefix)
	// we only care about zone-level updates, everything underneath doesn't matter
	if strings.Contains(relKey, "/") {
		return false
	}
	if t == mvccpb.DELETE {
		r.log.WithField("name", r.zones[relKey].Name).Trace("removed zone")
		r.zones[relKey].StopWatchingRecords()
		delete(r.zones, relKey)
	} else if t == mvccpb.PUT {
		z, err := r.zoneFromKV(kv)
		if err != nil {
			r.log.WithError(err).Warning("failed to convert zone from event")
		} else {
			if oldZone, ok := r.zones[z.Name]; ok {
				oldZone.StopWatchingRecords()
			}
			if !strings.HasSuffix(z.Name, ".") {
				r.log.WithField("name", z.Name).Warning("Zone is missing trailing preiod, most likely configured incorrectly")
			}
			r.log.WithField("name", z.Name).Debug("added zone")
			r.zones[z.Name] = z
		}
	}
	return true
}

func (r *Role) loadInitialZones() {
	zones, err := r.i.KV().Get(context.Background(), r.i.KV().Key(types.KeyRole, types.KeyZones).Prefix(true).String(), clientv3.WithPrefix())
	if err != nil {
		r.log.WithError(err).Warning("failed to list initial zones")
		time.Sleep(5 * time.Second)
		r.startWatchZones()
		return
	}
	for _, zone := range zones.Kvs {
		r.handleZoneOp(mvccpb.PUT, zone)
	}
}

func (r *Role) startWatchZones() {
	watchChan := r.i.KV().Watch(
		r.ctx,
		r.i.KV().Key(types.KeyRole, types.KeyZones).Prefix(true).String(),
		clientv3.WithPrefix(),
		clientv3.WithProgressNotify(),
	)
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			if r.handleZoneOp(event.Type, event.Kv) {
				r.log.WithField("key", string(event.Kv.Key)).Trace("zone watch update")
			}
		}
	}
}
