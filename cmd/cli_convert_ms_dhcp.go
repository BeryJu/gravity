package cmd

import (
	"context"

	"beryju.io/gravity/pkg/convert/ms_dhcp"
	"beryju.io/gravity/pkg/instance"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var convertMsDHCP = &cobra.Command{
	Use:   "ms_dhcp [input_xml]",
	Short: "Import Microsoft DHCP leases/reservations into gravity",
	Run: func(cmd *cobra.Command, args []string) {
		rootInst := instance.New()
		conv, err := ms_dhcp.New(rootInst.ForRole("convert.ms_dhcp"), args[0])
		if err != nil {
			rootInst.Log().Warn("Failed to initialise converter", zap.Error(err))
			return
		}
		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()
		conv.Run(ctx)
	},
}

func init() {
	cliConvertCmd.AddCommand(convertMsDHCP)
}
