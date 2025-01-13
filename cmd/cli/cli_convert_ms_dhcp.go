package cli

import (
	"context"
	"os"

	"beryju.io/gravity/pkg/convert/ms_dhcp"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var cliConverrtMSDHCPCmd = &cobra.Command{
	Use:   "ms_dhcp [input_xml [input_xml]]",
	Short: "Import Microsoft DHCP leases/reservations into gravity",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()
		for _, xml := range args {
			x, err := os.ReadFile(xml)
			if err != nil {
				logger.Warn("failed to read file", zap.Error(err), zap.String("file", xml))
				continue
			}
			conv, err := ms_dhcp.New(apiClient, string(x))
			if err != nil {
				logger.Warn("failed to convert", zap.String("file", xml), zap.Error(err))
				continue
			}
			err = conv.Run(ctx)
			if err != nil {
				logger.Warn("failed to convert file", zap.Error(err))
			}
		}
	},
}

func init() {
	cliConvertCmd.AddCommand(cliConverrtMSDHCPCmd)
}
