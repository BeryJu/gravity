package extconfig

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

func init() {
	promauto.NewGaugeFunc(prometheus.GaugeOpts{
		Name: "gravity_instance",
		Help: "The infos about this instance",
		ConstLabels: prometheus.Labels{
			"instance": Get().Instance.Identifier,
			"version":  Version,
			"build":    BuildHash,
		},
	}, func() float64 {
		return 1
	})
}
