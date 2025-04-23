package middleware

import (
	"bufio"
	"errors"
	"net"
	"net/http"
	"time"

	"beryju.io/gravity/pkg/roles/api/auth"
	"beryju.io/gravity/pkg/roles/api/types"
	"github.com/gorilla/sessions"
	"go.uber.org/zap"
)

// responseLogger is wrapper of http.ResponseWriter that keeps track of its HTTP status
// code and body size
type responseLogger struct {
	w      http.ResponseWriter
	status int
	size   int
}

// Header returns the ResponseWriter's Header
func (l *responseLogger) Header() http.Header {
	return l.w.Header()
}

// Support Websocket
func (l *responseLogger) Hijack() (rwc net.Conn, buf *bufio.ReadWriter, err error) {
	if hj, ok := l.w.(http.Hijacker); ok {
		return hj.Hijack()
	}
	return nil, nil, errors.New("http.Hijacker is not available on writer")
}

// Write writes the response using the ResponseWriter
func (l *responseLogger) Write(b []byte) (int, error) {
	if l.status == 0 {
		// The status will be StatusOK if WriteHeader has not been called yet
		l.status = http.StatusOK
	}
	size, err := l.w.Write(b)
	l.size += size
	return size, err
}

// WriteHeader writes the status code for the Response
func (l *responseLogger) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}

// Status returns the response status code
func (l *responseLogger) Status() int {
	return l.status
}

// Size returns the response size
func (l *responseLogger) Size() int {
	return l.size
}

// Flush sends any buffered data to the client
func (l *responseLogger) Flush() {
	if flusher, ok := l.w.(http.Flusher); ok {
		flusher.Flush()
	}
}

// loggingHandler is the http.Handler implementation for LoggingHandler
type loggingHandler struct {
	handler      http.Handler
	logger       *zap.Logger
	afterHandler afterHandler
}

type afterHandler func(l *zap.Logger, r *http.Request) *zap.Logger

// NewLoggingMiddleware provides an http.Handler which logs requests to the HTTP server
func NewLoggingMiddleware(logger *zap.Logger, after afterHandler) func(h http.Handler) http.Handler {
	if after == nil {
		after = func(l *zap.Logger, r *http.Request) *zap.Logger {
			return l
		}
	}
	return func(h http.Handler) http.Handler {
		return loggingHandler{
			handler:      h,
			logger:       logger,
			afterHandler: after,
		}
	}
}

func (h loggingHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	t := time.Now()
	url := *req.URL
	responseLogger := &responseLogger{w: w}
	h.handler.ServeHTTP(responseLogger, req)
	fields := []zap.Field{
		zap.String("host", req.Host),
		zap.String("remote", req.RemoteAddr),
		zap.Duration("runtime", time.Since(t)),
		zap.String("method", req.Method),
		zap.Int("size", responseLogger.Size()),
		zap.Int("status", responseLogger.Status()),
		zap.String("userAgent", req.UserAgent()),
		zap.String("url", url.RequestURI()),
	}
	se := req.Context().Value(types.RequestSession)
	if se != nil {
		session := se.(*sessions.Session)
		u, ok := session.Values[types.SessionKeyUser]
		if ok && u != nil {
			if uu, castOk := u.(auth.User); castOk {
				fields = append(fields, zap.String("user", uu.Username))
			}
		}
	}
	h.afterHandler(h.logger.With(fields...), req).Info("HTTP Request")
}
