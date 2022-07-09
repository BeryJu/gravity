package types

const (
	EventTopicDHCPLeaseGiven = "roles.dhcp.lease.given"
)

type EventLeaseGiven struct {
	Hostname string
	Address  string
}
