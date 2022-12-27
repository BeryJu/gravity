package instance

import (
	"fmt"

	"beryju.io/gravity/pkg/extconfig"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

func (i *Instance) startSentry() {
	transport := sentry.NewHTTPTransport()
	transport.Configure(sentry.ClientOptions{
		HTTPTransport: extconfig.Transport(),
	})
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              "https://ccd520a9a2b8458ca1e82108a8afb801@sentry.beryju.org/17",
		Environment:      "",
		Release:          fmt.Sprintf("gravity@%s", extconfig.FullVersion()),
		EnableTracing:    true,
		TracesSampleRate: 0.5,
		Transport:        transport,
		Debug:            extconfig.Get().Debug,
		DebugWriter:      NewSentryWriter(i.log.Named("sentry")),
	})
	if err != nil {
		i.log.Warn("failed to init sentry", zap.Error(err))
		return
	}
	sentry.ConfigureScope(func(scope *sentry.Scope) {
		scope.SetTag("gravity.instance", extconfig.Get().Instance.Identifier)
		scope.SetTag("gravity.version", extconfig.Version)
		scope.SetTag("gravity.hash", extconfig.BuildHash)
	})
}

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
