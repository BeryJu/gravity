package middleware

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
				var er error
				if r.Header.Get("Accept") == "application/json" {
					_, er = w.Write([]byte("{\"error\": \"internal error\"}"))
				} else {
					_, er = w.Write([]byte("internal error"))
				}
				if er != nil {
					l.Warn("failed to write error message", zap.Error(er))
				}
			}()
			h.ServeHTTP(w, r)
		})
	}
}
