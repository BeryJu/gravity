package cli

import (
	"context"
	"os"

	"beryju.io/gravity/pkg/convert/opnsense"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var cliConverrtOpnsense = &cobra.Command{
	Use:   "opnsense [config.xml [config.xml]]",
	Short: "Import Opnsense config files into Gravity",
	Run: func(cmd *cobra.Command, args []string) {
		ctx, cancel := context.WithCancel(cmd.Context())
		defer cancel()
		for _, file := range args {
			x, err := os.ReadFile(file)
			if err != nil {
				logger.Warn("failed to open file", zap.Error(err), zap.String("file", file))
				continue
			}

			conv, err := opnsense.New(apiClient, x)
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
	cliConvertCmd.AddCommand(cliConverrtOpnsense)
}
