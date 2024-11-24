package log_iml

import (
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type InMemoryLogger struct {
	msgs []zapcore.Entry
	msgM sync.RWMutex
	max  int
}

func (iml *InMemoryLogger) Hook() zap.Option {
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

func (iml *InMemoryLogger) MaxSize() int {
	return iml.max
}

func (iml *InMemoryLogger) Flush() {
	iml.msgM.Lock()
	iml.msgs = make([]zapcore.Entry, 0)
	iml.msgM.Unlock()
}

func (iml *InMemoryLogger) Messages() []zapcore.Entry {
	iml.msgM.RLock()
	defer iml.msgM.RUnlock()
	return iml.msgs
}

var iml *InMemoryLogger

func init() {
	iml = &InMemoryLogger{
		msgs: make([]zapcore.Entry, 0),
		max:  300,
	}
}

func Get() *InMemoryLogger {
	return iml
}
