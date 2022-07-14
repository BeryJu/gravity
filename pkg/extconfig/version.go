package extconfig

import (
	"strings"
)

var Version = "0.0.1"
var BuildHash = ""

func FullVersion() string {
	version := strings.Builder{}
	version.WriteString(Version)
	if BuildHash != "" {
		version.WriteRune('-')
		if len(BuildHash) >= 8 {
			version.WriteString(BuildHash[:8])
		} else {
			version.WriteString(BuildHash)
		}
	}
	return version.String()
}

func init() {
	if BuildHash == "" {
		BuildHash = "dev"
	}
}
