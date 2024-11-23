package log_iml

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type inMemoryLogger struct {
	msgs []zapcore.Entry
	msgM sync.RWMutex
	max  int
}

func (iml *inMemoryLogger) Hook() zap.Option {
	return zap.Hooks(func(e zapcore.Entry) error {
		iml.msgM.Lock()
		defer iml.msgM.Unlock()
		iml.msgs = append(iml.msgs, e)
		if len(iml.msgs) > iml.max {
			iml.msgs = iml.msgs[1:]
		}
		return nil
	})
}

func (iml *inMemoryLogger) Messages() []zapcore.Entry {
	iml.msgM.RLock()
	defer iml.msgM.RUnlock()
	return iml.msgs
}

var iml *inMemoryLogger

func init() {
	iml = &inMemoryLogger{
		msgs: make([]zapcore.Entry, 0),
		max:  300,
	}
}

type InMemoryLogger interface {
	Messages() []zapcore.Entry
	Hook() zap.Option
}

func Get() InMemoryLogger {
	return iml
}
