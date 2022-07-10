package etcd

import (
	"encoding/json"
	"fmt"
	"net/http"

	"beryju.io/ddet/pkg/extconfig"
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
	mux.Post("/v0/etcd/join", func(w http.ResponseWriter, r *http.Request) {
		_, err := ro.i.KV().MemberAddAsLearner(r.Context(), []string{r.URL.Query().Get("peer")})
		if err != nil {
			ro.log.WithError(err).Warning("added member")
			return
		}
		env := fmt.Sprintf(
			"ETCD_JOIN_CLUSTER='%s,%s'",
			extconfig.Get().Instance.Identifier,
			extconfig.Get().Instance.IP,
		)
		w.Write([]byte(env))
	})
}
