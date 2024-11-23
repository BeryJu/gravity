package cli

import "github.com/spf13/cobra"

var cliConvertCmd = &cobra.Command{
	Use:   "convert",
	Short: "Conversion of existing DHCP/DNS data",
}

func init() {
	CLICmd.AddCommand(cliConvertCmd)
}
