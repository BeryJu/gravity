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
	Use:   "generateSchema",
	Short: "Generate OpenAPI Schema",
	Run: func(cmd *cobra.Command, args []string) {
		inst := instance.NewInstance()
		inst.ForRole("api").AddEventListener(instance.EventTopicInstanceBootstrapped, func(ev *roles.Event) {
			defer inst.Stop()
			api := api.New(inst.ForRole("api"))
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
				inst.Log().WithError(err).Warning("failed to generate schema")
				return
			}
			os.Stdout.Write(out)
		})
		inst.Start()
	},
}

func init() {
	rootCmd.AddCommand(generateSchemaCmd)
	addUserCmd.PersistentFlags().StringVarP(&schemaFormat, "format", "f", "yaml", "Output format (yaml/json)")
}
