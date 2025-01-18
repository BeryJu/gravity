package server

import (
	"os"

	"beryju.io/gravity/cmd/cli"
	"beryju.io/gravity/pkg/extconfig"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "gravity",
	Short:   "Start gravity instance",
	Version: extconfig.FullVersion(),
}

func init() {
	rootCmd.AddCommand(cli.CLICmd)
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
