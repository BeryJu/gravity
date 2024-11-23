package cli

import "github.com/spf13/cobra"

var cliUsersCmd = &cobra.Command{
	Use:   "users",
	Short: "Commands related to user management",
}

func init() {
	CLICmd.AddCommand(cliUsersCmd)
}
