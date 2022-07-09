package dist

import (
	"embed"
)

//go:embed *
var Static embed.FS
