package instance

import (
	"fmt"

	"beryju.io/gravity/pkg/extconfig"
	"github.com/getsentry/sentry-go"
	"go.uber.org/zap"
)

func (i *Instance) startSentry() {
	if !extconfig.Get().Observability.Sentry.Enabled || extconfig.Get().CI {
		return
	}
	release := fmt.Sprintf("gravity@%s", extconfig.FullVersion())
	rate := 0.5
	if extconfig.Get().Debug {
		rate = 1
	}
	err := sentry.Init(sentry.ClientOptions{
		Dsn:              extconfig.Get().Observability.Sentry.DSN,
		Release:          release,
		EnableTracing:    true,
		TracesSampleRate: rate,
		HTTPTransport:    extconfig.NewUserAgentTransport(release, extconfig.Transport()),
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
