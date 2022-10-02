package cmd

import (
	"os"
	"strings"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/instance/types"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var schemaFormat = ""

func GenerateSchema(format string, callback func(schema []byte)) {
	rootInst := instance.New()
	inst := rootInst.ForRole("api")
	inst.AddEventListener(types.EventTopicInstanceBootstrapped, func(ev *roles.Event) {
		defer rootInst.Stop()
		api := api.New(inst)
		schema := api.Schema()
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
			rootInst.Log().WithError(err).Warning("failed to generate schema")
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
		GenerateSchema(schemaFormat, func(schema []byte) {
			if len(args) > 0 {
				err := os.WriteFile(args[0], schema, 0644)
				if err != nil {
					log.WithError(err).Warning("failed to write schema")
					return
				}
				log.Infof("Successfully wrote schema to %s", args[0])
			} else {
				cmd.OutOrStdout().Write(schema)
				log.Info("Successfully wrote schema to stdout")
			}
		})
	},
}

func init() {
	rootCmd.AddCommand(generateSchemaCmd)
	generateSchemaCmd.PersistentFlags().StringVarP(&schemaFormat, "format", "f", "yaml", "Output format (yaml/json)")
}
