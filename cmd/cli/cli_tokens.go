package cli

import "github.com/spf13/cobra"

var cliTokensCmd = &cobra.Command{
	Use:   "tokens",
	Short: "Commands related to token management",
}

func init() {
	CLICmd.AddCommand(cliTokensCmd)
}
