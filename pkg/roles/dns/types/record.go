package types

type DNSRecordType string

const (
	DNSRecordTypeA     = "A"
	DNSRecordTypeAAAA  = "AAAA"
	DNSRecordTypePTR   = "PTR"
	DNSRecordTypeSRV   = "SRV"
	DNSRecordTypeMX    = "MX"
	DNSRecordTypeCNAME = "CNAME"
	DNSRecordTypeTXT   = "TXT"
	DNSRecordTypeSOA   = "SOA"
)

func (DNSRecordType) Enum() []interface{} {
	return []interface{}{
		DNSRecordTypeA,
		DNSRecordTypeAAAA,
		DNSRecordTypePTR,
		DNSRecordTypeSRV,
		DNSRecordTypeMX,
		DNSRecordTypeCNAME,
		DNSRecordTypeTXT,
		DNSRecordTypeSOA,
	}
}
