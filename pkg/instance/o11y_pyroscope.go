package instance

import (
	"fmt"
	"net/http"
	"runtime"

	"beryju.io/gravity/pkg/extconfig"
	"github.com/grafana/pyroscope-go"
	"go.uber.org/zap"
)

var profiler *pyroscope.Profiler

func (i *Instance) startPyroscope() {
	if !extconfig.Get().Observability.Pyroscope.Enabled || extconfig.Get().CI {
		return
	}
	release := fmt.Sprintf("gravity@%s", extconfig.FullVersion())
	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)
	pr, err := pyroscope.Start(pyroscope.Config{
		ApplicationName:   "gravity.beryju.io",
		ServerAddress:     extconfig.Get().Observability.Pyroscope.Server,
		BasicAuthUser:     extconfig.Get().Observability.Pyroscope.Username,
		BasicAuthPassword: extconfig.Get().Observability.Pyroscope.Password,
		Logger:            i.log.Named("o11y.pyroscope").Sugar(),
		Tags: map[string]string{
			"gravity_instance":   extconfig.Get().Instance.Identifier,
			"gravity_version":    extconfig.Version,
			"gravity_hash":       extconfig.BuildHash,
			"service_repository": "https://github.com/BeryJu/gravity",
			"service_git_ref":    extconfig.BuildHash,
		},
		DisableGCRuns: true,
		HTTPClient: &http.Client{
			Transport: extconfig.NewUserAgentTransport(release, extconfig.Transport()),
		},
		ProfileTypes: []pyroscope.ProfileType{
			pyroscope.ProfileCPU,
			pyroscope.ProfileAllocObjects,
			pyroscope.ProfileAllocSpace,
			pyroscope.ProfileInuseObjects,
			pyroscope.ProfileInuseSpace,
			pyroscope.ProfileGoroutines,
			pyroscope.ProfileMutexCount,
			pyroscope.ProfileMutexDuration,
			pyroscope.ProfileBlockCount,
			pyroscope.ProfileBlockDuration,
		},
	})
	if err != nil {
		i.log.Warn("failed to init pyroscope", zap.Error(err))
		return
	}
	profiler = pr
}

func (i *Instance) stopPyroscope() {
	if profiler == nil {
		return
	}
	profiler.Flush(true)
	err := profiler.Stop()
	if err != nil {
		i.log.Warn("failed to flush pyroscope", zap.Error(err))
	}
}
