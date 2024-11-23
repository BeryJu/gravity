package cli

import (
	"os"

	"github.com/spf13/cobra"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check health and version",
	Run: func(cmd *cobra.Command, args []string) {
		v, hr, err := apiClient.ClusterInstancesApi.ClusterGetInfo(cmd.Context()).Execute()
		if err != nil {
			checkApiError(hr, err)
			os.Exit(1)
		}
		logger.Info(v.Version)
	},
}

func init() {
	CLICmd.AddCommand(healthCmd)
}
