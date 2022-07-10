package dns

import (
	"encoding/json"
	"net/http"

	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/dns/types"
	"github.com/gorilla/mux"
)

func (ro *DNSRole) eventHandlerAPIMux(ev *roles.Event) {
	m := ev.Payload.Data["mux"].(*mux.Router)
	m.Path("/v0/dns/zones").Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ro.zones)
	})
	m.Path("/v0/dns/zones/{zone}/recods").Methods(http.MethodGet).HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		zoneName := vars["zone"]
		zone, ok := ro.zones[zoneName]
		if !ok {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		rawRecords, err := ro.i.KV().Get(r.Context(), ro.i.KV().Key(
			types.KeyRole,
			types.KeyZones,
			zone.Name,
			"",
		))
		if err != nil {
			http.Error(w, "failed to get records", http.StatusInternalServerError)
			return
		}
		records := make([]Record, len(rawRecords.Kvs))
		for idx, rec := range rawRecords.Kvs {
			records[idx] = *zone.recordFromKV(rec)
		}
		json.NewEncoder(w).Encode(records)
	})
}
