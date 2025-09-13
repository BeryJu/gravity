package externaldns

import (
	"context"
	"slices"
	"strings"

	"beryju.io/gravity/api"
	"beryju.io/gravity/pkg/roles/dns/types"
	"beryju.io/gravity/pkg/roles/dns/utils"
)

func (s *Server) prefetchZones(ctx context.Context) error {
	zones, hr, err := s.api.RolesDnsAPI.DnsGetZones(ctx).Execute()
	if err != nil {
		return s.apiError(hr, err)
	}
	fz := []api.DnsAPIZone{}
	for _, zone := range zones.Zones {
		if len(Get().DomainFilter) > 0 && !slices.Contains(Get().DomainFilter, zone.Name) {
			continue
		}
		fz = append(fz, zone)
	}
	s.zones = fz
	return nil
}

func (s *Server) findZoneForRecord(fqdn string) (*api.DnsAPIZone, string) {
	var longestMatch *api.DnsAPIZone
	longestMatchLength := 0
	for _, zone := range s.zones {
		if !strings.HasSuffix(utils.EnsureTrailingPeriod(fqdn), zone.Name) {
			continue
		}
		if len(zone.Name) <= longestMatchLength {
			continue
		}
		longestMatchLength = len(zone.Name)
		longestMatch = &zone
	}
	if longestMatch != nil {
		hostname := strings.TrimSuffix(strings.Replace(fqdn, strings.TrimSuffix(longestMatch.Name, types.DNSSep), "", 1), types.DNSSep)
		return longestMatch, hostname
	}
	return longestMatch, ""
}
