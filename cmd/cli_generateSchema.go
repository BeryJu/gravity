package cmd

import (
	"os"
	"strings"

	"beryju.io/gravity/pkg/instance"
	"beryju.io/gravity/pkg/roles"
	"beryju.io/gravity/pkg/roles/api"
	"github.com/spf13/cobra"
)

var schemaFormat = ""

// generateSchemaCmd represents the generateSchema command
var generateSchemaCmd = &cobra.Command{
	Use:   "generateSchema [output_file]",
	Short: "Generate OpenAPI Schema",
	Run: func(cmd *cobra.Command, args []string) {
		rootInst := instance.NewInstance()
		inst := rootInst.ForRole("api")
		inst.AddEventListener(instance.EventTopicInstanceBootstrapped, func(ev *roles.Event) {
			defer rootInst.Stop()
			api := api.New(inst)
			schema := api.Schema()
			var out []byte
			var err error
			switch strings.ToLower(schemaFormat) {
			case "yaml":
				out, err = schema.MarshalYAML()
			case "json":
			default:
				out, err = schema.MarshalJSON()
			}
			if err != nil {
				rootInst.Log().WithError(err).Warning("failed to generate schema")
				return
			}
			if len(args) > 0 {
				err := os.WriteFile(args[0], out, 0644)
				if err != nil {
					rootInst.Log().WithError(err).Warning("failed to write schema")
				}
				rootInst.Log().Infof("Successfully wrote schema to %s", args[0])
			} else {
				os.Stdout.Write(out)
				rootInst.Log().Info("Successfully wrote schema to stdout")
			}
		})
		rootInst.Start()
	},
}

func init() {
	cliCmd.AddCommand(generateSchemaCmd)
	addUserCmd.PersistentFlags().StringVarP(&schemaFormat, "format", "f", "yaml", "Output format (yaml/json)")
}
