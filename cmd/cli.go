package cmd

import (
	"github.com/spf13/cobra"
)

var cliCmd = &cobra.Command{
	Use:   "cli",
	Short: "Interact with a running Gravity server",
}

func init() {
	rootCmd.AddCommand(cliCmd)
}
