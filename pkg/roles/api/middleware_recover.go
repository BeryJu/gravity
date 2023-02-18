package api

import (
	"net/http"

	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

func NewRecoverMiddleware(l *zap.Logger) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				err := recover()
				if err == nil {
					return
				}
				if e, ok := err.(error); ok {
					l.Error("recover in API handler", zap.Error(e))
					sentry.CaptureException(e)
				} else {
					l.Error("recover in API Handler", zap.Any("panic", err))
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
