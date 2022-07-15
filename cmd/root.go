package cmd

import (
	"os"

	"beryju.io/gravity/pkg/extconfig"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:     "gravity",
	Short:   "Start gravity instance",
	Version: extconfig.FullVersion(),
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
