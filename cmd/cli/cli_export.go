package cli

import (
	"encoding/json"
	"os"

	"beryju.io/gravity/api"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var exportSafe = false

var exportCmd = &cobra.Command{
	Use:   "export [file]",
	Short: "Output entire database into JSON file",
	Run: func(cmd *cobra.Command, args []string) {
		exp, hr, err := apiClient.RolesApiApi.ApiExport(cmd.Context()).ApiAPIExportInput(api.ApiAPIExportInput{
			Safe: &exportSafe,
		}).Execute()
		if err != nil {
			checkApiError(hr, err)
			return
		}
		raw, err := json.Marshal(exp)
		if err != nil {
			logger.Error("failed to json marshal", zap.Error(err))
			return
		}
		if len(args) > 0 {
			err = os.WriteFile(args[0], raw, 0o644)
			if err != nil {
				logger.Error("failed to write export", zap.Error(err))
			}
		} else {
			_, err = os.Stdout.Write(raw)
			if err != nil {
				logger.Error("failed to write export", zap.Error(err))
			}
		}
	},
}

func init() {
	exportCmd.Flags().BoolVar(&exportSafe, "safe", false, "Export only safe values")
	CLICmd.AddCommand(exportCmd)
}
