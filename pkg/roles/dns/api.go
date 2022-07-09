package dns

import (
	"encoding/json"
	"net/http"

	"beryju.io/ddet/pkg/roles"
	"beryju.io/ddet/pkg/roles/dns/types"
	"github.com/go-chi/chi/v5"
)

func (ro *DNSRole) eventHandlerAPIMux(ev *roles.Event) {
	mux := ev.Payload.Data["mux"].(*chi.Mux)
	mux.Get("/v0/dns/zones", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(ro.zones)
	})
	mux.Get("/v0/dns/zones/{zone}/recods", func(w http.ResponseWriter, r *http.Request) {
		zoneName := chi.URLParam(r, "zone")
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
