package watcher

import (
	"beryju.io/gravity/pkg/extconfig"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	watcherEvents = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gravity_watcher_events",
		Help: "Watch events",
		ConstLabels: prometheus.Labels{
			"instance": extconfig.Get().Instance.Identifier,
		},
	}, []string{"key_prefix", "type"})
)
