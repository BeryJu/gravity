package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export [file]",
	Short: "Output entire database into JSON file",
	Run: func(cmd *cobra.Command, args []string) {
		exp, _, err := apiClient.RolesApiApi.ApiExport(context.Background()).Execute()
		if err != nil {
			fmt.Println(err)
			return
		}
		raw, err := json.Marshal(exp)
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		err = os.WriteFile(args[0], raw, 0644)
		if err != nil {
			fmt.Println(err.Error())
		}
	},
}

func init() {
	cliCmd.AddCommand(exportCmd)
}
