package api

import (
	"context"
	"net/http"

	"beryju.io/gravity/pkg/roles/api/types"
)

func (r *Role) SessionMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, re *http.Request) {
		session, _ := r.sessions.Get(re, types.SessionName)
		c := context.WithValue(re.Context(), types.RequestSession, session)
		req := re.Clone(c)
		h.ServeHTTP(w, req)
		err := session.Save(req, w)
		if err != nil {
			r.log.WithError(err).Warning("failed to save session")
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
	})
}
