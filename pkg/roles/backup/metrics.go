package backup

import (
	"beryju.io/gravity/pkg/extconfig"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	backupStatus = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gravity_backup_status",
		Help: "Backup status",
		ConstLabels: prometheus.Labels{
			"instance": extconfig.Get().Instance.Identifier,
		},
	}, []string{"status"})
	backupSize = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "gravity_backup_size",
		Help: "Backup size",
		ConstLabels: prometheus.Labels{
			"instance": extconfig.Get().Instance.Identifier,
		},
	})
	backupDuration = promauto.NewGauge(prometheus.GaugeOpts{
		Name: "gravity_backup_duration",
		Help: "Backup duration",
		ConstLabels: prometheus.Labels{
			"instance": extconfig.Get().Instance.Identifier,
		},
	})
)
