package log_iml

import (
	"net/url"
	"sync"

	"go.uber.org/zap"
)

type inMemoryLogger struct {
	msgs []string
	msgM sync.RWMutex
	max  int
}

func (iml *inMemoryLogger) Close() error { return nil }
func (iml *inMemoryLogger) Sync() error  { return nil }
func (iml *inMemoryLogger) Write(log []byte) (int, error) {
	go func() {
		iml.msgM.Lock()
		iml.msgs = append(iml.msgs, string(log))
		if len(iml.msgs) > iml.max {
			iml.msgs = iml.msgs[1:]
		}
		iml.msgM.Unlock()
	}()
	return 0, nil
}
func (iml *inMemoryLogger) Messages() []string {
	iml.msgM.RLock()
	defer iml.msgM.RUnlock()
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
