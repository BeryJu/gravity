package log_iml

import (
	"go.uber.org/zap/zapcore"
)

type ZapCore struct {
	zapcore.Core
	fields []zapcore.Field
}

func (h *ZapCore) Check(entry zapcore.Entry, ce *zapcore.CheckedEntry) *zapcore.CheckedEntry {
	if h.Enabled(entry.Level) {
		return ce.AddCore(entry, h)
	}
	return ce
}

func (h *ZapCore) Write(entry zapcore.Entry, fields []zapcore.Field) error {
	Get().Write(entry, h.fields, fields)
	return h.Core.Write(entry, fields)
}

func (h *ZapCore) With(fields []zapcore.Field) zapcore.Core {
	return &ZapCore{
		Core:   h.Core.With(fields),
		fields: fields,
	}
}
