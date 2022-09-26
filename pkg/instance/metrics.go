package instance

import (
	"beryju.io/gravity/pkg/extconfig"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	instanceRoles = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gravity_instance_roles",
		Help: "The configured roles for this instance",
		ConstLabels: prometheus.Labels{
			"instance": extconfig.Get().Instance.Identifier,
		},
	}, []string{"roleId"})
	instanceRoleStarted = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gravity_instance_roles_started",
		Help: "The time when a role was last started",
		ConstLabels: prometheus.Labels{
			"instance": extconfig.Get().Instance.Identifier,
		},
	}, []string{"roleId"})
)
