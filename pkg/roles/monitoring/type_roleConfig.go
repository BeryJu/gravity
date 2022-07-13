package monitoring

import "encoding/json"

type MonitoringRoleConfig struct {
	Port int32 `json:"port"`
}

func (r *MonitoringRole) decodeMonitoringRoleConfig(raw []byte) *MonitoringRoleConfig {
	def := MonitoringRoleConfig{
		Port: 8009,
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
