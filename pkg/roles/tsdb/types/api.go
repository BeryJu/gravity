package types

import (
	"time"

	"beryju.io/gravity/pkg/roles/dns/types"
)

type APIMetricsRole string

func (APIMetricsRole) Enum() []interface{} {
	return []interface{}{
		KeySystem,
		// TODO: Do this more dynamically?
		types.KeyRole,
	}
}

type APIMetricsGetInput struct {
	Role     APIMetricsRole `query:"role" required:"true"`
	Category *string        `query:"category"`
	Since    *time.Time     `query:"since" description:"Optionally set a start time for which to return datapoints after"`
}

type APIMetricsRecord struct {
	Time  time.Time `json:"time" required:"true"`
	Keys  []string  `json:"keys" required:"true"`
	Node  string    `json:"node" required:"true"`
	Value int64     `json:"value" required:"true"`
}

type APIMetricsGetOutput struct {
	Records []APIMetricsRecord `json:"records" required:"true"`
}
