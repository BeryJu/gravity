package dhcp

import "encoding/json"

type RoleConfig struct {
	Port                  int `json:"port"`
	LeaseNegotiateTimeout int `json:"leaseNegotiateTimeout"`
}

func (r *Role) decodeRoleConfig(raw []byte) *RoleConfig {
	def := RoleConfig{
		Port:                  67,
		LeaseNegotiateTimeout: 30,
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
