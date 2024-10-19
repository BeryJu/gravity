package cmd

import (
	"encoding/json"
	"os"

	"beryju.io/gravity/api"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var importCmd = &cobra.Command{
	Use:   "import [file [file]]",
	Args:  cobra.MinimumNArgs(1),
	Short: "Import JSON file created with `export` into database",
	Run: func(cmd *cobra.Command, args []string) {
		for _, path := range args {
			cont, err := os.ReadFile(path)
			if err != nil {
				logger.Error("failed to read import", zap.Error(err))
				continue
			}
			var entries api.ApiAPIImportInput
			err = json.Unmarshal(cont, &entries)
			if err != nil {
				logger.Error("failed to unmarshal", zap.Error(err))
				continue
			}
			hr, err := apiClient.RolesApiApi.ApiImport(cmd.Context()).ApiAPIImportInput(entries).Execute()
			if err != nil {
				checkApiError(hr, err)
				continue
			}
		}
	},
}

func init() {
	cliCmd.AddCommand(importCmd)
}
