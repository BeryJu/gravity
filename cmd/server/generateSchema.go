package server

import (
	"context"
	"os"
	"strings"

	"beryju.io/gravity/pkg/extconfig"
	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/instance/types"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api"
	"github.com/spf13/cobra"
	"go.uber.org/zap"
)

var schemaFormat = ""

func GenerateSchema(ctx context.Context, format string, callback func(schema []byte)) {
	rootInst := instance.New()
	inst := rootInst.ForRole("api", ctx)
	inst.AddEventListener(types.EventTopicInstanceBootstrapped, func(ev *roles.Event) {
		defer rootInst.Stop()
		api := rootInst.Role("api").(*api.Role)
		schema := api.Schema(ctx)
		var out []byte
		var err error
		switch strings.ToLower(format) {
		case "yaml":
			out, err = schema.MarshalYAML()
		case "json":
			fallthrough
		default:
			out, err = schema.MarshalJSON()
		}
		if err != nil {
			rootInst.Log().Warn("failed to generate schema", zap.Error(err))
			return
		}
		callback(out)
	})
	rootInst.Start()
}

// generateSchemaCmd represents the generateSchema command
var generateSchemaCmd = &cobra.Command{
	Use:   "generateSchema [output_file]",
	Short: "Generate OpenAPI Schema",
	Run: func(cmd *cobra.Command, args []string) {
		logger := extconfig.Get().Logger()
		GenerateSchema(cmd.Context(), schemaFormat, func(schema []byte) {
			if len(args) > 0 {
				err := os.WriteFile(args[0], schema, 0o644)
				if err != nil {
					logger.Warn("failed to write schema", zap.Error(err))
					return
				}
				logger.Info("successfully wrote schema", zap.String("to", args[0]))
			} else {
				_, err := cmd.OutOrStdout().Write(schema)
				if err != nil {
					logger.Warn("failed to write schema", zap.Error(err))
					return
				}
				logger.Info("Successfully wrote schema to stdout")
			}
		})
	},
}

func init() {
	if extconfig.Get().Debug {
		rootCmd.AddCommand(generateSchemaCmd)
	}
	generateSchemaCmd.PersistentFlags().StringVarP(&schemaFormat, "format", "f", "yaml", "Output format (yaml/json)")
}
