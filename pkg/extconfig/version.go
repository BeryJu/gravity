package extconfig

import "fmt"

var Version = "0.0.1"
var BuildHash = ""

func FullVersion() string {
	return fmt.Sprintf("%s-%s", Version, BuildHash)
}

func init() {
	if BuildHash == "" {
		BuildHash = "dev"
	}
}
