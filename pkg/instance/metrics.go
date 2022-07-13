package instance

import (
	"beryju.io/gravity/pkg/extconfig"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	instanceRoles = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gravity_instance_roles",
		Help: "The active roles for this instance",
		ConstLabels: prometheus.Labels{
			"instance": extconfig.Get().Instance.Identifier,
		},
	}, []string{"roleId"})
)
