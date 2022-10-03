package types

type APIMetricsRecord struct {
	Time    string `json:"time" required:"true"`
	Node    string `json:"node" required:"true"`
	Handler string `json:"handler" required:"true"`
	Value   int64  `json:"value" required:"true"`
}

type APIMetricsGetOutput struct {
	Records []APIMetricsRecord `json:"records" required:"true"`
}
