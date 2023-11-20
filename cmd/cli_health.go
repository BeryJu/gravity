package cmd

import (
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check health and version",
	Run: func(cmd *cobra.Command, args []string) {
		v, _, err := apiClient.ClusterInstancesApi.ClusterGetInfo(cmd.Context()).Execute()
		if err != nil {
			logger.Error("failed to get status", zap.Error(err))
			os.Exit(1)
		}
		logger.Info(v.Version)
	},
}

func init() {
	cliCmd.AddCommand(healthCmd)
}
