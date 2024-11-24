package extconfig

import (
	"os"
	"strings"
)

// Set via ldflags
var (
	Version   = ""
	BuildHash = ""
)

func FullVersion() string {
	if os.Getenv("CI") == "true" {
		Version = "99.99.99"
	}
	version := strings.Builder{}
	version.WriteString(Version)
	if BuildHash != "" {
		version.WriteRune('+')
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
