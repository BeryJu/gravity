package discovery

import (
	"encoding/json"
	"net/http"

	"beryju.io/gravity/pkg/roles/discovery/types"
	clientv3 "go.etcd.io/etcd/client/v3"
)

func (ro *DiscoveryRole) apiHandlerSubnetList(w http.ResponseWriter, r *http.Request) {
	rawSubnets, err := ro.i.KV().Get(r.Context(), ro.i.KV().Key(
		types.KeyRole,
		types.KeySubnets,
		"",
	), clientv3.WithPrefix())
	if err != nil || len(rawSubnets.Kvs) < 1 {
		http.Error(w, "device not found", 404)
		return
	}
	subnets := make([]*Subnet, len(rawSubnets.Kvs))
	for _, rawSub := range rawSubnets.Kvs {
		sub, err := ro.subnetFromKV(rawSub)
		if err != nil {
			ro.log.WithError(err).Warning("failed to load subnet")
			continue
		}
		subnets = append(subnets, sub)
	}
	json.NewEncoder(w).Encode(subnets)
}
