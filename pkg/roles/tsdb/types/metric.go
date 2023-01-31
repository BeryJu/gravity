package types

type Metric struct {
	Value        int  `json:"value"`
	ResetOnWrite bool `json:"resetOnWrite"`
}
