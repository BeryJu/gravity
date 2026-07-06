package extconfig

import (
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

func LogPanic(err any, fields ...zap.Field) {
	if err == nil {
		return
	}
	l := Get().Logger()
	fields = append(fields, zap.Stack("stack"))
	if e, ok := err.(error); ok {
		fields = append(fields, zap.Error(e))
		sentry.CaptureException(e)
	} else {
		fields = append(fields, zap.Any("panic", err))
	}
	l.Error("recover", fields...)
}
