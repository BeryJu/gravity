package cmd

import (
	"context"
	"encoding/json"
	"os"

	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var exportCmd = &cobra.Command{
	Use:   "export [file]",
	Short: "Output entire database into JSON file",
	Run: func(cmd *cobra.Command, args []string) {
		exp, _, err := apiClient.RolesApiApi.ApiExport(context.Background()).Execute()
		if err != nil {
			logger.Error("failed to export", zap.Error(err))
			return
		}
		raw, err := json.Marshal(exp)
		if err != nil {
			logger.Error("failed to json marshal", zap.Error(err))
			return
		}
		err = os.WriteFile(args[0], raw, 0644)
		if err != nil {
			logger.Error("failed to write export", zap.Error(err))
		}
	},
}

func init() {
	cliCmd.AddCommand(exportCmd)
}
