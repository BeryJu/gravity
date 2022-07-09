package utils

import "strings"

func EnsureTrailingPeriod(name string) string {
	if strings.HasSuffix(name, ".") {
		return name
	}
	return name + "."
}

func EnsureLeadingPeriod(name string) string {
	if strings.HasPrefix(name, ".") {
		return name
	}
	return "." + name
}
