package etcd

import (
	"encoding/json"
	"net/http"

	"beryju.io/ddet/pkg/roles"
	"github.com/go-chi/chi/v5"
)

func (ro *EmbeddedEtcd) eventHandlerAPIMux(ev *roles.Event) {
	mux := ev.Payload.Data["mux"].(*chi.Mux)
	mux.Get("/v0/etcd/members", func(w http.ResponseWriter, r *http.Request) {
		members, err := ro.i.KV().MemberList(r.Context())
		if err != nil {
			ro.log.WithError(err).Warning("failed to list members")
			return
		}
		json.NewEncoder(w).Encode(members.Members)
	})
}
