package dhcp

import "encoding/json"

type DHCPRoleConfig struct {
	ListenOnly            bool `json:"listenOnly"`
	Port                  int  `json:"port"`
	LeaseNegotiateTimeout int  `json:"leaseNegotiateTimeout"`
}

func (r *DHCPRole) decodeDHCPRoleConfig(raw []byte) *DHCPRoleConfig {
	def := DHCPRoleConfig{
		ListenOnly:            false,
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
