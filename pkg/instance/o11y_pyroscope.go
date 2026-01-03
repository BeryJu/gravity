package instance

import (
	"fmt"
	"net/http"
	"runtime"

	"beryju.io/gravity/pkg/extconfig"
	"github.com/grafana/pyroscope-go"
)

func (i *Instance) startPyroscope() {
	if !extconfig.Get().Observability.Pyroscope.Enabled || extconfig.Get().CI {
		return
	}
	release := fmt.Sprintf("gravity@%s", extconfig.FullVersion())
	runtime.SetMutexProfileFraction(5)
	runtime.SetBlockProfileRate(5)
	pyroscope.Start(pyroscope.Config{
		ApplicationName:   "gravity.beryju.io",
		ServerAddress:     extconfig.Get().Observability.Pyroscope.Server,
		BasicAuthUser:     extconfig.Get().Observability.Pyroscope.Username,
		BasicAuthPassword: extconfig.Get().Observability.Pyroscope.Password,
		Logger:            i.log.Sugar(),
		Tags: map[string]string{
			"gravity.instance": extconfig.Get().Instance.Identifier,
			"gravity.version":  extconfig.Version,
			"gravity.hash":     extconfig.BuildHash,
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

}
