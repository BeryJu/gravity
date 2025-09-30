//go:build debug
// +build debug

package cli

import (
	"github.com/go-delve/delve/cmd/dlv/cmds"
	"github.com/spf13/cobra"
)

var DebugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Utils to debug",
}

func init() {
	DebugCmd.AddCommand(cmds.New(false))
	CLICmd.AddCommand(DebugCmd)
}
