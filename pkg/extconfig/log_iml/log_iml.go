package log_iml

import (
	"fmt"
	"sync"

	"go.uber.org/zap/zapcore"
)

type Message struct {
	Entry        zapcore.Entry
	Fields       []zapcore.Field
	LoggerFields []zapcore.Field
}

type InMemoryLogger struct {
	msgs []Message
	msgM sync.RWMutex
	max  int
}

func (iml *InMemoryLogger) Write(entry zapcore.Entry, logFields []zapcore.Field, fields []zapcore.Field) {
	iml.msgM.Lock()
	defer iml.msgM.Unlock()
	iml.msgs = append(iml.msgs, Message{entry, fields, logFields})
	if len(iml.msgs) > iml.max {
		iml.msgs = iml.msgs[1:]
	}
}

func (iml *InMemoryLogger) MaxSize() int {
	return iml.max
}

func (iml *InMemoryLogger) Flush() {
	iml.msgM.Lock()
	iml.msgs = make([]Message, 0)
	iml.msgM.Unlock()
}

func (iml *InMemoryLogger) Messages() []Message {
	iml.msgM.RLock()
	defer iml.msgM.RUnlock()
	return iml.msgs
}

var iml *InMemoryLogger

func init() {
	iml = &InMemoryLogger{
		msgs: make([]Message, 0),
		max:  300,
	}
}

func Get() *InMemoryLogger {
	return iml
}

func (msg Message) FieldsToMap() map[string]interface{} {
	result := make(map[string]interface{})

	convertSingleField := func(field zapcore.Field) {
		switch field.Type {
		case zapcore.BoolType:
			result[field.Key] = field.Integer == 1
		case zapcore.Int8Type, zapcore.Int16Type, zapcore.Int32Type, zapcore.Int64Type:
			result[field.Key] = field.Integer
		case zapcore.Uint8Type, zapcore.Uint16Type, zapcore.Uint32Type, zapcore.Uint64Type:
			result[field.Key] = uint64(field.Integer)
		case zapcore.Float32Type:
			result[field.Key] = float32(field.Integer)
		case zapcore.Float64Type:
			result[field.Key] = float64(field.Integer)
		case zapcore.StringType:
			result[field.Key] = field.String
		case zapcore.TimeType:
			if field.Interface != nil {
				result[field.Key] = field.Interface
			}
		case zapcore.DurationType:
			result[field.Key] = field.Integer // Duration as nanoseconds
		case zapcore.ErrorType:
			if field.Interface != nil {
				result[field.Key] = field.Interface.(error).Error()
			}
		case zapcore.StringerType:
			if field.Interface != nil {
				result[field.Key] = field.Interface.(fmt.Stringer).String()
			}
		default:
			result[field.Key] = fmt.Sprintf("%+v", field.Interface)
		}
	}

	for _, field := range msg.LoggerFields {
		convertSingleField(field)
	}
	for _, field := range msg.Fields {
		convertSingleField(field)
	}
	return result
}
