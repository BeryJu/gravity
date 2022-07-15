package cmd

import (
	"beryju.io/gravity/pkg/instance"
	"github.com/spf13/cobra"
)

var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Run Gravity server",
	Run: func(cmd *cobra.Command, args []string) {
		inst := instance.NewInstance()
		defer inst.Stop()
		inst.Start()
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
}
