package instance

import (
	log "github.com/sirupsen/logrus"
)

type sentryWriter struct {
	entry *log.Entry
}

func NewSentryWriter(log *log.Entry) sentryWriter {
	return sentryWriter{
		entry: log,
	}
}

func (sw sentryWriter) Write(p []byte) (n int, err error) {
	sw.entry.Debug(string(p))
	return len(p), nil
}
