package dns

import (
	"context"
	"strings"
	"time"

	"beryju.io/gravity/pkg/roles/dns/types"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"go.uber.org/zap"
)

func (r *Role) handleZoneOp(t mvccpb.Event_EventType, kv *mvccpb.KeyValue, ctx context.Context) bool {
	prefix := r.i.KV().Key(types.KeyRole, types.KeyZones).Prefix(true).String()
	relKey := strings.TrimPrefix(string(kv.Key), prefix)
	// we only care about zone-level updates, everything underneath doesn't matter
	if strings.Contains(relKey, "/") {
		return false
	}
	if t == mvccpb.DELETE {
		r.log.Debug("removed zone", zap.String("key", relKey))
		r.zones[relKey].StopWatchingRecords()
		r.zonesM.Lock()
		defer r.zonesM.Unlock()
		delete(r.zones, relKey)
	} else if t == mvccpb.PUT {
		z, err := r.zoneFromKV(kv)
		if err != nil {
			r.log.Warn("failed to convert zone from event", zap.Error(err))
			return true
		}
		z.Init(ctx)
		if oldZone, ok := r.zones[z.Name]; ok {
			oldZone.StopWatchingRecords()
		}
		if !strings.HasSuffix(z.Name, ".") {
			r.log.Warn("Zone is missing trailing preiod, most likely configured incorrectly", zap.String("name", z.Name))
		}
		r.log.Debug("added zone", zap.String("name", z.Name))
		r.zonesM.Lock()
		defer r.zonesM.Unlock()
		r.zones[z.Name] = z
	}
	return true
}

func (r *Role) loadInitialZones(ctx context.Context) {
	zones, err := r.i.KV().Get(
		ctx,
		r.i.KV().Key(
			types.KeyRole,
			types.KeyZones,
		).Prefix(true).String(),
		clientv3.WithPrefix(),
	)
	if err != nil {
		r.log.Warn("failed to list initial zones", zap.Error(err))
		time.Sleep(5 * time.Second)
		r.loadInitialZones(ctx)
		return
	}
	for _, zone := range zones.Kvs {
		r.handleZoneOp(mvccpb.PUT, zone, ctx)
	}
}

func (r *Role) startWatchZones(ctx context.Context) {
	watchChan := r.i.KV().Watch(
		ctx,
		r.i.KV().Key(types.KeyRole, types.KeyZones).Prefix(true).String(),
		clientv3.WithPrefix(),
	)
	for watchResp := range watchChan {
		for _, event := range watchResp.Events {
			if r.handleZoneOp(event.Type, event.Kv, ctx) {
				r.log.Debug("zone watch update", zap.String("key", string(event.Kv.Key)))
			}
		}
	}
}
