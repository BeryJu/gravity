package types

type APIMetricsRecord struct {
	Node    string `json:"node" required:"true"`
	Handler string `json:"handler" required:"true"`
	Value   int64  `json:"value" required:"true"`
}
