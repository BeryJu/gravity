package cmd

import (
	"encoding/json"
	"fmt"
	"os"

	"beryju.io/gravity/pkg/instance"
	"github.com/spf13/cobra"
)

var exportCmd = &cobra.Command{
	Use:   "export [file]",
	Short: "Output entire database into JSON file",
	Run: func(cmd *cobra.Command, args []string) {
		rootInst := instance.NewInstance()
		entries, err := rootInst.Export()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		raw, err := json.Marshal(entries)
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
