package instance

import (
	"go.uber.org/zap"
)

type sentryWriter struct {
	logger *zap.Logger
}

func NewSentryWriter(log *zap.Logger) sentryWriter {
	return sentryWriter{
		logger: log,
	}
}

func (sw sentryWriter) Write(p []byte) (n int, err error) {
	sw.logger.Debug(string(p))
	return len(p), nil
}
