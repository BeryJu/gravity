package dns

import (
	"beryju.io/gravity/pkg/extconfig"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	dnsQueries = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gravity_dns_queries",
		Help: "DNS queries",
		ConstLabels: prometheus.Labels{
			"instance": extconfig.Get().Instance.Identifier,
		},
	}, []string{"queryType", "handler", "zone"})
	dnsQueryDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name: "gravity_dns_query_duration",
		Help: "DNS queries duration",
		ConstLabels: prometheus.Labels{
			"instance": extconfig.Get().Instance.Identifier,
		},
	}, []string{"queryType", "handler", "zone"})
)

var (
	dnsRecordsMetric = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gravity_dns_records",
		Help: "DNS records",
		ConstLabels: prometheus.Labels{
			"instance": extconfig.Get().Instance.Identifier,
		},
	}, []string{"zone"})
)
