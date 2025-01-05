package types

import "time"

const (
	KeyRole  = "dns"
	KeyZones = "zones"
)

const (
	DNSRecordTypeA     = "A"
	DNSRecordTypeAAAA  = "AAAA"
	DNSRecordTypePTR   = "PTR"
	DNSRecordTypeCNAME = "CNAME"
)

const (
	DNSWildcard   = "*"
	DNSRootRecord = "@"
	DNSRootZone   = "."
)

const (
	DefaultUpstreamTimeout = time.Second * 2
)
