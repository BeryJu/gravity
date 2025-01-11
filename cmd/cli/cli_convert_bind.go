package cli

import (
	"context"

	"beryju.io/gravity/pkg/convert/bind"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var cliConverrtBind = &cobra.Command{
	Use:   "bind [zone_file [zone_file]]",
	Short: "Import Bind zone files into Gravity",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()
		for _, file := range args {
			conv, err := bind.New(apiClient, file)
			if err != nil {
				logger.Warn("failed to convert", zap.String("file", file), zap.Error(err))
				continue
			}
			err = conv.Run(ctx)
			if err != nil {
				logger.Warn("failed to convert zone", zap.Error(err))
			}
		}
	},
}

func init() {
	cliConvertCmd.AddCommand(cliConverrtBind)
}
