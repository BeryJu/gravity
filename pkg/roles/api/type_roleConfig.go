package api

import "encoding/json"

type APIRoleConfig struct {
	Port int32 `json:"port"`
}

func (r *APIRole) decodeAPIRoleConfig(raw []byte) *APIRoleConfig {
	def := APIRoleConfig{
		Port: 8008,
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
