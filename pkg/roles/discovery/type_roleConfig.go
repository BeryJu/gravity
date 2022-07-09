package discovery

import "encoding/json"

type DiscoveryRoleConfig struct {
	Enabled bool `json:"enabled"`
}

func (r *DiscoveryRole) decodeDiscoveryRoleConfig(raw []byte) *DiscoveryRoleConfig {
	def := DiscoveryRoleConfig{
		Enabled: true,
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
