package cli

import (
	"context"
	"os"

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
			x, err := os.Open(file)
			if err != nil {
				logger.Warn("failed to open file", zap.Error(err), zap.String("file", file))
				continue
			}
			defer func() {
				err := x.Close()
				if err != nil {
					logger.Warn("failed to close file", zap.Error(err))
				}
			}()

			conv, err := bind.New(apiClient, x)
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
