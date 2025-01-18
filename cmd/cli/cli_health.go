package cli

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var healthCmd = &cobra.Command{
	Use:   "health",
	Short: "Check health and version",
	Run: func(cmd *cobra.Command, args []string) {
		c, hr, err := apiClient.ClusterApi.ClusterGetClusterInfo(cmd.Context()).Execute()
		if err != nil {
			checkApiError(hr, err)
			os.Exit(1)
		}
		m := map[string]interface{}{
			"clusterVersion": c.ClusterVersion,
			"instances":      c.Instances,
		}
		b, err := json.MarshalIndent(m, "", "\t")
		if err != nil {
			logger.Warn("failed to render JSON", zap.Error(err))
			return
		}
		fmt.Println(string(b))
	},
}

func init() {
	CLICmd.AddCommand(healthCmd)
}
