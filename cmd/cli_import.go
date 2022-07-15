package cmd

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"beryju.io/gravity/pkg/instance"
	"github.com/spf13/cobra"
)

var importCmd = &cobra.Command{
	Use:   "import [file [file]]",
	Args:  cobra.MinimumNArgs(1),
	Short: "Import JSON file created with `export` into database",
	Run: func(cmd *cobra.Command, args []string) {
		rootInst := instance.NewInstance()
		for _, path := range args {
			cont, err := ioutil.ReadFile(path)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			var entries []instance.ExportEntry
			err = json.Unmarshal(cont, &entries)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
			err = rootInst.Import(entries)
			if err != nil {
				fmt.Println(err.Error())
				continue
			}
		}
	},
}

func init() {
	cliCmd.AddCommand(importCmd)
}
