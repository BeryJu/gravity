package backup

import "encoding/json"

type BackupRoleConfig struct {
	Enabled   bool
	AccessKey string
	SecretKey string
	Endpoint  string
	Bucket    string
	CronExpr  string
}

func (r *BackupRole) decodeBackupRoleConfig(raw []byte) *BackupRoleConfig {
	def := BackupRoleConfig{
		Enabled:  true,
		CronExpr: "0 */24 * * *",
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
