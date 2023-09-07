package log_iml

import (
	"bytes"
	"net/url"
	"strings"
	"sync"

	"go.uber.org/zap"
)

type inMemoryLogger struct {
	msgs    []string
	msgM    sync.RWMutex
	lineBuf *bytes.Buffer
	max     int
}

func (iml *inMemoryLogger) Close() error { return nil }
func (iml *inMemoryLogger) Sync() error  { return nil }
func (iml *inMemoryLogger) flush() {
	go func() {
		iml.msgM.Lock()
		iml.msgs = append(iml.msgs, iml.lineBuf.String())
		iml.lineBuf.Reset()
		if len(iml.msgs) > iml.max {
			iml.msgs = iml.msgs[1:]
		}
		iml.msgM.Unlock()
	}()
}
func (iml *inMemoryLogger) Write(log []byte) (int, error) {
	n, err := iml.lineBuf.Write(log)
	if strings.Contains(string(log), "\n") {
		iml.flush()
	}
	return n, err
}
func (iml *inMemoryLogger) Messages() []string {
	iml.msgM.RLock()
	defer iml.msgM.RUnlock()
	return iml.msgs
}

var iml *inMemoryLogger

func init() {
	iml = &inMemoryLogger{
		msgs:    make([]string, 0),
		max:     300,
		lineBuf: new(bytes.Buffer),
	}
	err := zap.RegisterSink("gravity-in-memory", func(u *url.URL) (zap.Sink, error) {
		return iml, nil
	})
	if err != nil {
		panic(err)
	}
}

type InMemoryLogger interface {
	Messages() []string
}

func Get() InMemoryLogger {
	return iml
}
