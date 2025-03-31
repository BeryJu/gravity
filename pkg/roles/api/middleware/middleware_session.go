package middleware

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

// sessionHandler is the http.Handler implementation for sessionHandler
type sessionHandler struct {
	handler http.Handler
	store   sessions.Store
	log     *zap.Logger
}

func NewSessionMiddleware(store sessions.Store, log *zap.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return &sessionHandler{
			handler: h,
			store:   store,
			log:     log,
		}
	}
}

func (sh *sessionHandler) ServeHTTP(w http.ResponseWriter, re *http.Request) {
	session, _ := sh.store.Get(re, types.SessionName)
	c := re.Context()
	if c.Value(types.RequestSession) == nil {
		c = context.WithValue(c, types.RequestSession, session)
	}
	req := re.WithContext(c)

	sw := &sessionWriter{
		w:       w,
		session: session,
		req:     req,
		log:     sh.log,
	}
	sh.handler.ServeHTTP(sw, req)
}
