package resources

import (
	"embed"
)

//go:embed blocky/*.txt
var BlockyLists embed.FS

//go:embed macoui/db.txt
var MacOUIDB []byte

//go:embed tftp/*
var TFTPRoot embed.FS
