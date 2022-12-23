package cmd

import (
	"context"

	"beryju.io/gravity/pkg/convert/ms_dhcp"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var convertMsDHCP = &cobra.Command{
	Use:   "ms_dhcp [input_xml [import_xml]]",
	Short: "Import Microsoft DHCP leases/reservations into gravity",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		for _, xml := range args {
			conv, err := ms_dhcp.New(apiClient, xml)
			if err != nil {
				logger.Warn("failed to convert", zap.String("file", xml), zap.Error(err))
				continue
			}
			conv.Run(ctx)
		}
	},
}

func init() {
	cliConvertCmd.AddCommand(convertMsDHCP)
}
