package log_iml

import (
	"net/url"

	"go.uber.org/zap"
)

type inMemoryLogger struct {
	msgs []string
	max  int
}

func (iml *inMemoryLogger) Close() error { return nil }
func (iml *inMemoryLogger) Sync() error  { return nil }
func (iml *inMemoryLogger) Write(log []byte) (int, error) {
	iml.msgs = append(iml.msgs, string(log))
	if len(iml.msgs) > iml.max {
		iml.msgs = iml.msgs[1:]
	}
	return 0, nil
}
func (iml *inMemoryLogger) Messages() []string {
	return iml.msgs
}

var iml *inMemoryLogger

func init() {
	iml = &inMemoryLogger{
		msgs: make([]string, 0),
		max:  300,
	}
	zap.RegisterSink("gravity-in-memory", func(u *url.URL) (zap.Sink, error) {
		return iml, nil
	})
}

type InMemoryLogger interface {
	Messages() []string
}

func Get() InMemoryLogger {
	return iml
}
