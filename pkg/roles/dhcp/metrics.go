package dhcp

import (
	"math/big"

	"beryju.io/gravity/pkg/extconfig"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	dhcpRequests = promauto.NewCounterVec(prometheus.CounterOpts{
		Name: "gravity_dhcp_requests",
		Help: "DHCP Requests",
		ConstLabels: prometheus.Labels{
			"instance": extconfig.Get().Instance.Identifier,
		},
	}, []string{"type", "scope"})
	dhcpScopeSize = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gravity_dhcp_scope_size",
		Help: "Total free IP addresses in a scope",
		ConstLabels: prometheus.Labels{
			"instance": extconfig.Get().Instance.Identifier,
		},
	}, []string{"scope"})
	dhcpScopeUsage = promauto.NewGaugeVec(prometheus.GaugeOpts{
		Name: "gravity_dhcp_scope_usage",
		Help: "Used IP addresses in a scope",
		ConstLabels: prometheus.Labels{
			"instance": extconfig.Get().Instance.Identifier,
		},
	}, []string{"scope"})
)

// calculateUsage Calculate scope usage for prometheus metrics
func (s *Scope) calculateUsage() {
	usable := s.ipam.UsableSize()
	dhcpScopeSize.WithLabelValues(s.Name).Set(float64(usable.Uint64()))
	used := big.NewInt(0)
	for _, lease := range s.role.leases.Iter() {
		if lease.ScopeKey != s.Name {
			continue
		}
		used = used.Add(used, big.NewInt(1))
	}
	dhcpScopeUsage.WithLabelValues(s.Name).Set(float64(used.Uint64()))
}
