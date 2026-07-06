package types

import "time"

const (
	KeyRole        = "dns"
	KeyZones       = "zones"
	KeyHandlerType = "handler"
)

const (
	DNSWildcard = "*"
	// Special name for DNS records at the zone apex
	DNSRootRecord = "@"
	// Special name for the root zone in gravity
	DNSRootZone = "."
	// Separator between DNS labels
	DNSSep = "."
	// Separator between multiple TXT values
	TXTSeparator = "\n"
)

const (
	DefaultUpstreamTimeout = time.Second * 2
)
