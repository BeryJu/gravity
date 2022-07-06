package dns

import "encoding/json"

type DNSRoleConfig struct {
	Port int32 `json:"port"`
}

func (r *DNSServerRole) decodeDNSRoleConfig(raw []byte) *DNSRoleConfig {
	def := DNSRoleConfig{
		Port: 53,
	}
	if len(raw) < 1 {
		return &def
	}
	err := json.Unmarshal(raw, &def)
	if err != nil {
		r.log.WithError(err).Warning("failed to decode role config")
	}
	return &def
}
