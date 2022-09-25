package api

import (
	"net/http"

	"beryju.io/gravity/pkg/extconfig"
	"github.com/getsentry/sentry-go"
	log "github.com/sirupsen/logrus"
)

func NewRecoverMiddleware(l *log.Entry) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				err := extconfig.Recover()
				if err == nil {
					return
				}
				if e, ok := err.(error); ok {
					l.WithError(e).Warning("recover in API handler")
					sentry.CaptureException(e)
				} else {
					l.WithField("panic", err).Warning("recover in API handler")
				}
				w.WriteHeader(http.StatusInternalServerError)
				if r.Header.Get("Accept") == "application/json" {
					w.Write([]byte("{\"error\": \"internal error\"}"))
				} else {
					w.Write([]byte("internal error"))
				}
			}()
			h.ServeHTTP(w, r)
		})
	}
}
