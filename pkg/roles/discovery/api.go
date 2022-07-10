package discovery

import (
	"net/http"
	"strings"

	"beryju.io/ddet/pkg/roles/discovery/types"
)

func (ro *DiscoveryRole) apiHandlerApply(w http.ResponseWriter, r *http.Request) {
	relKey := r.URL.Query().Get("relKey")
	by := strings.SplitN(relKey, "/", 1)[0]
	if by != types.KeyDevicesByMAC && by != types.KeyDevicesByIP {
		http.Error(w, "invalid key", 400)
		return
	}

	rawDevice, err := ro.i.KV().Get(r.Context(), ro.i.KV().Key(
		types.KeyRole,
		types.KeyDevices,
		relKey,
	))
	if err != nil || len(rawDevice.Kvs) < 1 {
		http.Error(w, "device not found", 404)
		return
	}

	device := ro.deviceFromKV(rawDevice.Kvs[0])
	if by == types.KeyDevicesByIP {
		err = device.toDHCP(r.URL.Query().Get("dhcpScope"))
	} else {
		err = device.toDNS(r.URL.Query().Get("dnsZone"))
	}
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
}
