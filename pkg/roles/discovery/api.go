package discovery

import (
	"net/http"

	"beryju.io/ddet/pkg/roles/discovery/types"
)

func (ro *DiscoveryRole) apiHandlerApply(w http.ResponseWriter, r *http.Request) {
	relKey := r.URL.Query().Get("relKey")

	rawDevice, err := ro.i.KV().Get(r.Context(), ro.i.KV().Key(
		types.KeyRole,
		types.KeyDevices,
		relKey,
	))
	if err != nil || len(rawDevice.Kvs) < 1 {
		http.Error(w, "device not found", 404)
		return
	}

}
