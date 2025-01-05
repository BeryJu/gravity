package utils

import (
	"strings"

	"beryju.io/gravity/pkg/roles/dns/types"
)

func EnsureTrailingPeriod(name string) string {
	if strings.HasSuffix(name, types.DNSSep) {
		return name
	}
	return name + types.DNSSep
}

func EnsureLeadingPeriod(name string) string {
	if strings.HasPrefix(name, types.DNSSep) {
		return name
	}
	return types.DNSSep + name
}
