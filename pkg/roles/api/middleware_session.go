package api

import (
	"context"
	"net/http"

	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

type sessionWriter struct {
	w       http.ResponseWriter
	session *sessions.Session
	req     *http.Request
	log     *zap.Logger
}

func (sw *sessionWriter) WriteHeader(statusCode int) {
	if dirty, ok := sw.session.Values[types.SessionKeyDirty]; ok && dirty == true {
		sw.log.Debug("session is dirty, writing")
		sw.session.Values[types.SessionKeyDirty] = false
		err := sw.session.Save(sw.req, sw.w)
		if err != nil {
			sw.log.Warn("failed to save session", zap.Error(err))
		}
	}
	sw.w.WriteHeader(statusCode)
}

func (sw *sessionWriter) Header() http.Header {
	return sw.w.Header()
}

func (sw *sessionWriter) Write(data []byte) (int, error) {
	return sw.w.Write(data)
}

func (r *Role) SessionMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, re *http.Request) {
		session, _ := r.sessions.Get(re, types.SessionName)
		c := context.WithValue(re.Context(), types.RequestSession, session)
		req := re.Clone(c)

		sw := &sessionWriter{
			w:       w,
			session: session,
			req:     req,
			log:     r.log,
		}
		h.ServeHTTP(sw, req)
	})
}
